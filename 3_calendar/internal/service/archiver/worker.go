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

type EventRepository interface {
	GetOldEvents(ctx context.Context, limit int) ([]models.Event, error)
	DeleteEvent(ctx context.Context, id string) error
}

type ArchiveRepository interface {
	ArchiveEvent(ctx context.Context, event models.Event) error
}

type Tx interface {
	EventRepository
	ArchiveRepository
	Commit() error
	Rollback() error
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error
}

type Service struct {
	txmgr   TxManager
	logs    chan<- log.Entry
	mainCtx context.Context
	cfg     config.ArchiverConfig
}

func New(ctx context.Context, cfg config.ArchiverConfig, txmgr TxManager) *Service {
	return &Service{
		txmgr:   txmgr,
		mainCtx: ctx,
		cfg:     cfg,
	}
}

func (s *Service) Logs() <-chan log.Entry {
	logs := make(chan log.Entry, log.ChannelCapacity)
	s.logs = logs
	return logs
}

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
