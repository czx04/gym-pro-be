package handler

import (
	"strconv"
	"strings"

	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	socialdomain "gym-pro-2026-ptit/internal/domain/social"
	socialuc "gym-pro-2026-ptit/internal/usecase/social"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"

	"github.com/gin-gonic/gin"
)

type SocialHandler struct {
	socialUC *socialuc.SocialUseCases
}

type createCommentRequest struct {
	Text    *string `json:"text"`
	Content *string `json:"content"`
	Media   []struct {
		PublicID     string `json:"public_id"`
		ResourceType string `json:"resource_type"`
	} `json:"media"`
	ParentCommentID *string `json:"parent_comment_id"`
	ParentID        *string `json:"parent_id"`
	ParentCommentId *string `json:"parentCommentId"`
	ParentId        *string `json:"parentId"`
}

type reportPostRequest struct {
	Reason      string  `json:"reason"`
	Description *string `json:"description"`
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

func (h *SocialHandler) FollowUser(c *gin.Context) {
	currentUserID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.socialUC.FollowUser(c.Request.Context(), currentUserID, c.Param("userId"))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *SocialHandler) UnfollowUser(c *gin.Context) {
	currentUserID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.socialUC.UnfollowUser(c.Request.Context(), currentUserID, c.Param("userId"))
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

func (h *SocialHandler) EditPost(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input socialuc.UpdatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	result, err := h.socialUC.EditPost(c.Request.Context(), userID, c.Param("postId"), input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *SocialHandler) DeletePost(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	if err := h.socialUC.DeletePost(c.Request.Context(), userID, c.Param("postId")); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"deleted": true})
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

func (h *SocialHandler) MarkInterested(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.socialUC.MarkInterested(c.Request.Context(), userID, c.Param("postId"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

func (h *SocialHandler) UnmarkInterested(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.socialUC.UnmarkInterested(c.Request.Context(), userID, c.Param("postId"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

func (h *SocialHandler) MarkNotInterested(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.socialUC.MarkNotInterested(c.Request.Context(), userID, c.Param("postId"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

func (h *SocialHandler) UnmarkNotInterested(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.socialUC.UnmarkNotInterested(c.Request.Context(), userID, c.Param("postId"))
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
	if content == nil && len(req.Media) == 0 {
		response.ValidationError(c, "validation failed", map[string]interface{}{
			"content": "content or media is required",
		})
		return
	}

	input := socialuc.CreateCommentInput{
		Content:  content,
		Media:    make([]socialdomain.CreatePostMediaInput, 0, len(req.Media)),
		ParentID: firstNonEmptyString(req.ParentID, req.ParentId, req.ParentCommentID, req.ParentCommentId),
	}
	for _, media := range req.Media {
		input.Media = append(input.Media, socialdomain.CreatePostMediaInput{
			PublicID:     media.PublicID,
			ResourceType: media.ResourceType,
		})
	}

	result, err := h.socialUC.CreateComment(c.Request.Context(), userID, c.Param("postId"), input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, result)
}

func (h *SocialHandler) DeleteComment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	if err := h.socialUC.DeleteComment(c.Request.Context(), userID, c.Param("postId"), c.Param("commentId")); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

func (h *SocialHandler) ReportPost(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req reportPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}
	result, err := h.socialUC.ReportPost(c.Request.Context(), userID, c.Param("postId"), socialuc.ReportPostInput{
		Reason:      req.Reason,
		Description: req.Description,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

func (h *SocialHandler) BlockUser(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.socialUC.BlockUser(c.Request.Context(), userID, c.Param("userId"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

func (h *SocialHandler) UnblockUser(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.socialUC.UnblockUser(c.Request.Context(), userID, c.Param("userId"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
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
