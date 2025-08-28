package notify

import (
	"calendar/internal/models"
	"calendar/internal/models/log"
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
)

// Worker is the interface for the worker
type Worker interface {
	AddNotification(ctx context.Context, notification models.Notification) error
	DeleteNotification(ctx context.Context, id string) error
}

// NotificationRepository is the interface for the notification repository
type NotificationRepository interface {
	Create(ctx context.Context, notification models.CreateNotificationRequest) (string, error)
	GetIDsByEventID(ctx context.Context, eventID string) ([]string, error)
	DeleteAllByEventID(ctx context.Context, eventID string) error
}

// Service is the struct for the notification service
type Service struct {
	repo    NotificationRepository
	worker  Worker
	logs    chan<- log.Entry
	mainCtx context.Context
}

// Logs returns the channel for the logs
func (s *Service) Logs() <-chan log.Entry {
	logs := make(chan log.Entry, log.ChannelCapacity)
	s.logs = logs
	return logs
}

// NewService creates a new notification service
func NewService(ctx context.Context, repo NotificationRepository, worker Worker) *Service {
	return &Service{repo: repo, worker: worker, mainCtx: ctx}
}

// Process processes the notifications
func (s *Service) Process(jobs <-chan any) {
	for {
		select {
		case <-s.mainCtx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			switch j := job.(type) {
			case models.CreateNotificationRequest:
				err := s.createNotification(j)
				if err != nil {
					s.sendLog(log.Error(err, "failed to create notification", echo.Map{
						"op": "createNotification",
					}))
				}
			case models.DeleteNotificationsRequest:
				err := s.deleteNotifications(j.EventID)
				if err != nil {
					s.sendLog(log.Error(err, "failed to delete notifications", echo.Map{
						"op": "deleteNotifications",
					}))
				}
			default:
				s.sendLog(log.Error(fmt.Errorf("unknown job type: %T", j), "unknown job type", echo.Map{
					"op": "process",
				}))
			}
		}
	}
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

// createNotification creates a notification
func (s *Service) createNotification(notification models.CreateNotificationRequest) error {
	id, err := s.repo.Create(s.mainCtx, notification)
	if err != nil {
		s.sendLog(log.Error(err, "failed to create notification", echo.Map{
			"op": "createNotification",
		}))

		return err
	}

	err = s.worker.AddNotification(s.mainCtx, models.Notification{
		EventID:   notification.EventID,
		ID:        id,
		Message:   notification.Message,
		When:      notification.When,
		Channel:   notification.Channel,
		Recipient: notification.Recipient,
	})

	if err != nil {
		s.sendLog(log.Error(err, "failed to send notification", echo.Map{
			"op": "createNotification",
		}))

		return err
	}

	return nil
}

// deleteNotifications deletes all notifications for an event
func (s *Service) deleteNotifications(eventID string) error {
	ids, err := s.repo.GetIDsByEventID(s.mainCtx, eventID)
	if err != nil {
		s.sendLog(log.Error(err, "failed to get notification IDs", echo.Map{
			"op": "deleteNotifications",
		}))

		return err
	}

	for _, id := range ids {
		err := s.worker.DeleteNotification(s.mainCtx, id)
		if err != nil {
			s.sendLog(log.Error(err, "failed to cancel notification", echo.Map{
				"op": "deleteNotifications",
			}))
		}
	}

	err = s.repo.DeleteAllByEventID(s.mainCtx, eventID)
	if err != nil {
		s.sendLog(log.Error(err, "failed to delete notifications", echo.Map{
			"op": "deleteNotifications",
		}))
	}

	return nil
}
