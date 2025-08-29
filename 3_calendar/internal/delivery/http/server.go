package http

import (
	"calendar/internal/config"
	"calendar/internal/delivery/validate"
	"calendar/internal/models"
	"calendar/internal/models/log"
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server is the HTTP server for the application
type Server struct {
	e       *echo.Echo
	logs    chan<- log.Entry
	port    int
	svc     Service
	mainCtx context.Context
}

// Service is the interface for the application services
type Service interface {
	CreateEvent(ctx context.Context, event models.CreateEventRequest) (string, error)
	GetEvents(ctx context.Context) ([]models.Event, error)
	GetEvent(ctx context.Context, id string) (models.Event, error)
	UpdateEvent(ctx context.Context, id string, event models.UpdateEventRequest) error
	DeleteEvent(ctx context.Context, id string) error
}

// NewServer creates a new HTTP server
func NewServer(ctx context.Context, sCfg *config.ServerConfig, service Service) *Server {
	e := echo.New()
	e.Validator = validate.New()
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

// Logs returns the channel for the logs
func (s *Server) Logs() <-chan log.Entry {
	logs := make(chan log.Entry, log.ChannelCapacity)
	s.logs = logs
	return logs
}

// setup sets up the server
func (s *Server) setup() {
	s.e.Use(middleware.Recover())

	s.setupRoutes()
}

// setupRoutes sets up the routes
func (s *Server) setupRoutes() {
	s.e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	s.e.POST("/events", s.createEvent)
	s.e.GET("/events", s.getEvents)
	s.e.GET("/events/:id", s.getEvent)
	s.e.PUT("/events/:id", s.updateEvent)
	s.e.DELETE("/events/:id", s.deleteEvent)

	s.e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

}

// Start starts the server
func (s *Server) Start() error {
	return s.e.Start(fmt.Sprintf(":%d", s.port))
}
