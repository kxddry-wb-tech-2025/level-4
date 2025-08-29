package handlers

import (
	"context"
	"errors"
	"fmt"
	"l0/internal/models"
	"l0/internal/storage"
	"net/http"

	"github.com/labstack/echo/v4"
)

// OrderGetter gets orders
type OrderGetter interface {
	GetOrder(context.Context, string) (*models.Order, error)
}

// Cacher saves orders and fetches them quicker than persistent storage
type Cacher interface {
	OrderGetter
	OrderSaver
	LoadOrders(context.Context, []*models.Order) error
}

// GetOrderHandler handles GET requests
func GetOrderHandler(getter OrderGetter, cacher Cacher) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		id := c.Param("id")
		if id == "" {
			return echo.ErrNotFound
		}

		cache, err := cacher.GetOrder(ctx, id)
		if err == nil {
			return c.JSON(http.StatusOK, cache)
		}

		order, err := getter.GetOrder(ctx, id)
		if err != nil {
			if errors.Is(err, storage.ErrOrderNotFound) {
				return c.String(http.StatusNotFound, fmt.Sprintf("order %s not found", id))
			}
			return c.String(http.StatusInternalServerError, err.Error())
		}
		_ = cacher.SaveOrder(ctx, order) // nil always
		return c.JSON(http.StatusOK, order)
	}
}
