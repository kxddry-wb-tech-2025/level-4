package notify

import (
	"calendar/internal/models"
	"calendar/internal/models/log"
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
)

type NotificationRepository interface {
	Create(notification models.Notification) error
	GetIDsByEventID(eventID string) ([]string, error)
	DeleteAllByEventID(eventID string) error
}

type NotificationSender interface {
	Send(ctx context.Context, notification models.CreateNotificationRequest) (string, error)
	Cancel(ctx context.Context, id string) error
}

type Service struct {
	sender NotificationSender
	repo   NotificationRepository
	logs   chan<- log.Entry
}

func NewService(sender NotificationSender, repo NotificationRepository) *Service {
	return &Service{sender: sender, repo: repo}
}

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

func (s *Service) createNotification(ctx context.Context, notification models.CreateNotificationRequest) error {
	id, err := s.sender.Send(ctx, notification)
	if err != nil {
		s.sendLog(ctx, log.Error(err, "failed to create notification", echo.Map{
			"op": "createNotification",
		}))

		return err
	}

	err = s.repo.Create(models.Notification{
		EventID:   notification.EventID,
		ID:        id,
		Message:   notification.Message,
		When:      notification.When,
		Channel:   notification.Channel,
		Recipient: notification.Recipient,
	})
	if err != nil {
		s.sendLog(ctx, log.Error(err, "failed to create notification", echo.Map{
			"op": "createNotification",
		}))

		return err
	}

	return nil
}

func (s *Service) deleteNotifications(ctx context.Context, eventID string) error {
	ids, err := s.repo.GetIDsByEventID(eventID)
	if err != nil {
		s.sendLog(ctx, log.Error(err, "failed to get notification IDs", echo.Map{
			"op": "deleteNotifications",
		}))

		return err
	}

	for _, id := range ids {
		err := s.sender.Cancel(ctx, id)
		if err != nil {
			s.sendLog(ctx, log.Error(err, "failed to cancel notification", echo.Map{
				"op": "deleteNotifications",
			}))
		}
	}

	err = s.repo.DeleteAllByEventID(eventID)
	if err != nil {
		s.sendLog(ctx, log.Error(err, "failed to delete notifications", echo.Map{
			"op": "deleteNotifications",
		}))
	}

	return nil
}
