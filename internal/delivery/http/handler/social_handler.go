package handler

import (
	"strconv"
	"strings"

	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	socialuc "gym-pro-2026-ptit/internal/usecase/social"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"

	"github.com/gin-gonic/gin"
)

type SocialHandler struct {
	socialUC *socialuc.SocialUseCases
}

type createCommentRequest struct {
	Text            *string `json:"text"`
	Content         *string `json:"content"`
	ParentCommentID *string `json:"parent_comment_id"`
	ParentID        *string `json:"parent_id"`
	ParentCommentId *string `json:"parentCommentId"`
	ParentId        *string `json:"parentId"`
}

func NewSocialHandler(socialUC *socialuc.SocialUseCases) *SocialHandler {
	return &SocialHandler{socialUC: socialUC}
}

func (h *SocialHandler) GetFeed(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")

	result, err := h.socialUC.GetFeed(c.Request.Context(), userID, cursor, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *SocialHandler) GetPostByID(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.socialUC.GetPostByID(c.Request.Context(), userID, c.Param("postId"))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *SocialHandler) GetUserProfile(c *gin.Context) {
	currentUserID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.socialUC.GetUserProfile(c.Request.Context(), currentUserID, c.Param("userId"))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *SocialHandler) GetUserPosts(c *gin.Context) {
	currentUserID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")

	result, err := h.socialUC.GetUserPosts(c.Request.Context(), currentUserID, c.Param("userId"), cursor, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *SocialHandler) CreatePost(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input socialuc.CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	result, err := h.socialUC.CreatePost(c.Request.Context(), userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, result)
}

func (h *SocialHandler) CreateMediaSignature(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input socialuc.CreateMediaSignatureInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	result, err := h.socialUC.CreateMediaSignature(c.Request.Context(), userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *SocialHandler) ConfirmMedia(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input socialuc.ConfirmMediaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	result, err := h.socialUC.ConfirmMedia(c.Request.Context(), userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *SocialHandler) LikePost(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.socialUC.LikePost(c.Request.Context(), userID, c.Param("postId"))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *SocialHandler) UnlikePost(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.socialUC.UnlikePost(c.Request.Context(), userID, c.Param("postId"))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *SocialHandler) CreateComment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	content := firstNonEmptyString(req.Content, req.Text)
	if content == nil {
		response.ValidationError(c, "validation failed", map[string]interface{}{
			"content": "content is required",
		})
		return
	}

	input := socialuc.CreateCommentInput{
		Content:  *content,
		ParentID: firstNonEmptyString(req.ParentID, req.ParentId, req.ParentCommentID, req.ParentCommentId),
	}

	result, err := h.socialUC.CreateComment(c.Request.Context(), userID, c.Param("postId"), input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, result)
}

func firstNonEmptyString(values ...*string) *string {
	for _, value := range values {
		if value == nil {
			continue
		}
		trimmed := strings.TrimSpace(*value)
		if trimmed == "" {
			continue
		}
		return &trimmed
	}
	return nil
}

func (h *SocialHandler) GetPostComments(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")

	result, err := h.socialUC.GetPostComments(c.Request.Context(), c.Param("postId"), cursor, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *SocialHandler) GetCommentReplies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")

	result, err := h.socialUC.GetCommentReplies(c.Request.Context(), c.Param("postId"), c.Param("commentId"), cursor, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}
