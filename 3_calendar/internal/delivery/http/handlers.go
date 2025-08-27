package http

import (
	"calendar/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) createEvent(c echo.Context) error {
	var req models.CreateEventRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": err.Error()})
	}

	event, err := s.svc.CreateEvent(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, event)
}

func (s *Server) getEvents(c echo.Context) error {
	return nil
}

func (s *Server) getEvent(c echo.Context) error {
	return nil
}

func (s *Server) updateEvent(c echo.Context) error {
	return nil
}

func (s *Server) deleteEvent(c echo.Context) error {
	return nil
}
