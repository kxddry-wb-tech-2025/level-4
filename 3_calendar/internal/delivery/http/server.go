package http

import (
	"calendar/internal/models/log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e    *echo.Echo
	logs chan<- log.Entry
}

func NewServer(logs chan<- log.Entry) *Server {
	e := echo.New()
	s := &Server{
		e:    e,
		logs: logs,
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
}

func (s *Server) Start(addr string) error {
	return s.e.Start(addr)
}
