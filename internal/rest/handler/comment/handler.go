package comment

import (
	"github.com/gin-gonic/gin"
	"net/http"
	commentRepo "twitter-api/internal/repo/comment"
	commentService "twitter-api/internal/service/comment"
	"twitter-api/pkg/logger"
)

type Handler struct {
	commentSvc *commentService.Service
	l          *logger.Logger
}

func NewHandler(
	commentSvc *commentService.Service,
	l *logger.Logger,
) *Handler {
	return &Handler{
		commentSvc: commentSvc,
		l:          l,
	}
}

func (h *Handler) CreateNewComment(c *gin.Context) {
	var commentDTO commentRepo.CommentDTO
	if err := c.ShouldBindJSON(&commentDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	_, err := h.commentSvc.CreateComment(c.Request.Context(), &commentDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "new comment is added successfully"})
}

func (h *Handler) CommentInfo(c *gin.Context) {
	commentID, ok := c.Get("commentID") // should be corrected
	if !ok {
		h.l.Error("comment ID is not found in context")
		c.JSON(http.StatusNotFound, gin.H{"error": "no comment added"})
		return
	}

	id, ok := commentID.(int)
	if !ok {
		h.l.Error("post ID is not an integer")
		c.JSON(http.StatusNotFound, gin.H{"error": "no comment added"})
		return
	}

	comment, err := h.commentSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		h.l.Error("failed to get comment by ID", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get comment by ID"})
		return
	}

	c.JSON(http.StatusOK, comment)
}
