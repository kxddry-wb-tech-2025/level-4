package http

import (
	"calendar/internal/models"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	v10 "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type fakeSvc struct{}

func (fakeSvc) CreateEvent(ctx context.Context, event models.CreateEventRequest) (string, error) {
	return "id1", nil
}
func (fakeSvc) GetEvents(ctx context.Context) ([]models.Event, error) {
	return []models.Event{{ID: "id1"}}, nil
}
func (fakeSvc) GetEvent(ctx context.Context, id string) (models.Event, error) {
	return models.Event{ID: id}, nil
}
func (fakeSvc) UpdateEvent(ctx context.Context, id string, event models.UpdateEventRequest) error {
	return nil
}
func (fakeSvc) DeleteEvent(ctx context.Context, id string) error { return nil }

func TestRoutes_BasicFlow(t *testing.T) {
	e := echo.New()
	s := &Server{e: e, svc: fakeSvc{}, port: 0, mainCtx: context.Background()}
	e.Validator = &Validator{validator: v10.New()}
	e.POST("/events", s.createEvent)
	e.GET("/events", s.getEvents)
	e.GET("/events/:id", s.getEvent)
	e.PUT("/events/:id", s.updateEvent)
	e.DELETE("/events/:id", s.deleteEvent)

	// create
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(`{"title":"t","description":"d","start":"2025-01-01T00:00:00Z","end":"2025-01-01T01:00:00Z","notify":false,"email":"a@b.c"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: %d body=%s", rec.Code, rec.Body.String())
	}

	// list
	req = httptest.NewRequest(http.MethodGet, "/events", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list: %d", rec.Code)
	}

	// get
	req = httptest.NewRequest(http.MethodGet, "/events/id1", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("get: %d", rec.Code)
	}

	// update
	req = httptest.NewRequest(http.MethodPut, "/events/id1", strings.NewReader(`{"title":"t","description":"d","start":"2025-01-01T00:00:00Z","end":"2025-01-01T01:00:00Z","notify":false,"email":"a@b.c"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("update: %d", rec.Code)
	}

	// delete
	req = httptest.NewRequest(http.MethodDelete, "/events/id1", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("delete: %d", rec.Code)
	}
}
