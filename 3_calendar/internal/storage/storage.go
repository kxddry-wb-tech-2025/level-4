package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// ErrNotFound is the error for not found
var (
	ErrNotFound = errors.New("not found")
)

type txKey struct{}

// WithTx adds a transaction to the context.
func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// GetTx gets a transaction from the context.
func GetTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}
