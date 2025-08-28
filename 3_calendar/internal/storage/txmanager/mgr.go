package txmanager

import (
	"calendar/internal/lib/fanin"
	"calendar/internal/models/log"
	"calendar/internal/storage/repos/archives"
	"calendar/internal/storage/repos/events"
	"calendar/internal/storage/repos/notifications"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager struct {
	pool  *pgxpool.Pool
	repos *Repositories
}

type Repositories struct {
	Events        *events.Repository
	Archives      *archives.Repository
	Notifications *notifications.Repository
}

func New(ctx context.Context, dsn string) (*TxManager, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &TxManager{
		pool: pool,
		repos: &Repositories{
			Events:        events.NewRepository(ctx, pool),
			Archives:      archives.NewRepository(ctx, pool),
			Notifications: notifications.NewRepository(ctx, pool),
		},
	}, nil
}

func (m *TxManager) Logs(ctx context.Context) <-chan log.Entry {
	logs := fanin.FanIn(ctx, m.repos.Events.Logs(), m.repos.Archives.Logs(), m.repos.Notifications.Logs())
	return logs
}

func (m *TxManager) Close() {
	m.pool.Close()
}

type tx struct {
	ctx   context.Context
	repos *Repositories
	tx    pgx.Tx
}

func (t *tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *tx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (m *TxManager) Do(ctx context.Context, fn func(ctx context.Context, tx *tx) error) error {
	pgTx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}

	t := &tx{
		ctx:   ctx,
		repos: m.repos,
		tx:    pgTx,
	}

	if err := fn(ctx, t); err != nil {
		_ = t.Rollback(ctx)
		return err
	}

	return t.Commit(ctx)
}
