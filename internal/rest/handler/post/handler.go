package post

import (
	"context"
	"net/http"
	"strconv"
	"time"
	postService "twitter-api/internal/service/post"
	"twitter-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/uploadcare/uploadcare-go/ucare"
	"github.com/uploadcare/uploadcare-go/upload"
)

type Handler struct {
	ucareClient ucare.Client
	postSvc     *postService.Service
	l           *logger.Logger
}

func NewHandler(
	ucareClient ucare.Client,
	postSvc *postService.Service,
	l *logger.Logger,
) *Handler {
	return &Handler{
		ucareClient: ucareClient,
		postSvc:     postSvc,
		l:           l,
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

	content := c.PostForm("content")
	if content == "" {
		h.l.Error("content is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
		return
	}

	file, err := c.FormFile("file")
	var fileURL string

	if err == nil {
		src, err := file.Open()
		if err != nil {
			h.l.Error("failed to open file", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
			return
		}
		defer func() {
			if err := src.Close(); err != nil {
				h.l.Error("failed to close file", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to close file"})
				return
			}
		}()

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*60*time.Second)
		defer cancel()

		uploadSvc := upload.NewService(h.ucareClient)

		params := upload.FileParams{
			Data: src,
			Name: file.Filename,
		}

		fileID, err := uploadSvc.File(ctx, params)
		if err != nil {
			h.l.Error("failed to upload file", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
			return
		}

		fileURL = "https://ucarecdn.com/" + fileID + "/" + file.Filename
	} else {
		fileURL = ""
	}

	_, err = h.postSvc.CreatePost(c.Request.Context(), userID, content, fileURL)
	if err != nil {
		h.l.Error("failed to create post", err)
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

func getPageFromQuery(c *gin.Context) int {
	pageStr := c.Query("page")
	if pageStr == "" {
		return 1
	}

	if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
		return page
	}

	return 1
}

func (h *Handler) GetAllPosts(c *gin.Context) {
	page := getPageFromQuery(c)

	limit := 10
	offset := (page - 1) * limit

	posts, err := h.postSvc.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		h.l.Error("failed to get all posts", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) UpdateExistingPost(c *gin.Context) {
	postIDStr, ok := c.Params.Get("postID")
	if !ok {
		h.l.Error("post ID not found in context")
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.l.Error("post ID is not an integer")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "post not found"})
		return
	}

	content := c.PostForm("content")
	if content == "" {
		h.l.Error("content is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
		return
	}

	err = h.postSvc.UpdatePostContent(c.Request.Context(), postID, content)
	if err != nil {
		h.l.Error("failed to update post content", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "post updated successfully"})
		return
	}

	src, err := file.Open()
	if err != nil {
		h.l.Error("failed to open file", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer func() {
		if cerr := src.Close(); cerr != nil {
			h.l.Error("failed to close file", cerr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to close file"})
		}
	}()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*60*time.Second)
	defer cancel()

	uploadSvc := upload.NewService(h.ucareClient)

	params := upload.FileParams{
		Data: src,
		Name: file.Filename,
	}

	fileID, err := uploadSvc.File(ctx, params)
	if err != nil {
		h.l.Error("failed to upload file", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
		return
	}

	fileURL := "https://ucarecdn.com/" + fileID + "/" + file.Filename

	err = h.postSvc.UpdatePostFileURL(c.Request.Context(), postID, fileURL)
	if err != nil {
		h.l.Error("failed to update post file URL", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update post file URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post updated successfully"})
}

func (h *Handler) SoftDeleteByID(c *gin.Context) {
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

	isAuthor, err := h.postSvc.IsAuthor(c.Request.Context(), postID, userID)
	if err != nil {
		h.l.Error("failed to check if user is author", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check author"})
		return
	}

	if !isAuthor {
		h.l.Error("user is not the author of the post")
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not the author of this post"})
		return
	}

	err = h.postSvc.SoftDeletePost(c.Request.Context(), postID)
	if err != nil {
		h.l.Error("failed to soft delete post", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to soft delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post soft deleted successfully"})
}

func (h *Handler) HardDeleteByID(c *gin.Context) {
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

	post, err := h.postSvc.HardDeletePost(c.Request.Context(), postID)
	if err != nil {
		h.l.Error("failed to get post by ID", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get post"})
		return
	}

	c.JSON(http.StatusOK, post)
}
