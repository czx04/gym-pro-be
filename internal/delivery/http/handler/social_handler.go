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
		PublicID      string `json:"public_id"`
		PublicId      string `json:"publicId"`
		ResourceType  string `json:"resource_type"`
		ResourceType2 string `json:"resourceType"`
	} `json:"media"`
	ParentCommentID *string `json:"parent_comment_id"`
	ParentID        *string `json:"parent_id"`
	ParentCommentId *string `json:"parentCommentId"`
	ParentId        *string `json:"parentId"`
}

type createMediaSignatureRequest struct {
	ResourceType  string `json:"resource_type"`
	ResourceType2 string `json:"resourceType"`
	Folder        string `json:"folder"`
}

type confirmMediaRequest struct {
	PublicID      string `json:"public_id"`
	PublicId      string `json:"publicId"`
	SecureURL     string `json:"secure_url"`
	SecureUrl     string `json:"secureUrl"`
	ResourceType  string `json:"resource_type"`
	ResourceType2 string `json:"resourceType"`
	Bytes         int64  `json:"bytes"`
}

type reportPostRequest struct {
	Reason      string  `json:"reason"`
	Description *string `json:"description"`
}

type setFollowStateRequest struct {
	Following *bool `json:"following"`
}

type setBlockStateRequest struct {
	Blocked *bool `json:"blocked"`
}

type setPostPreferenceRequest struct {
	Like          *bool `json:"like"`
	Interested    *bool `json:"interested"`
	NotInterested *bool `json:"notInterested"`
	NotCare       *bool `json:"notCare"`
	NotCare2      *bool `json:"not_care"`
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

func (h *SocialHandler) Search(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	q := strings.TrimSpace(c.Query("q"))
	typ := strings.ToLower(strings.TrimSpace(c.DefaultQuery("type", "posts")))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")

	result, err := h.socialUC.Search(c.Request.Context(), userID, q, typ, cursor, limit)
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

func (h *SocialHandler) GetPostAttachment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	postID := c.Param("postId")
	kind := strings.ToLower(strings.TrimSpace(c.Query("kind")))
	switch kind {
	case "meal":
		result, err := h.socialUC.GetPostAttachedMealLog(c.Request.Context(), userID, postID)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.Success(c, result)
	case "workout":
		result, err := h.socialUC.GetPostAttachedWorkoutSession(c.Request.Context(), userID, postID)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.Success(c, result)
	default:
		response.Error(c, errors.BadRequest("kind must be meal or workout"))
	}
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

// UpdatePost is the normalized v2 endpoint that replaces PATCH EditPost with PUT.
func (h *SocialHandler) UpdatePost(c *gin.Context) {
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

	var req createMediaSignatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}
	input := socialuc.CreateMediaSignatureInput{
		ResourceType: strings.TrimSpace(firstNonEmpty(req.ResourceType, req.ResourceType2)),
		Folder:       strings.TrimSpace(req.Folder),
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

	var req confirmMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}
	input := socialuc.ConfirmMediaInput{
		PublicID:     strings.TrimSpace(firstNonEmpty(req.PublicID, req.PublicId)),
		SecureURL:    strings.TrimSpace(firstNonEmpty(req.SecureURL, req.SecureUrl)),
		ResourceType: strings.TrimSpace(firstNonEmpty(req.ResourceType, req.ResourceType2)),
		Bytes:        req.Bytes,
	}

	result, err := h.socialUC.ConfirmMedia(c.Request.Context(), userID, input)
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
			PublicID:     strings.TrimSpace(firstNonEmpty(media.PublicID, media.PublicId)),
			ResourceType: strings.TrimSpace(firstNonEmpty(media.ResourceType, media.ResourceType2)),
		})
	}

	result, err := h.socialUC.CreateComment(c.Request.Context(), userID, c.Param("postId"), input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, result)
}

func (h *SocialHandler) UpdateComment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	result, err := h.socialUC.UpdateComment(c.Request.Context(), userID, c.Param("postId"), c.Param("commentId"), socialuc.UpdateCommentInput{
		Content: strings.TrimSpace(req.Content),
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *SocialHandler) DeleteComment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	out, err := h.socialUC.DeleteComment(c.Request.Context(), userID, c.Param("postId"), c.Param("commentId"))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"deleted": true, "deletedByRole": out.DeletedByRole})
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

func (h *SocialHandler) SetFollowState(c *gin.Context) {
	currentUserID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req setFollowStateRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Following == nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	if *req.Following {
		result, err := h.socialUC.FollowUser(c.Request.Context(), currentUserID, c.Param("userId"))
		if err != nil {
			response.Error(c, err)
			return
		}
		response.Success(c, result)
		return
	}

	result, err := h.socialUC.UnfollowUser(c.Request.Context(), currentUserID, c.Param("userId"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

func (h *SocialHandler) SetBlockState(c *gin.Context) {
	currentUserID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req setBlockStateRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Blocked == nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	if *req.Blocked {
		result, err := h.socialUC.BlockUser(c.Request.Context(), currentUserID, c.Param("userId"))
		if err != nil {
			response.Error(c, err)
			return
		}
		response.Success(c, result)
		return
	}

	result, err := h.socialUC.UnblockUser(c.Request.Context(), currentUserID, c.Param("userId"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

func (h *SocialHandler) SetPostPreference(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req setPostPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	notInterested := req.NotInterested
	if notInterested == nil {
		notInterested = firstNonNilBool(req.NotCare, req.NotCare2)
	}

	// validate mutually exclusive preferences
	if req.Interested != nil && notInterested != nil && *req.Interested && *notInterested {
		response.Error(c, errors.Validation("interested and notInterested/notCare cannot both be true"))
		return
	}

	// Require at least one field
	if req.Like == nil && req.Interested == nil && notInterested == nil {
		response.Error(c, errors.Validation("at least one preference field must be provided"))
		return
	}

	type preferenceResponse struct {
		PostID              string `json:"post_id"`
		IsLikedByMe         *bool  `json:"is_liked_by_me,omitempty"`
		LikeCount           *int   `json:"like_count,omitempty"`
		IsInterestedByMe    *bool  `json:"is_interested_by_me,omitempty"`
		IsNotInterestedByMe *bool  `json:"is_not_interested_by_me,omitempty"`
	}

	out := preferenceResponse{PostID: c.Param("postId")}

	// Like state
	if req.Like != nil {
		if *req.Like {
			likeRes, err := h.socialUC.LikePost(c.Request.Context(), userID, c.Param("postId"))
			if err != nil {
				response.Error(c, err)
				return
			}
			out.IsLikedByMe = &likeRes.IsLikedByMe
			out.LikeCount = &likeRes.LikeCount
		} else {
			likeRes, err := h.socialUC.UnlikePost(c.Request.Context(), userID, c.Param("postId"))
			if err != nil {
				response.Error(c, err)
				return
			}
			out.IsLikedByMe = &likeRes.IsLikedByMe
			out.LikeCount = &likeRes.LikeCount
		}
	}

	// Interested / Not interested
	if req.Interested != nil {
		if *req.Interested {
			prefRes, err := h.socialUC.MarkInterested(c.Request.Context(), userID, c.Param("postId"))
			if err != nil {
				response.Error(c, err)
				return
			}
			out.IsInterestedByMe = &prefRes.IsInterestedByMe
			out.IsNotInterestedByMe = &prefRes.IsNotInterestedByMe
		} else {
			prefRes, err := h.socialUC.UnmarkInterested(c.Request.Context(), userID, c.Param("postId"))
			if err != nil {
				response.Error(c, err)
				return
			}
			out.IsInterestedByMe = &prefRes.IsInterestedByMe
			out.IsNotInterestedByMe = &prefRes.IsNotInterestedByMe
		}
	}

	if notInterested != nil {
		if *notInterested {
			prefRes, err := h.socialUC.MarkNotInterested(c.Request.Context(), userID, c.Param("postId"))
			if err != nil {
				response.Error(c, err)
				return
			}
			out.IsInterestedByMe = &prefRes.IsInterestedByMe
			out.IsNotInterestedByMe = &prefRes.IsNotInterestedByMe
		} else {
			prefRes, err := h.socialUC.UnmarkNotInterested(c.Request.Context(), userID, c.Param("postId"))
			if err != nil {
				response.Error(c, err)
				return
			}
			out.IsInterestedByMe = &prefRes.IsInterestedByMe
			out.IsNotInterestedByMe = &prefRes.IsNotInterestedByMe
		}
	}

	response.Success(c, out)
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func firstNonNilBool(values ...*bool) *bool {
	for _, value := range values {
		if value != nil {
			return value
		}
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

func (h *SocialHandler) SocialNotifications(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	typ := strings.ToLower(strings.TrimSpace(c.DefaultQuery("type", "list")))
	switch typ {
	case "unread_count":
		out, err := h.socialUC.GetNotificationsUnreadCount(c.Request.Context(), userID)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.Success(c, out)
		return
	case "list":
		filter := c.DefaultQuery("filter", "all")
		cursor := c.Query("cursor")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		out, err := h.socialUC.ListNotifications(c.Request.Context(), userID, filter, cursor, limit)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.Success(c, out)
		return
	default:
		response.Error(c, errors.BadRequest("invalid notifications type"))
	}
}

type socialNotificationsPostBody struct {
	Type string   `json:"type"`
	IDs  []string `json:"ids"`
}

func (h *SocialHandler) SocialNotificationsWrite(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req socialNotificationsPostBody
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}
	if strings.TrimSpace(strings.ToLower(req.Type)) != "mark_read" {
		response.Error(c, errors.BadRequest("invalid type"))
		return
	}

	out, err := h.socialUC.MarkNotificationsRead(c.Request.Context(), userID, req.IDs)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, out)
}
