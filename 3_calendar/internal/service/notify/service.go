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
	repo   NotificationRepository
	worker Worker
	logs   chan<- log.Entry
}

// Logs returns the channel for the logs
func (s *Service) Logs() <-chan log.Entry {
	logs := make(chan log.Entry, 100)
	s.logs = logs
	return logs
}

// NewService creates a new notification service
func NewService(repo NotificationRepository, worker Worker) *Service {
	return &Service{repo: repo, worker: worker}
}

// Process processes the notifications
func (s *Service) Process(ctx context.Context, jobs <-chan any) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			switch j := job.(type) {
			case models.CreateNotificationRequest:
				err := s.createNotification(ctx, j)
				if err != nil {
					s.sendLog(ctx, log.Error(err, "failed to create notification", echo.Map{
						"op": "createNotification",
					}))
				}
			case models.DeleteNotificationsRequest:
				err := s.deleteNotifications(ctx, j.EventID)
				if err != nil {
					s.sendLog(ctx, log.Error(err, "failed to delete notifications", echo.Map{
						"op": "deleteNotifications",
					}))
				}
			default:
				s.sendLog(ctx, log.Error(fmt.Errorf("unknown job type: %T", j), "unknown job type", echo.Map{
					"op": "process",
				}))
			}
		}
	}
}

// sendLog sends a log entry
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

// createNotification creates a notification
func (s *Service) createNotification(ctx context.Context, notification models.CreateNotificationRequest) error {
	id, err := s.repo.Create(ctx, notification)
	if err != nil {
		s.sendLog(ctx, log.Error(err, "failed to create notification", echo.Map{
			"op": "createNotification",
		}))

		return err
	}

	err = s.worker.AddNotification(ctx, models.Notification{
		EventID:   notification.EventID,
		ID:        id,
		Message:   notification.Message,
		When:      notification.When,
		Channel:   notification.Channel,
		Recipient: notification.Recipient,
	})

	if err != nil {
		s.sendLog(ctx, log.Error(err, "failed to send notification", echo.Map{
			"op": "createNotification",
		}))

		return err
	}

	return nil
}

// deleteNotifications deletes all notifications for an event
func (s *Service) deleteNotifications(ctx context.Context, eventID string) error {
	ids, err := s.repo.GetIDsByEventID(ctx, eventID)
	if err != nil {
		s.sendLog(ctx, log.Error(err, "failed to get notification IDs", echo.Map{
			"op": "deleteNotifications",
		}))

		return err
	}

	for _, id := range ids {
		err := s.worker.DeleteNotification(ctx, id)
		if err != nil {
			s.sendLog(ctx, log.Error(err, "failed to cancel notification", echo.Map{
				"op": "deleteNotifications",
			}))
		}
	}

	err = s.repo.DeleteAllByEventID(ctx, eventID)
	if err != nil {
		s.sendLog(ctx, log.Error(err, "failed to delete notifications", echo.Map{
			"op": "deleteNotifications",
		}))
	}

	return nil
}
