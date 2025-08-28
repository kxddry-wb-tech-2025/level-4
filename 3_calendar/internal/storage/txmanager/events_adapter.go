package txmanager

import (
	"calendar/internal/models"
	eventspkg "calendar/internal/service/events"
	"context"
)

// eventsTx is the transaction for the events repository.
type eventsTx struct {
	ctx context.Context
	t   *tx
}

// Event repository passthrough methods
func (e *eventsTx) CreateEvent(ctx context.Context, event models.CreateEventRequest) (string, error) {
	return e.t.CreateEvent(ctx, event)
}

// GetAllEvents gets all events.
func (e *eventsTx) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	return e.t.GetAllEvents(ctx)
}

// GetEvent gets an event by id.
func (e *eventsTx) GetEvent(ctx context.Context, id string) (models.Event, error) {
	return e.t.GetEvent(ctx, id)
}

// UpdateEvent updates an event.
func (e *eventsTx) UpdateEvent(ctx context.Context, id string, event models.UpdateEventRequest) error {
	return e.t.UpdateEvent(ctx, id, event)
}

// DeleteEvent deletes an event.
func (e *eventsTx) DeleteEvent(ctx context.Context, id string) error {
	return e.t.DeleteEvent(ctx, id)
}

// Transaction control without explicit ctx, as required by events.Tx
func (e *eventsTx) Commit() error {
	return e.t.Commit(e.ctx)
}

// Rollback rolls back the transaction.
func (e *eventsTx) Rollback() error {
	return e.t.Rollback(e.ctx)
}

// EventsTxManagerAdapter adapts TxManager to events.TxManager
type EventsTxManagerAdapter struct {
	m *TxManager
}

// Do executes a function within a transaction.
func (a *EventsTxManagerAdapter) Do(ctx context.Context, fn func(ctx context.Context, tx eventspkg.Tx) error) error {
	return a.m.Do(ctx, func(ctx context.Context, t *tx) error {
		return fn(ctx, &eventsTx{ctx: ctx, t: t})
	})
}

// AsEventsTxManager returns an adapter that satisfies events.TxManager
func (m *TxManager) AsEventsTxManager() eventspkg.TxManager {
	return &EventsTxManagerAdapter{m: m}
}
