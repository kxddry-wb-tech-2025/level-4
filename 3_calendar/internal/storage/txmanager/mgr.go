package txmanager

import (
	"calendar/internal/lib/fanin"
	"calendar/internal/models/log"
	"calendar/internal/storage"
	"calendar/internal/storage/repos/archives"
	"calendar/internal/storage/repos/events"
	"calendar/internal/storage/repos/notifications"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TxManager is the transaction manager.
type TxManager struct {
	pool  *pgxpool.Pool
	repos *Repositories
}

// Repositories is the repositories for the transaction manager.
type Repositories struct {
	Events        *events.Repository
	Archives      *archives.Repository
	Notifications *notifications.Repository
}

// New creates a new transaction manager.
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

// Logs returns the logs channel.
func (m *TxManager) Logs(ctx context.Context) <-chan log.Entry {
	logs := fanin.FanIn(ctx, m.repos.Events.Logs(), m.repos.Archives.Logs(), m.repos.Notifications.Logs())
	return logs
}

// Close closes the transaction manager.
func (m *TxManager) Close() {
	m.pool.Close()
}

// tx is the transaction for the transaction manager.
type tx struct {
	ctx   context.Context
	repos *Repositories
	tx    pgx.Tx
}

// Rollback rolls back the transaction.
func (t *tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

// Commit commits the transaction.
func (t *tx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

// Do executes a function within a transaction.
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

// DoWith starts a transaction, injects it into the context, and runs the provided function.
// It commits on success and rolls back on error.
//
//nolint:gocognit
func (m *TxManager) DoWith(ctx context.Context, fn func(ctx context.Context) error) error {
	pgTx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}

	ctxWithTx := storage.WithTx(ctx, pgTx)

	if err := fn(ctxWithTx); err != nil {
		_ = pgTx.Rollback(ctx)
		return err
	}

	return pgTx.Commit(ctx)
}
