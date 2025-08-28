package events

import (
	"calendar/internal/models"
	"calendar/internal/models/log"
	"calendar/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"
)

// EventRepository is the interface for the event repository
type EventRepository interface {
	Create(ctx context.Context, event models.CreateEventRequest) (string, error)
	GetAll(ctx context.Context) ([]models.Event, error)
	Get(ctx context.Context, id string) (models.Event, error)
	Update(ctx context.Context, id string, event models.UpdateEventRequest) error
	Delete(ctx context.Context, id string) error
}

// Service is the struct for the event service
type Service struct {
	repo    EventRepository
	logs    chan<- log.Entry
	jobs    chan<- any
	mainCtx context.Context
}

// Logs returns the channel for the logs
func (s *Service) Logs() <-chan log.Entry {
	logs := make(chan log.Entry, log.ChannelCapacity)
	s.logs = logs
	return logs
}

// NewService creates a new event service
func NewService(ctx context.Context, repo EventRepository, jobs chan<- any) *Service {
	return &Service{repo: repo, jobs: jobs, mainCtx: ctx}
}

// sendLog sends a log entry
func (s *Service) sendLog(entry log.Entry) {
	if s.logs == nil {
		return
	}
	if entry.Level >= log.LevelWarn {
		go func() {
			select {
			case s.logs <- entry:
			case <-s.mainCtx.Done():
			}
		}()
	} else {
		select {
		case s.logs <- entry:
		default:
		}
	}
}

// sendJob sends a job to the job channel
func (s *Service) sendJob(ctx context.Context, job any) {
	go func(ctx context.Context) {
		select {
		case s.jobs <- job:
		case <-ctx.Done():
		}
	}(ctx)
}

// CreateEvent creates a new event
func (s *Service) CreateEvent(ctx context.Context, event models.CreateEventRequest) (string, error) {
	id, err := s.repo.Create(ctx, event)
	if err != nil {
		s.sendLog(log.Error(err, "failed to create event", map[string]any{
			"op": "createEvent",
		}))
		return "", err
	}

	if event.Notify {
		s.sendJob(ctx, models.CreateNotificationRequest{
			EventID:   id,
			Message:   fmt.Sprintf(models.MessageTemplate, event.Title, event.Start.Format(time.RFC1123)),
			When:      models.NotifyTime(event.Start),
			Channel:   "email",
			Recipient: event.Email,
		})
	}

	return id, nil
}

// GetEvents gets all events
func (s *Service) GetEvents(ctx context.Context) ([]models.Event, error) {
	events, err := s.repo.GetAll(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return []models.Event{}, nil
		}
		s.sendLog(log.Error(err, "failed to get events", map[string]any{
			"op": "getEvents",
		}))

		return nil, err
	}

	return events, nil
}

// GetEvent gets an event by id
func (s *Service) GetEvent(ctx context.Context, id string) (models.Event, error) {
	event, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.Event{}, models.ErrNotFound
		}
		s.sendLog(log.Error(err, "failed to get event", map[string]any{
			"op": "getEvent",
		}))

		return models.Event{}, err
	}

	return event, nil
}

// UpdateEvent updates an event
func (s *Service) UpdateEvent(ctx context.Context, id string, new models.UpdateEventRequest) error {
	old, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.ErrNotFound
		}
		s.sendLog(log.Error(err, "failed to get event", map[string]any{
			"op": "getEvent",
		}))

		return err
	}
	if old.Notify {
		if (new.Email != "" && new.Email != old.Email) ||
			(new.Start != old.Start) {
			s.sendJob(ctx, models.DeleteNotificationsRequest{
				EventID: id,
			})
		}
	}

	if new.Notify {
		s.sendJob(ctx, models.CreateNotificationRequest{
			EventID:   id,
			Message:   fmt.Sprintf(models.MessageTemplate, new.Title, new.Start.Format(time.RFC1123)),
			When:      models.NotifyTime(new.Start),
			Channel:   "email",
			Recipient: new.Email,
		})
	}

	return s.repo.Update(ctx, id, new)
}

// DeleteEvent deletes an event
func (s *Service) DeleteEvent(ctx context.Context, id string) error {
	s.sendJob(ctx, models.DeleteNotificationsRequest{
		EventID: id,
	})

	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.ErrNotFound
		}
		s.sendLog(log.Error(err, "failed to delete event", map[string]any{
			"op": "deleteEvent",
		}))

		return err
	}

	return nil
}
