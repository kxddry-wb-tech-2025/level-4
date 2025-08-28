package archiver

import (
	"calendar/internal/config"
	"calendar/internal/models"
	"calendar/internal/models/log"
	"calendar/internal/storage"
	"context"
	"testing"
)

type fakeTx struct {
	events   []models.Event
	archived []models.Event
	getErr   error
}

func (f *fakeTx) GetOldEvents(ctx context.Context, limit int) ([]models.Event, error) {
	return f.events, f.getErr
}
func (f *fakeTx) DeleteEvent(ctx context.Context, id string) error { return nil }
func (f *fakeTx) ArchiveEvent(ctx context.Context, event models.Event) error {
	f.archived = append(f.archived, event)
	return nil
}
func (f *fakeTx) Commit() error   { return nil }
func (f *fakeTx) Rollback() error { return nil }

type fakeTxMgr struct{ tx *fakeTx }

func (m fakeTxMgr) Do(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error {
	return fn(ctx, m.tx)
}

func TestArchive_NoEvents_LogsInfo(t *testing.T) {
	ctx := context.Background()
	cfg := config.ArchiverConfig{BatchSize: 10}
	s := New(ctx, cfg, fakeTxMgr{tx: &fakeTx{getErr: storage.ErrNotFound}})
	logs := s.Logs()
	s.archive()
	e := <-logs
	if e.Level != log.LevelInfo {
		t.Fatalf("expected info level, got %v", e.Level)
	}
}

func TestArchive_ArchivesEvents(t *testing.T) {
	ctx := context.Background()
	tx := &fakeTx{events: []models.Event{{ID: "id1"}}}
	cfg := config.ArchiverConfig{BatchSize: 10}
	s := New(ctx, cfg, fakeTxMgr{tx: tx})
	_ = s.Logs()
	s.archive()
	if len(tx.archived) != 1 {
		t.Fatalf("expected 1 archived, got %d", len(tx.archived))
	}
}
