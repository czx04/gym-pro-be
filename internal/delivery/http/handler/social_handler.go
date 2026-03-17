package handler

import (
	"strconv"

	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	socialuc "gym-pro-2026-ptit/internal/usecase/social"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"

	"github.com/gin-gonic/gin"
)

type SocialHandler struct {
	socialUC *socialuc.SocialUseCases
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
