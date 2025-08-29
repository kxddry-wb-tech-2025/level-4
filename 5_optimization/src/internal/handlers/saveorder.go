package handlers

import (
	"context"
	"l0/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// OrderSaver is an interface that can save orders
type OrderSaver interface {
	SaveOrder(context.Context, *models.Order) error
}

func CreateOrderHandler(saver OrderSaver) echo.HandlerFunc {
	return func(c echo.Context) error {
		var order models.Order
		if err := c.Bind(&order); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := c.Validate(order); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := saver.SaveOrder(c.Request().Context(), &order); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, order)
	}
}
