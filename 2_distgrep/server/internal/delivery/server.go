package delivery

import (
	"fmt"
	"grep-server/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	e    *echo.Echo
	srvc Service
}

type Service interface {
	Grep(req models.Request) (models.Response, error)
}

func NewServer(srvc Service) *Server {
	e := echo.New()

	s := &Server{e: e, srvc: srvc}

	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.e.POST("/grep", s.grep)
	s.e.GET("/health", s.health)
}

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

func (s *Server) health(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) Start(port int) error {
	return s.e.Start(fmt.Sprintf(":%d", port))
}
