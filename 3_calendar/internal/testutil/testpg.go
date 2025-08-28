package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// TestPostgres holds a running postgres container and DSN
type TestPostgres struct {
	Ctn *postgres.PostgresContainer
	DSN string
}

// StartPostgres starts a postgres testcontainer and returns DSN
func StartPostgres(ctx context.Context) (*TestPostgres, error) {
	pg, err := postgres.RunContainer(ctx,
		tc.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithInitScripts(),
	)
	if err != nil {
		return nil, err
	}

	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = pg.Terminate(ctx)
		return nil, err
	}

	return &TestPostgres{Ctn: pg, DSN: dsn}, nil
}

// ApplySchema applies minimal schema required by repositories
func (tp *TestPostgres) ApplySchema(ctx context.Context) error {
	cfg, err := pgxpool.ParseConfig(tp.DSN)
	if err != nil {
		return err
	}
	cfg.MaxConns = 4
	cfg.MaxConnLifetime = time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	stmts := []string{
		// events table
		`CREATE TABLE IF NOT EXISTS events (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            title TEXT NOT NULL,
            description TEXT NOT NULL,
            start TIMESTAMPTZ NOT NULL,
            "end" TIMESTAMPTZ NOT NULL,
            notify BOOLEAN NOT NULL,
            email TEXT NOT NULL
        );`,
		// archives table
		`CREATE TABLE IF NOT EXISTS archives (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            event_id UUID NOT NULL,
            title TEXT NOT NULL,
            description TEXT NOT NULL,
            start TIMESTAMPTZ NOT NULL,
            "end" TIMESTAMPTZ NOT NULL,
            notify BOOLEAN NOT NULL,
            email TEXT NOT NULL
        );`,
		// notifications table
		`CREATE TABLE IF NOT EXISTS notifications (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            event_id UUID NOT NULL,
            message TEXT NOT NULL,
            "when" TIMESTAMPTZ NOT NULL,
            channel TEXT NOT NULL,
            recipient TEXT NOT NULL
        );`,
		// required extension for gen_random_uuid()
		`CREATE EXTENSION IF NOT EXISTS pgcrypto;`,
	}

	for i, s := range stmts {
		if _, err := pool.Exec(ctx, s); err != nil {
			return fmt.Errorf("apply schema step %d: %w", i, err)
		}
	}
	return nil
}

// Stop terminates the container
func (tp *TestPostgres) Stop(ctx context.Context) error {
	if tp.Ctn == nil {
		return nil
	}
	return tp.Ctn.Terminate(ctx)
}
