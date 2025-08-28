package worker

import (
	"calendar/internal/config"
	"calendar/internal/models"
	"calendar/internal/models/log"
	"calendar/internal/service"
	"calendar/internal/storage"
	"context"
	"errors"
	"time"

	"github.com/labstack/echo/v4"
)

// Storage is the interface for the storage
type Storage interface {
	PopDue(ctx context.Context, key string, limit int64) ([]string, error)
	Enqueue(ctx context.Context, id string, at time.Time) error
	Remove(ctx context.Context, id string) error
}

// NotificationRepository is the interface for the notification repository
type NotificationRepository interface {
	service.NotificationCreator
	service.NotificationDeleter
	service.NotificationGetter
}

// Tx is the interface for the transaction.
type Tx interface {
	NotificationRepository
	Commit() error
	Rollback() error
}

// TxManager is the interface for the transaction manager.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error
}

// NotificationSender is the interface for the notification sender
type NotificationSender interface {
	Send(ctx context.Context, notifications []models.Notification) error
}

// Worker is the struct for the worker
type Worker struct {
	st       Storage
	txmgr    TxManager
	sender   NotificationSender
	logs     chan<- log.Entry
	interval time.Duration
	limit    int64
}

// Logs returns the channel for the logs
func (w *Worker) Logs() <-chan log.Entry {
	logs := make(chan log.Entry, log.ChannelCapacity)
	w.logs = logs
	return logs
}

// NewWorker creates a new worker
func NewWorker(st Storage, txmgr TxManager, sender NotificationSender, cfg *config.WorkerConfig) *Worker {
	return &Worker{st: st, txmgr: txmgr, sender: sender, interval: cfg.Interval, limit: cfg.Limit}
}

// sendLog sends a log entry
func (w *Worker) sendLog(ctx context.Context, entry log.Entry) {
	if w.logs == nil {
		return
	}

	if entry.Level >= log.LevelWarn {
		go func() {
			select {
			case w.logs <- entry:
			case <-ctx.Done():
			}
		}()
	} else {
		select {
		case w.logs <- entry:
		default:
		}
	}
}

// AddNotification schedules a notification to be sent at a specific time.
func (w *Worker) AddNotification(ctx context.Context, n models.CreateNotificationRequest) error {
	return w.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
		id, err := tx.CreateNotification(ctx, n)
		if err != nil {
			return err
		}

		if err := w.st.Enqueue(ctx, id, n.When); err != nil {
			w.sendLog(ctx, log.Error(err, "failed to enqueue notification", echo.Map{
				"op":      "AddNotification",
				"notifID": id,
			}))
			return err
		}
		return nil
	})
}

// DeleteNotificationByID deletes a specific notification by id and removes it from the queue
func (w *Worker) DeleteNotificationByID(ctx context.Context, id string) error {
	// best-effort removal from queue
	_ = w.st.Remove(ctx, id)
	return w.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
		return tx.DeleteNotificationByID(ctx, id)
	})
}

// Handle handles the worker
func (w *Worker) Handle(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ids, err := w.st.PopDue(ctx, "notify:due", w.limit)
			if err != nil {
				w.sendLog(ctx, log.Error(err, "failed to pop due", echo.Map{
					"op": "Handle",
				}))
				continue
			}

			if err := w.sendNotifications(ctx, ids); err != nil {
				w.sendLog(ctx, log.Error(err, "failed to send notifications", echo.Map{
					"op": "Handle",
				}))
				continue
			}
		}
	}
}

func (w *Worker) sendNotifications(ctx context.Context, ids []string) error {
	return w.txmgr.Do(ctx, func(ctx context.Context, tx Tx) error {
		notifications := []models.Notification{}
		for _, id := range ids {
			notification, err := tx.GetNotificationByID(ctx, id)
			if err != nil {
				// ON DELETE CASCADE deletes the notification from the database
				if errors.Is(err, storage.ErrNotFound) {
					continue
				}
				w.sendLog(ctx, log.Error(err, "failed to get notification", echo.Map{
					"op":      "Handle",
					"notifID": id,
				}))
				continue
			}

			notifications = append(notifications, notification)
		}

		if err := w.sender.Send(ctx, notifications); err != nil {
			w.sendLog(ctx, log.Error(err, "failed to send notifications", echo.Map{
				"op": "Handle",
			}))
		}

		for _, notification := range notifications {
			if err := tx.DeleteNotificationByID(ctx, notification.ID); err != nil {
				w.sendLog(ctx, log.Error(err, "failed to delete notification after send", echo.Map{
					"op":      "Handle",
					"notifID": notification.ID,
				}))
			}
		}

		return nil
	})
}
