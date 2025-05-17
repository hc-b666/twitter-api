package rest

import (
	"net/http"
	"twitter-api/internal/rest/handler/health"

	"github.com/gin-gonic/gin"
)

type Server struct {
	mux           *gin.Engine
	healthHandler *health.Handler
}

func NewServer(
	mux *gin.Engine,
	healthHandler *health.Handler,
) *Server {
	return &Server{
		mux:           mux,
		healthHandler: healthHandler,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Init() {
	const (
		baseURL = "/api/v1"
	)
	s.mux.Use(gin.Logger())

	group := s.mux.Group(baseURL)

	// Public routes
	group.GET("/health", s.healthHandler.HealthCheck)

}
