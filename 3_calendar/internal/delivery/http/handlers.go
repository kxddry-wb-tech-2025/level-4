package http

import (
	"calendar/internal/models"
	"errors"
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
	events, err := s.svc.GetEvents(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, events)
}

func (s *Server) getEvent(c echo.Context) error {
	id := c.Param("id")
	event, err := s.svc.GetEvent(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, event)
}

func (s *Server) updateEvent(c echo.Context) error {
	id := c.Param("id")
	var req models.UpdateEventRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": err.Error()})
	}

	event, err := s.svc.UpdateEvent(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, event)
}

func (s *Server) deleteEvent(c echo.Context) error {
	id := c.Param("id")
	err := s.svc.DeleteEvent(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
