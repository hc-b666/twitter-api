package rest

import (
	"net/http"
	"twitter-api/internal/rest/handler/health"
	"twitter-api/internal/rest/handler/token"
	"twitter-api/internal/rest/handler/user"
	"twitter-api/internal/rest/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	mux           *gin.Engine
	healthHandler *health.Handler
	userHandler   *user.Handler
	tokenHandler  *token.Handler
	mw            *middleware.Middleware
}

func NewServer(
	mux *gin.Engine,
	healthHandler *health.Handler,
	userHandler *user.Handler,
	tokenHandler *token.Handler,
	mw *middleware.Middleware,
) *Server {
	return &Server{
		mux:           mux,
		healthHandler: healthHandler,
		userHandler:   userHandler,
		tokenHandler:  tokenHandler,
		mw:            mw,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Init() {
	const (
		baseURL = "/api/v1"
		auth    = "/auth"
		user    = "/user"
	)
	s.mux.Use(gin.Logger())

	group := s.mux.Group(baseURL)

	// Public routes
	group.GET("/health", s.healthHandler.HealthCheck)

	// Auth routes
	authGroup := group.Group(auth)
	authGroup.POST("/register", s.userHandler.Register)
	authGroup.POST("/login", s.userHandler.Login)
	authGroup.POST("/refresh", s.tokenHandler.Refresh)

	// User routes
	userGroup := group.Group(user)
	userGroup.Use(s.mw.Authenticate())
	userGroup.GET("/profile", s.userHandler.Profile)
}
