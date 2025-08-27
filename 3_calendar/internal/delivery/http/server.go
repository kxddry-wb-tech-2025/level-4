package http

import (
	"calendar/internal/config"
	"calendar/internal/models"
	"calendar/internal/models/log"
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e       *echo.Echo
	logs    chan<- log.Entry
	port    int
	svc     Service
	mainCtx context.Context
}

type Service interface {
	CreateEvent(ctx context.Context, event models.CreateEventRequest) (string, error)
	GetEvents(ctx context.Context) ([]models.Event, error)
	GetEvent(ctx context.Context, id string) (models.Event, error)
	UpdateEvent(ctx context.Context, id string, event models.UpdateEventRequest) error
	DeleteEvent(ctx context.Context, id string) error
}

func NewServer(ctx context.Context, sCfg config.ServerConfig, service Service) *Server {
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}
	e.Server.ReadTimeout = sCfg.Timeout
	e.Server.WriteTimeout = sCfg.Timeout
	e.Server.IdleTimeout = sCfg.IdleTimeout

	s := &Server{
		e:       e,
		port:    sCfg.Port,
		svc:     service,
		mainCtx: ctx,
	}

	s.setup()

	return s
}

func (s *Server) Logs() <-chan log.Entry {
	logs := make(chan log.Entry, 100)
	s.logs = logs
	return logs
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func (s *Server) setup() {
	s.e.Use(middleware.Recover())

	s.setupRoutes()
}

func (s *Server) setupRoutes() {
	s.e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	s.e.POST("/events", s.createEvent)
	s.e.GET("/events", s.getEvents)
	s.e.GET("/events/:id", s.getEvent)
	s.e.PUT("/events/:id", s.updateEvent)
	s.e.DELETE("/events/:id", s.deleteEvent)

}

func (s *Server) Start() error {
	return s.e.Start(fmt.Sprintf(":%d", s.port))
}
