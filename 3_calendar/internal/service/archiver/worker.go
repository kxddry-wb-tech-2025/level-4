package archiver

import (
	"calendar/internal/config"
	"calendar/internal/models"
	"calendar/internal/models/log"
	"calendar/internal/storage"
	"context"
	"errors"
	"time"
)

// EventRepository is the interface for the event repository.
type EventRepository interface {
	GetOldEvents(ctx context.Context, limit int) ([]models.Event, error)
	DeleteEvent(ctx context.Context, id string) error
}

// ArchiveRepository is the interface for the archive repository.
type ArchiveRepository interface {
	ArchiveEvent(ctx context.Context, event models.Event) error
}

// Tx is the interface for the transaction.
type Tx interface {
	EventRepository
	ArchiveRepository
	Commit() error
	Rollback() error
}

// TxManager is the interface for the transaction manager.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error
}

// Service is the service for the archiver.
type Service struct {
	txmgr   TxManager
	logs    chan<- log.Entry
	mainCtx context.Context
	cfg     config.ArchiverConfig
}

// New creates a new archiver service.
func New(ctx context.Context, cfg config.ArchiverConfig, txmgr TxManager) *Service {
	return &Service{
		txmgr:   txmgr,
		mainCtx: ctx,
		cfg:     cfg,
	}
}

// Logs returns the logs channel.
func (s *Service) Logs() <-chan log.Entry {
	logs := make(chan log.Entry, log.ChannelCapacity)
	s.logs = logs
	return logs
}

// sendLog sends a log entry to the logs channel.
func (s *Service) sendLog(e log.Entry) {
	if s.logs == nil {
		return
	}

	if e.Level >= log.LevelWarn {
		go func() {
			select {
			case s.logs <- e:
			case <-s.mainCtx.Done():
			}
		}()
	} else {
		select {
		case s.logs <- e:
		default:
		}
	}
}

// Run starts the archiver service.
func (s *Service) Run() {
	go func() {
		ticker := time.NewTicker(s.cfg.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-s.mainCtx.Done():
				return
			case <-ticker.C:
				s.archive()
			}
		}
	}()
}

// archive archives the events.
func (s *Service) archive() {
	if err := s.txmgr.Do(s.mainCtx, func(ctx context.Context, tx Tx) error {
		events, err := tx.GetOldEvents(ctx, s.cfg.BatchSize)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				return models.ErrNotFound
			}
			return err
		}

		for _, event := range events {
			if err := tx.ArchiveEvent(ctx, event); err != nil {
				return err
			}
		}

		return nil

	}); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			s.sendLog(log.Info("no events to archive", map[string]any{
				"batch_size": s.cfg.BatchSize,
				"op":         "archive",
			}))
			return
		}
		s.sendLog(log.Error(err, "failed to archive events", map[string]any{
			"batch_size": s.cfg.BatchSize,
			"op":         "archive",
		}))
		return
	}
}
