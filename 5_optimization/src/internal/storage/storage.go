package storage

import (
	c "context"
	"errors"
	"l0/internal/models"
)

var (
	// ErrOrderNotFound explicitly states the order was not found
	ErrOrderNotFound = errors.New("order not found")
)

// Storage can save and get orders
type Storage interface {
	SaveOrder(c.Context, *models.Order) error
	GetOrder(c.Context, string) (*models.Order, error)
}
