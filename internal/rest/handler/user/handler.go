package user

import (
	"fmt"
	"net/http"
	userRepo "twitter-api/internal/repo/user"
	tokenService "twitter-api/internal/service/token"
	userService "twitter-api/internal/service/user"
	"twitter-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userSvc  *userService.Service
	tokenSvc *tokenService.Service
	l        *logger.Logger
}

func NewHandler(
	userSvc *userService.Service,
	tokenSvc *tokenService.Service,
	l *logger.Logger,
) *Handler {
	return &Handler{
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
		l:        l,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var userDTO userRepo.RegisterUserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	_, err := h.userSvc.CreateUser(c.Request.Context(), &userDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var userDTO userRepo.LoginUserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	u, err := h.userSvc.ValidateUser(c.Request.Context(), &userDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, refreshToken, err := h.tokenSvc.CreateTokens(c.Request.Context(), u.ID, u.Role)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) Profile(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		h.l.Error("user ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, ok := userID.(int)
	if !ok {
		h.l.Error("user ID is not an integer")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		h.l.Error("failed to get user by ID", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
