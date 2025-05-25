package comment

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	commentIDStr, ok := c.Params.Get("commentID")
	if !ok {
		h.l.Error("comment ID is not found in context")
		c.JSON(http.StatusNotFound, gin.H{"error": "no comment added"})
		return
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if !ok {
		h.l.Error("comment ID is not an integer")
		c.JSON(http.StatusNotFound, gin.H{"error": "no comment added"})
		return
	}

	comment, err := h.commentSvc.GetByID(c.Request.Context(), commentID)
	if err != nil {
		h.l.Error("failed to get comment by ID", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get comment by ID"})
		return
	}

	c.JSON(http.StatusOK, comment)
}
func (h *Handler) UpdateExistingComment(c *gin.Context) {
	commentIDStr, ok := c.Params.Get("commentID")
	if !ok {
		h.l.Error("comment with this ID not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	//should be corrected
	content := c.PostForm("content")
	if content == "" {
		h.l.Error("content is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
		return
	}
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		h.l.Error("comment with this ID not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	comment, err := h.commentSvc.UpdateComment(c.Request.Context(), commentID, content)
	if err != nil {
		h.l.Error("failed to update comment by ID", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get post"})
		return
	}

	c.JSON(http.StatusOK, comment)
}
func (h *Handler) GetAllCommentsTOPosts(c *gin.Context) {
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

	comments, err := h.commentSvc.GetALlPostComments(c.Request.Context(), postID)
	if err != nil {
		h.l.Error("failed to get post comment", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to post comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *Handler) GetAllComments(c *gin.Context) {
	comments, err := h.commentSvc.GetALlCommentsByAdmin(c.Request.Context())
	if err != nil {
		h.l.Error("failed to get all comments", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}
func (h *Handler) SoftDeleteComment(c *gin.Context) {
	commentIDStr, ok := c.Params.Get("commentID")
	if !ok {
		h.l.Error("comment with this ID not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		h.l.Error("comment with this ID not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	comment, err := h.commentSvc.SoftDeleteComment(c.Request.Context(), commentID)
	if err != nil {
		h.l.Error("failed to get post comment", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to post comments"})
		return
	}

	c.JSON(http.StatusOK, comment)
}
func (h *Handler) HardDeleteComment(c *gin.Context) {
	commentIDStr, ok := c.Params.Get("commentID")
	if !ok {
		h.l.Error("comment with this ID not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		h.l.Error("comment with this ID not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	err = h.commentSvc.HardDeleteComment(c.Request.Context(), commentID)
	if err != nil {
		h.l.Error("failed to get post comment", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to post comments"})
		return
	}

	c.JSON(http.StatusOK, nil)
}
