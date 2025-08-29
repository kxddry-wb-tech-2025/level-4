package http

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server is the HTTP server for the application
type Server struct {
	e *echo.Echo
}

// New creates a new Server
func New() *Server {
	e := echo.New()

	s := &Server{
		e: e,
	}

	s.routes()

	return s
}

// routes sets up the routes for the server
func (s *Server) routes() {
	s.e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}

func (s *Server) Start(port string) error {
	return s.e.Start(":" + port)
}
