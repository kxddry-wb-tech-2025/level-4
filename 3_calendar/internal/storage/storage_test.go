package storage

import (
	"context"
	"testing"
)

type dummyTx struct{}

func TestWithTxAndGetTx(t *testing.T) {
	ctx := context.Background()
	// we don't have a real pgx.Tx; ensure GetTx returns false on empty ctx
	if _, ok := GetTx(ctx); ok {
		t.Fatalf("expected no tx in context")
	}

	// Storing nil would not pass type assertion, so ok should remain false
	ctx = context.WithValue(ctx, txKey{}, nil)
	if _, ok := GetTx(ctx); ok {
		t.Fatalf("expected ok=false when tx is nil of wrong type")
	}
}
