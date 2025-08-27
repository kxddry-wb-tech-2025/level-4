package http

import (
	"calendar/internal/models"
	"calendar/internal/models/log"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) sendLog(entry log.Entry) {
	if s.logs == nil {
		return
	}

	go func() {
		s.logs <- entry
	}()
}

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
		s.sendLog(log.Error(err, "failed to create event", echo.Map{
			"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
			"op":         "createEvent",
		}))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, event)
}

func (s *Server) getEvents(c echo.Context) error {
	events, err := s.svc.GetEvents(c.Request().Context())
	if err != nil {
		s.sendLog(log.Error(err, "failed to get events", echo.Map{
			"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
			"op":         "getEvents",
		}))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, events)
}

func (s *Server) getEvent(c echo.Context) error {
	id := c.Param("id")
	event, err := s.svc.GetEvent(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
		}

		s.sendLog(log.Error(err, "failed to get event", echo.Map{
			"id":         id,
			"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
			"op":         "getEvent",
		}))
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

		s.sendLog(log.Error(err, "failed to update event", echo.Map{
			"id":         id,
			"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
			"op":         "updateEvent",
		}))
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

		s.sendLog(log.Error(err, "failed to delete event", echo.Map{
			"id":         id,
			"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
			"op":         "deleteEvent",
		}))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
