package txmanager

import (
	"calendar/internal/models"
	workerpkg "calendar/internal/service/worker"
	"context"
)

// workerTx is the transaction for the worker repository.
type workerTx struct {
	ctx context.Context
	t   *tx
}

// CreateNotification creates a notification.
func (w *workerTx) CreateNotification(ctx context.Context, n models.CreateNotificationRequest) (string, error) {
	return w.t.CreateNotification(ctx, n)
}

// DeleteNotificationByID deletes a notification by id.
func (w *workerTx) DeleteNotificationByID(ctx context.Context, id string) error {
	return w.t.DeleteNotificationByID(ctx, id)
}

// DeleteAllNotificationsByEventID deletes all notifications by event id.
func (w *workerTx) DeleteAllNotificationsByEventID(ctx context.Context, eventID string) error {
	return w.t.DeleteAllNotificationsByEventID(ctx, eventID)
}

// GetNotificationByID gets a notification by id.
func (w *workerTx) GetNotificationByID(ctx context.Context, id string) (models.Notification, error) {
	return w.t.GetNotificationByID(ctx, id)
}

// GetNotificationIDsByEventID gets all notification ids by event id.
func (w *workerTx) GetNotificationIDsByEventID(ctx context.Context, eventID string) ([]string, error) {
	return w.t.GetNotificationIDsByEventID(ctx, eventID)
}

// Commit commits the transaction.
func (w *workerTx) Commit() error {
	return w.t.Commit(w.ctx)
}

// Rollback rolls back the transaction.
func (w *workerTx) Rollback() error {
	return w.t.Rollback(w.ctx)
}

// WorkerTxManagerAdapter adapts TxManager to worker.TxManager.
type WorkerTxManagerAdapter struct {
	m *TxManager
}

// Do executes a function within a transaction.
func (a *WorkerTxManagerAdapter) Do(ctx context.Context, fn func(ctx context.Context, tx workerpkg.Tx) error) error {
	return a.m.Do(ctx, func(ctx context.Context, t *tx) error {
		return fn(ctx, &workerTx{ctx: ctx, t: t})
	})
}

// AsWorkerTxManager returns an adapter that satisfies worker.TxManager
func (m *TxManager) AsWorkerTxManager() workerpkg.TxManager {
	return &WorkerTxManagerAdapter{m: m}
}
