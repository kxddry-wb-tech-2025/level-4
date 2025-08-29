package events

import (
	"calendar/internal/models"
	"calendar/internal/models/log"
	"calendar/internal/service"
	"calendar/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"
)

// EventRepository is the interface for the event repository
type EventRepository interface {
	service.EventCreator
	service.EventCreatorID
	service.EventGetter
	service.EventUpdater
	service.EventDeleter
}

// Tx is the interface for the transaction.
type Tx interface {
	EventRepository
	Commit() error
	Rollback() error
}

// TxManager is the interface for the transaction manager.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error
}

// Service is the struct for the event service.
type Service struct {
	txmgr   TxManager
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
func NewService(ctx context.Context, txmgr TxManager, jobs chan<- any) *Service {
	return &Service{txmgr: txmgr, jobs: jobs, mainCtx: ctx}
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
	if !event.End.After(event.Start) {
		return "", fmt.Errorf("%w: %s", models.ErrInvalidEvent, "end time must be after start time")
	}

	var id string
	if err := s.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
		var err error
		id, err = tx.CreateEvent(ctx, event)
		return err
	}); err != nil {
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
			Channel:   "email", // we can make this configurable, but we have to change the user interface first. I'm too lazy to implement it now.
			Recipient: event.Email,
		})
	}

	return id, nil
}

// GetEvents gets all events
func (s *Service) GetEvents(ctx context.Context) ([]models.Event, error) {
	var events []models.Event
	if err := s.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
		var err error
		events, err = tx.GetAllEvents(ctx)
		return err
	}); err != nil {
		s.sendLog(log.Error(err, "failed to get events", map[string]any{
			"op": "getEvents",
		}))
		if errors.Is(err, storage.ErrNotFound) {
			return []models.Event{}, nil
		}
		return nil, err
	}

	return events, nil
}

// GetEvent gets an event by id
func (s *Service) GetEvent(ctx context.Context, id string) (models.Event, error) {
	var event models.Event
	if err := s.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
		var err error
		event, err = tx.GetEvent(ctx, id)
		return err
	}); err != nil {
		s.sendLog(log.Error(err, "failed to get event", map[string]any{
			"op": "getEvent",
		}))

		if errors.Is(err, storage.ErrNotFound) {
			return models.Event{}, fmt.Errorf("%w: %s", models.ErrNotFound, id)
		}

		return models.Event{}, err
	}

	return event, nil
}

// UpdateEvent updates an event
func (s *Service) UpdateEvent(ctx context.Context, id string, new models.UpdateEventRequest) error {
	var old models.Event

	if !new.End.After(new.Start) {
		return fmt.Errorf("%w: %s", models.ErrInvalidEvent, "end time must be after start time")
	}

	if err := s.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
		var err error
		var deleted bool
		old, err = tx.GetEvent(ctx, id)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				return fmt.Errorf("%w: event %s not found", models.ErrNotFound, id)
			}
			return err
		}

		if old.Notify {
			// if email was changed, delete the old notifications
			// if start time was changed, delete the old notifications
			// if notify was turned off, delete the old notifications
			if (new.Email != "" && new.Email != old.Email) ||
				(new.Start != old.Start) || !new.Notify {
				// delete the old event, DELETE ON CASCADE deletes the notifications anyway
				err = tx.DeleteEvent(ctx, id)
				if err != nil {
					return err
				}
				deleted = true
			}
		}

		if deleted {
			err = tx.CreateEventWithID(ctx, id, models.CreateEventRequest(new))
		} else {
			err = tx.UpdateEvent(ctx, id, new)
		}
		if err != nil {
			return err
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

		return tx.UpdateEvent(ctx, id, new)
	}); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.ErrNotFound
		}
		s.sendLog(log.Error(err, "failed to get event", map[string]any{
			"op": "getEvent",
		}))
		return err
	}
	return nil
}

// DeleteEvent deletes an event
func (s *Service) DeleteEvent(ctx context.Context, id string) error {
	err := s.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
		return tx.DeleteEvent(ctx, id)
	})
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
