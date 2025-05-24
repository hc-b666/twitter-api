package post

import (
	"net/http"
	"strconv"
	postRepo "twitter-api/internal/repo/post"
	postService "twitter-api/internal/service/post"
	"twitter-api/pkg/logger"

	"github.com/gin-gonic/gin"
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
	userIDStr, ok := c.Get("userID")
	if !ok {
		h.l.Error("user ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDStr.(int)
	if !ok {
		h.l.Error("user ID is not an integer")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var postDTO postRepo.PostDTO
	if err := c.ShouldBindJSON(&postDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	_, err := h.postSvc.CreatePost(c.Request.Context(), userID, &postDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "new post is created successfully"})
}

func (h *Handler) GetPostByID(c *gin.Context) {
	postIDStr, ok := c.Params.Get("postID")
	if !ok {
		h.l.Error("post with this ID not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.l.Error("post with this ID not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	post, err := h.postSvc.GetByID(c.Request.Context(), postID)
	if err != nil {
		h.l.Error("failed to get post by ID", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) GetUserPosts(c *gin.Context) {
	userIDStr, ok := c.Params.Get("userID")
	if !ok {
		h.l.Error("user with this ID not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.l.Error("user with this ID not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	post, err := h.postSvc.GetUserPosts(c.Request.Context(), userID)
	if err != nil {
		h.l.Error("failed to get user posts", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user posts"})
		return
	}

	c.JSON(http.StatusOK, post)
}
