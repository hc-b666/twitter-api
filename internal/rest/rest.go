package rest

import (
	"net/http"
	"time"
	"twitter-api/internal/rest/handler/comment"
	"twitter-api/internal/rest/handler/health"
	"twitter-api/internal/rest/handler/post"
	"twitter-api/internal/rest/handler/token"
	"twitter-api/internal/rest/handler/user"
	"twitter-api/internal/rest/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	mux            *gin.Engine
	healthHandler  *health.Handler
	userHandler    *user.Handler
	tokenHandler   *token.Handler
	postHandler    *post.Handler
	commentHandler *comment.Handler
	mw             *middleware.Middleware
}

func NewServer(
	mux *gin.Engine,
	healthHandler *health.Handler,
	userHandler *user.Handler,
	tokenHandler *token.Handler,
	postHandler *post.Handler,
	commentHandler *comment.Handler,
	mw *middleware.Middleware,
) *Server {
	return &Server{
		mux:            mux,
		healthHandler:  healthHandler,
		userHandler:    userHandler,
		tokenHandler:   tokenHandler,
		postHandler:    postHandler,
		commentHandler: commentHandler,
		mw:             mw,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Init() {
	const (
		baseURL    = "/api/v1"
		authURL    = "/auth"
		userURL    = "/user"
		postURL    = "/post"
		commentURL = "/comments"
	)
	s.mux.Use(gin.Logger())

	s.mux.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	group := s.mux.Group(baseURL)

	// Public routes
	group.GET("/health", s.healthHandler.HealthCheck)

	// Auth routes
	authGroup := group.Group(authURL)
	authGroup.POST("/register", s.userHandler.Register)
	authGroup.POST("/login", s.userHandler.Login)
	authGroup.POST("/refresh", s.tokenHandler.Refresh)

	// User routes
	userGroup := group.Group(userURL)
	userGroup.Use(s.mw.Authenticate())
	userGroup.GET("/profile", s.userHandler.Profile)

	// Post routes
	postGroup := group.Group(postURL)
	postGroup.POST(postURL, s.postHandler.CreateNewPost)

	// Comment routes
	commentGroup := group.Group(commentURL)
	commentGroup.POST(commentURL, s.commentHandler.CreateNewComment)
}
