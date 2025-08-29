package handlers

import (
	"context"
	"errors"
	"fmt"
	"l0/internal/metrics"
	"l0/internal/models"
	"l0/internal/storage"
	"net/http"
	"time"

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
		const route = "GET /order/:id"
		startDelivery := time.Now()
		startDomain := time.Now()
		ctx := c.Request().Context()

		id := c.Param("id")
		if id == "" {
			metrics.ObserveDomainDuration(route, time.Since(startDomain))
			metrics.ObserveDeliveryDuration(route, time.Since(startDelivery))
			return echo.ErrNotFound
		}

		cache, err := cacher.GetOrder(ctx, id)
		if err == nil {
			metrics.ObserveDomainDuration(route, time.Since(startDomain))
			metrics.ObserveDeliveryDuration(route, time.Since(startDelivery))
			return c.JSON(http.StatusOK, cache)
		}

		preRepoDomain := time.Since(startDomain)
		order, err := getter.GetOrder(ctx, id)
		if err != nil {
			if errors.Is(err, storage.ErrOrderNotFound) {
				metrics.ObserveDomainDuration(route, preRepoDomain)
				metrics.ObserveDeliveryDuration(route, time.Since(startDelivery))
				return c.String(http.StatusNotFound, fmt.Sprintf("order %s not found", id))
			}
			metrics.ObserveDomainDuration(route, preRepoDomain)
			metrics.ObserveDeliveryDuration(route, time.Since(startDelivery))
			return c.String(http.StatusInternalServerError, err.Error())
		}
		// domain after repo
		startDomain2 := time.Now()
		_ = cacher.SaveOrder(ctx, order) // nil always
		domainTotal := preRepoDomain + time.Since(startDomain2)
		metrics.ObserveDomainDuration(route, domainTotal)
		metrics.ObserveDeliveryDuration(route, time.Since(startDelivery))
		return c.JSON(http.StatusOK, order)
	}
}
