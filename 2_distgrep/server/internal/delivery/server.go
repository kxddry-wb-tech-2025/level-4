package delivery

import (
	"fmt"
	"grep-server/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Server is the main server struct
type Server struct {
	e    *echo.Echo
	srvc Service
}

// Service is the interface for the service layer
type Service interface {
	Grep(req models.Request) (models.Response, error)
}

// NewServer creates a new server
func NewServer(srvc Service) *Server {
	e := echo.New()

	s := &Server{e: e, srvc: srvc}

	s.registerRoutes()
	return s
}

// registerRoutes registers the routes for the server
func (s *Server) registerRoutes() {
	s.e.POST("/grep", s.grep)
	s.e.GET("/health", s.health)
}

// grep is the handler for the grep endpoint
func (s *Server) grep(c echo.Context) error {
	var req models.Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	resp, err := s.srvc.Grep(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// health is the handler for the health endpoint
func (s *Server) health(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// Start starts the server
func (s *Server) Start(port int) error {
	return s.e.Start(fmt.Sprintf(":%d", port))
}
