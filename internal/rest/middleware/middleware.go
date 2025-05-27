package middleware

import (
	"net/http"
	"strings"
	"twitter-api/internal/service/user"
	"twitter-api/pkg/db"
	"twitter-api/pkg/logger"
	"twitter-api/pkg/types"
	"twitter-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	db      db.Pool
	userSvc *user.Service
	l       *logger.Logger
}

func New(
	pool db.Pool,
	userSvc *user.Service,
	l *logger.Logger,
) *Middleware {
	return &Middleware{
		db:      pool,
		userSvc: userSvc,
		l:       l,
	}
}

func (m *Middleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			m.l.Error("missing or invalid Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

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

func (m *Middleware) Authorize(roles []types.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("userID")
		if !exists {
			m.l.Error("user ID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		userID, ok := userIDStr.(int)
		if !ok {
			m.l.Error("user ID is not an int")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
		}

		u, err := m.userSvc.GetByID(c.Request.Context(), userID)
		if err != nil {
			m.l.Error("failed to get user by ID", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		for _, role := range roles {
			if u.Role == role {
				c.Next()
				return
			}
		}

		m.l.Error("user does not have the required role", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		c.Abort()
	}
}

// func (m *Middleware) IsAdmin() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userID, exists := c.Get("userID")
// 		if !exists {
// 			m.l.Error("user ID not found in context")
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 			c.Abort()
// 			return
// 		}

// 		u, err := m.userSvc.GetByID(c.Request.Context(), userID.(int))
// 		if err != nil {
// 			m.l.Error("failed to get user by ID", err)
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 			c.Abort()
// 			return
// 		}

// 		if u.Role != types.Admin {
// 			m.l.Error("user is not an admin", err)
// 			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
// 			c.Abort()
// 			return
// 		}

// 		c.Next()
// 	}
// }
