package token

import (
	"net/http"
	tokenRepository "twitter-api/internal/repo/token"
	tokenService "twitter-api/internal/service/token"
	"twitter-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	tokenSvc *tokenService.Service
}

func NewHandler(
	tokenSvc *tokenService.Service,
) *Handler {
	return &Handler{
		tokenSvc: tokenSvc,
	}
}

func (h *Handler) Refresh(c *gin.Context) {
	var tkn tokenRepository.Token
	if err := c.ShouldBindJSON(&tkn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.tokenSvc.GetByToken(c.Request.Context(), tkn.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	newAccessToken, err := utils.CreateAccessToken(token.UserID, token.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
