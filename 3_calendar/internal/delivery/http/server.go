package http

import (
	"calendar/internal/config"
	"calendar/internal/models/log"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e    *echo.Echo
	logs chan<- log.Entry
	port int
}

func NewServer(logs chan<- log.Entry, sCfg config.ServerConfig) *Server {
	e := echo.New()
	e.Server.ReadTimeout = sCfg.Timeout
	e.Server.WriteTimeout = sCfg.Timeout
	e.Server.IdleTimeout = sCfg.IdleTimeout

	s := &Server{
		e:    e,
		logs: logs,
		port: sCfg.Port,
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

func (s *Server) Start() error {
	return s.e.Start(fmt.Sprintf(":%d", s.port))
}
