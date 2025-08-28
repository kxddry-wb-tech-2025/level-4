package worker

import (
	"calendar/internal/config"
	"calendar/internal/models"
	"context"
	"errors"
	"testing"
	"time"
)

type fakeStorage struct {
	q       []string
	removed map[string]bool
}

func (f *fakeStorage) PopDue(ctx context.Context, key string, limit int64) ([]string, error) {
	out := f.q
	f.q = nil
	return out, nil
}
func (f *fakeStorage) Enqueue(ctx context.Context, id string, at time.Time) error {
	f.q = append(f.q, id)
	return nil
}
func (f *fakeStorage) Remove(ctx context.Context, id string) error {
	if f.removed == nil {
		f.removed = map[string]bool{}
	}
	f.removed[id] = true
	return nil
}

type fakeTx struct {
	notifs map[string]models.Notification
}

func (f *fakeTx) CreateNotification(ctx context.Context, n models.CreateNotificationRequest) (string, error) {
	id := "n1"
	f.notifs[id] = models.Notification{ID: id, EventID: n.EventID, Message: n.Message, When: n.When, Channel: n.Channel, Recipient: n.Recipient}
	return id, nil
}
func (f *fakeTx) DeleteNotificationByID(ctx context.Context, id string) error {
	delete(f.notifs, id)
	return nil
}
func (f *fakeTx) DeleteAllNotificationsByEventID(ctx context.Context, eid string) error {
	for id, n := range f.notifs {
		if n.EventID == eid {
			delete(f.notifs, id)
		}
	}
	return nil
}
func (f *fakeTx) GetNotificationByID(ctx context.Context, id string) (models.Notification, error) {
	n, ok := f.notifs[id]
	if !ok {
		return models.Notification{}, errors.New("missing")
	}
	return n, nil
}
func (f *fakeTx) GetNotificationIDsByEventID(ctx context.Context, eid string) ([]string, error) {
	ids := []string{}
	for id, n := range f.notifs {
		if n.EventID == eid {
			ids = append(ids, id)
		}
	}
	return ids, nil
}
func (f *fakeTx) Commit() error   { return nil }
func (f *fakeTx) Rollback() error { return nil }

type fakeTxMgr struct{ tx *fakeTx }

func (m fakeTxMgr) Do(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error {
	if m.tx == nil {
		m.tx = &fakeTx{notifs: map[string]models.Notification{}}
	}
	return fn(ctx, m.tx)
}

type fakeSender struct {
	sent []models.Notification
	err  error
}

func (s *fakeSender) Send(ctx context.Context, ns []models.Notification) error {
	s.sent = append(s.sent, ns...)
	return s.err
}

func TestAddNotification_Enqueue(t *testing.T) {
	ctx := context.Background()
	st := &fakeStorage{}
	m := fakeTxMgr{tx: &fakeTx{notifs: map[string]models.Notification{}}}
	w := NewWorker(st, m, &fakeSender{}, &config.WorkerConfig{Interval: 1 * time.Second, Limit: 100})

	when := time.Now().Add(time.Minute)
	if err := w.AddNotification(ctx, models.CreateNotificationRequest{EventID: "e1", Message: "m", When: when, Channel: "email", Recipient: "a@b.c"}); err != nil {
		t.Fatalf("AddNotification: %v", err)
	}
	if len(st.q) == 0 {
		t.Fatalf("expected enqueue to be called")
	}
}

func TestHandle_SendsAndDeletes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	st := &fakeStorage{q: []string{"n1"}}
	tx := &fakeTx{notifs: map[string]models.Notification{"n1": {ID: "n1", EventID: "e1", Message: "m", When: time.Now(), Channel: "email", Recipient: "a@b.c"}}}
	m := fakeTxMgr{tx: tx}
	s := &fakeSender{}
	w := NewWorker(st, m, s, &config.WorkerConfig{Interval: 1 * time.Second, Limit: 100})
	go w.Handle(ctx)
	time.Sleep(1200 * time.Millisecond)
	if len(s.sent) == 0 {
		t.Fatalf("expected notification sent")
	}
	if _, ok := tx.notifs["n1"]; ok {
		t.Fatalf("expected notification deleted after send")
	}
}
