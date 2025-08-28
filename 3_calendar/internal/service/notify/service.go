package notify

import (
	"calendar/internal/models"
	"calendar/internal/models/log"
	"calendar/internal/service"
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
)

// Worker is the interface for the worker
type Worker interface {
	service.NotificationAdder
	service.NotificationDeleter
}

// Service is the struct for the notification service
type Service struct {
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
func NewService(ctx context.Context, worker Worker) *Service {
	return &Service{worker: worker, mainCtx: ctx}
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
	err := s.worker.AddNotification(s.mainCtx, notification)
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
	return s.worker.DeleteAllNotificationsByEventID(s.mainCtx, eventID)
}
