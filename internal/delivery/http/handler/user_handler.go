package handler

import (
	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	domainuser "gym-pro-2026-ptit/internal/domain/user"
	useruc "gym-pro-2026-ptit/internal/usecase/user"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"

	"github.com/gin-gonic/gin"
)

// Keeps domainuser in scope for swag (@Success domainuser.UserNutritionTarget).
var _ = domainuser.UserNutritionTarget{}

type UserHandler struct {
	userUC *useruc.UserUseCases
}

func NewUserHandler(userUC *useruc.UserUseCases) *UserHandler {
	return &UserHandler{userUC: userUC}
}

// GetUserNutritionTarget godoc
// @Summary Get user nutrition target
// @Description Get user nutrition target
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=domainuser.UserNutritionTarget}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/nutrition-target [get]
func (h *UserHandler) GetUserNutritionTarget(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	target, err := h.userUC.GetUserNutritionTarget(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, target)
}

// UpdateUserNutritionTarget godoc
// @Summary Update user nutrition target
// @Description Update user nutrition target
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body useruc.UpdateUserNutritionTargetInput true "Update user nutrition target"
// @Success 200 {object} response.Response{data=domainuser.UserNutritionTarget}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /users/nutrition-target [put]
func (h *UserHandler) UpdateUserNutritionTarget(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input useruc.UpdateUserNutritionTargetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	target, err := h.userUC.UpdateUserNutritionTarget(c.Request.Context(), userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, target)
}
