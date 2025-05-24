package post

import (
	"github.com/gin-gonic/gin"
	"net/http"
	postRepo "twitter-api/internal/repo/post"
	postService "twitter-api/internal/service/post"
	"twitter-api/pkg/logger"
)

type Handler struct {
	postSvc *postService.Service
	l       *logger.Logger
}

func NewHandler(
	postSvc *postService.Service,
	l *logger.Logger,
) *Handler {
	return &Handler{
		postSvc: postSvc,
		l:       l,
	}
}

func (h *Handler) CreateNewPost(c *gin.Context) {
	var postDTO postRepo.PostDTO
	if err := c.ShouldBindJSON(&postDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	_, err := h.postSvc.CreatePost(c.Request.Context(), &postDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "new post is created successfully"})
}

func (h *Handler) PostInfo(c *gin.Context) {
	postID, ok := c.Get("postID")
	if !ok {
		h.l.Error("post ID is not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no post created"})
		return
	}

	id, ok := postID.(int)
	if !ok {
		h.l.Error("post ID is not an integer")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no post created"})
		return
	}

	post, err := h.postSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		h.l.Error("failed to get post by ID", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get post"})
		return
	}

	c.JSON(http.StatusOK, post)
}
