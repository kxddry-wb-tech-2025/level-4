package txmanager

import (
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

func (m *TxManager) Close() {
	m.pool.Close()
}

type tx struct {
	ctx   context.Context
	repos *Repositories
	tx    pgx.Tx
}
