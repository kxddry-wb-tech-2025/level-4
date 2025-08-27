package http

import (
	"calendar/internal/config"
	"calendar/internal/models"
	"calendar/internal/models/log"
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e    *echo.Echo
	logs chan<- log.Entry
	port int
	svc  Service
}

type Service interface {
	CreateEvent(ctx context.Context, event models.CreateEventRequest) (models.Event, error)
	GetEvents(ctx context.Context) ([]models.Event, error)
	GetEvent(ctx context.Context, id string) (models.Event, error)
	UpdateEvent(ctx context.Context, id string, event models.UpdateEventRequest) (models.Event, error)
	DeleteEvent(ctx context.Context, id string) error
}

func NewServer(logs chan<- log.Entry, sCfg config.ServerConfig, service Service) *Server {
	e := echo.New()
	e.Server.ReadTimeout = sCfg.Timeout
	e.Server.WriteTimeout = sCfg.Timeout
	e.Server.IdleTimeout = sCfg.IdleTimeout

	s := &Server{
		e:    e,
		logs: logs,
		port: sCfg.Port,
		svc:  service,
	}

	s.setup()

	return s
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
