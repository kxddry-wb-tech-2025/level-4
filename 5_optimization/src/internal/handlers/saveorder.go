package handlers

import (
	"context"
	"l0/internal/metrics"
	"l0/internal/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// OrderSaver is an interface that can save orders
type OrderSaver interface {
	SaveOrder(context.Context, *models.Order) error
}

func CreateOrderHandler(saver OrderSaver) echo.HandlerFunc {
	return func(c echo.Context) error {
		const route = "POST /order"
		startDelivery := time.Now()
		startDomain := time.Now()
		var order models.Order
		if err := c.Bind(&order); err != nil {
			metrics.ObserveDomainDuration(route, time.Since(startDomain))
			metrics.ObserveDeliveryDuration(route, time.Since(startDelivery))
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := c.Validate(order); err != nil {
			metrics.ObserveDomainDuration(route, time.Since(startDomain))
			metrics.ObserveDeliveryDuration(route, time.Since(startDelivery))
			return c.JSON(http.StatusBadRequest, err)
		}

		preRepoDomain := time.Since(startDomain)
		if err := saver.SaveOrder(c.Request().Context(), &order); err != nil {
			metrics.ObserveDomainDuration(route, preRepoDomain)
			metrics.ObserveDeliveryDuration(route, time.Since(startDelivery))
			return c.JSON(http.StatusInternalServerError, err)
		}

		// domain after repo
		startDomain2 := time.Now()
		resp := c.JSON(http.StatusOK, order)
		domainTotal := preRepoDomain + time.Since(startDomain2)
		metrics.ObserveDomainDuration(route, domainTotal)
		metrics.ObserveDeliveryDuration(route, time.Since(startDelivery))
		return resp
	}
}
