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

type EventRepository interface {
	Create(event models.CreateEventRequest) (string, error)
	GetAll() ([]models.Event, error)
	Get(id string) (models.Event, error)
	Update(id string, event models.UpdateEventRequest) error
	Delete(id string) error
}

type Service struct {
	repo EventRepository
	logs chan<- log.Entry
	jobs chan<- any
}

func NewService(repo EventRepository, jobs chan<- any) *Service {
	return &Service{repo: repo, jobs: jobs}
}

func (s *Service) sendLog(ctx context.Context, entry log.Entry) {
	if s.logs == nil {
		return
	}
	go func(ctx context.Context) {
		select {
		case s.logs <- entry:
		case <-ctx.Done():
		}
	}(ctx)
}

func (s *Service) sendJob(ctx context.Context, job any) {
	go func(ctx context.Context) {
		select {
		case s.jobs <- job:
		case <-ctx.Done():
		}
	}(ctx)
}

func (s *Service) CreateEvent(ctx context.Context, event models.CreateEventRequest) (string, error) {
	id, err := s.repo.Create(event)
	if err != nil {
		s.sendLog(ctx, log.Error(err, "failed to create event", map[string]any{
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

func (s *Service) GetEvents(ctx context.Context) ([]models.Event, error) {
	events, err := s.repo.GetAll()
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return []models.Event{}, nil
		}
		s.sendLog(ctx, log.Error(err, "failed to get events", map[string]any{
			"op": "getEvents",
		}))

		return nil, err
	}

	return events, nil
}

func (s *Service) GetEvent(ctx context.Context, id string) (models.Event, error) {
	event, err := s.repo.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.Event{}, models.ErrNotFound
		}
		s.sendLog(ctx, log.Error(err, "failed to get event", map[string]any{
			"op": "getEvent",
		}))

		return models.Event{}, err
	}

	return event, nil
}

func (s *Service) UpdateEvent(ctx context.Context, id string, new models.UpdateEventRequest) error {
	old, err := s.repo.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.ErrNotFound
		}
		s.sendLog(ctx, log.Error(err, "failed to get event", map[string]any{
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

	return s.repo.Update(id, new)
}

func (s *Service) DeleteEvent(ctx context.Context, id string) error {
	s.sendJob(ctx, models.DeleteNotificationsRequest{
		EventID: id,
	})

	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.ErrNotFound
		}
		s.sendLog(ctx, log.Error(err, "failed to delete event", map[string]any{
			"op": "deleteEvent",
		}))

		return err
	}

	return nil
}
