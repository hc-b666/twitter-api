package middleware

import (
	"net/http"
	"twitter-api/internal/service/user"
	"twitter-api/pkg/logger"
	"twitter-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Middleware struct {
	db      *pgxpool.Pool
	userSvc *user.Service
	l       *logger.Logger
}

func New(
	db *pgxpool.Pool,
	userSvc *user.Service,
	l *logger.Logger,
) *Middleware {
	return &Middleware{
		db:      db,
		userSvc: userSvc,
		l:       l,
	}
}

func (m *Middleware) Authenticate() gin.HandlerFunc {
	const tokenHeaderLength = 7

	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < tokenHeaderLength || authHeader[:tokenHeaderLength] != "Bearer " {
			m.l.Error("missing or invalid Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		tokenStr := authHeader[tokenHeaderLength:]

		userID, userRole, err := utils.VerifyAccessToken(tokenStr)
		if err != nil {
			m.l.Error("failed to verify access token", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		u, err := m.userSvc.GetByID(c.Request.Context(), userID)
		if err != nil {
			m.l.Error("failed to get user by ID", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if u.Role != userRole {
			m.l.Error("user role does not match", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		c.Set("userID", userID)

		c.Next()
	}
}
