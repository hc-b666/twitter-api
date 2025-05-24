package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/uploadcare/uploadcare-go/ucare"
	"go.uber.org/dig"

	"twitter-api/internal/rest"
	"twitter-api/internal/rest/middleware"
	"twitter-api/pkg/db"
	logger "twitter-api/pkg/logger"

	healthHandler "twitter-api/internal/rest/handler/health"

	userRepo "twitter-api/internal/repo/user"
	userHandler "twitter-api/internal/rest/handler/user"
	userService "twitter-api/internal/service/user"

	tokenRepo "twitter-api/internal/repo/token"
	tokenHandler "twitter-api/internal/rest/handler/token"
	tokenService "twitter-api/internal/service/token"

	postRepo "twitter-api/internal/repo/post"
	postHandler "twitter-api/internal/rest/handler/post"
	postService "twitter-api/internal/service/post"

	commentRepo "twitter-api/internal/repo/comment"
	commentHandler "twitter-api/internal/rest/handler/comment"
	commentService "twitter-api/internal/service/comment"
)

func main() {
	var (
		port = "9999"
		host = "0.0.0.0"
		dsn  = "postgres://postgres:postgres@localhost:5432/db"
	)

	if err := execute(host, port, dsn); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(host, port, dsn string) error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	creds := ucare.APICreds{
		SecretKey: os.Getenv("UPLOADCARE_SECRET_KEY"),
		PublicKey: os.Getenv("UPLOADCARE_PUBLIC_KEY"),
	}

	conf := &ucare.Config{
		SignBasedAuthentication: true,
		APIVersion:              ucare.APIv06,
	}

	client, err := ucare.NewClient(creds, conf)
	if err != nil {
		return fmt.Errorf("failed to create uploadcare client: %w", err)
	}

	deps := []interface{}{
		func() (*pgxpool.Pool, error) {
			return db.NewDB(dsn)
		},
		gin.New,
		rest.NewServer,
		func(server *rest.Server) *http.Server {
			return &http.Server{
				Addr:              net.JoinHostPort(host, port),
				Handler:           server,
				ReadHeaderTimeout: 5 * time.Second,
			}
		},
		func() (*logger.Logger, error) {
			return logger.NewLogger("app.log")
		},
		func() ucare.Client {
			return client
		},
		middleware.New,
		healthHandler.NewHandler,
		userRepo.NewRepo,
		userService.NewService,
		userHandler.NewHandler,
		tokenRepo.NewRepo,
		tokenService.NewService,
		tokenHandler.NewHandler,
		postRepo.NewRepo,
		postService.NewService,
		postHandler.NewHandler,
		commentRepo.NewRepo,
		commentService.NewService,
		commentHandler.NewHandler,
	}

	container := dig.New()
	for _, dep := range deps {
		if err := container.Provide(dep); err != nil {
			return fmt.Errorf("failed to provide dependency: %w", err)
		}
	}

	err = container.Invoke(func(server *rest.Server) {
		server.Init()
	})
	if err != nil {
		return fmt.Errorf("failed to invoke server: %w", err)
	}

	if err := container.Invoke(func(server *http.Server) error {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("failed to start server: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("container.Invoke failed: %w", err)
	}

	return nil
}
