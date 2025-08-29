package handlers

import (
	"context"
	"fmt"
	"l0/internal/models"
)

// Database is an interface for a SQL or NoSQL database
type Database interface {
	OrderGetter
	AllOrders(ctx context.Context) ([]*models.Order, error)
}

// LoadCache fetches all orders from the database
func LoadCache(ctx context.Context, cacher Cacher, db Database) error {
	const op = "handlers.LoadCache"

	orders, err := db.AllOrders(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = cacher.LoadOrders(ctx, orders)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
