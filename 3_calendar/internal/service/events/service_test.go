package events

import (
	"calendar/internal/models"
	"context"
	"errors"
	"testing"
	"time"
)

type fakeTx struct {
	events map[string]models.Event
}

func (f *fakeTx) CreateEvent(ctx context.Context, req models.CreateEventRequest) (string, error) {
	id := "id-1"
	f.events[id] = models.Event{ID: id, Title: req.Title, Description: req.Description, Start: req.Start, End: req.End, Notify: req.Notify, Email: req.Email}
	return id, nil
}
func (f *fakeTx) CreateEventWithID(ctx context.Context, id string, req models.CreateEventRequest) error {
	f.events[id] = models.Event{ID: id, Title: req.Title, Description: req.Description, Start: req.Start, End: req.End, Notify: req.Notify, Email: req.Email}
	return nil
}
func (f *fakeTx) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	out := make([]models.Event, 0, len(f.events))
	for _, e := range f.events {
		out = append(out, e)
	}
	return out, nil
}
func (f *fakeTx) GetEvent(ctx context.Context, id string) (models.Event, error) {
	e, ok := f.events[id]
	if !ok {
		return models.Event{}, errors.New("missing")
	}
	return e, nil
}
func (f *fakeTx) UpdateEvent(ctx context.Context, id string, req models.UpdateEventRequest) error {
	if _, ok := f.events[id]; !ok {
		return errors.New("missing")
	}
	f.events[id] = models.Event{ID: id, Title: req.Title, Description: req.Description, Start: req.Start, End: req.End, Notify: req.Notify, Email: req.Email}
	return nil
}
func (f *fakeTx) DeleteEvent(ctx context.Context, id string) error {
	if _, ok := f.events[id]; !ok {
		return errors.New("missing")
	}
	delete(f.events, id)
	return nil
}
func (f *fakeTx) Commit() error   { return nil }
func (f *fakeTx) Rollback() error { return nil }

type fakeTxMgr struct{ tx *fakeTx }

func (m fakeTxMgr) Do(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error {
	if m.tx == nil {
		m.tx = &fakeTx{events: map[string]models.Event{}}
	}
	return fn(ctx, m.tx)
}

func TestCreateEvent_SendsJobWhenNotify(t *testing.T) {
	ctx := context.Background()
	jobs := make(chan any, 1)
	svc := NewService(ctx, fakeTxMgr{}, jobs)

	_, err := svc.CreateEvent(ctx, models.CreateEventRequest{
		Title: "t", Description: "d", Start: time.Now().Add(time.Hour), End: time.Now().Add(2 * time.Hour), Notify: true, Email: "a@b.c",
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	// job is sent via goroutine; allow scheduling
	time.Sleep(50 * time.Millisecond)
	select {
	case j := <-jobs:
		if _, ok := j.(models.CreateNotificationRequest); !ok {
			t.Fatalf("expected CreateNotificationRequest, got %T", j)
		}
	default:
		t.Fatalf("expected job to be enqueued")
	}
}

func TestGetEvents(t *testing.T) {
	ctx := context.Background()
	jobs := make(chan any, 1)
	svc := NewService(ctx, fakeTxMgr{tx: &fakeTx{events: map[string]models.Event{}}}, jobs)
	_, _ = svc.CreateEvent(ctx, models.CreateEventRequest{})
	evs, err := svc.GetEvents(ctx)
	if err != nil || len(evs) == 0 {
		t.Fatalf("unexpected: %v %v", evs, err)
	}
}
