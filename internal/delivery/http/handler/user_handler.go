package handler

import (
	"time"

	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	mealdomain "gym-pro-2026-ptit/internal/domain/meal"
	domainuser "gym-pro-2026-ptit/internal/domain/user"
	mealuc "gym-pro-2026-ptit/internal/usecase/meal"
	useruc "gym-pro-2026-ptit/internal/usecase/user"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"

	"github.com/gin-gonic/gin"
)

// Keeps domainuser in scope for swag (@Success domainuser.UserNutritionTarget).
var _ = domainuser.UserNutritionTarget{}

// Keeps domainuser.WeightHistoryPoint in scope for swag.
var _ = domainuser.WeightHistoryPoint{}

var _ = mealdomain.RegisterPushTokenInput{}

type UserHandler struct {
	userUC       *useruc.UserUseCases
	mealStreakUC *mealuc.MealStreakUseCases
	pushTokenUC  *mealuc.PushTokenUseCases
}

func NewUserHandler(
	userUC *useruc.UserUseCases,
	mealStreakUC *mealuc.MealStreakUseCases,
	pushTokenUC *mealuc.PushTokenUseCases,
) *UserHandler {
	return &UserHandler{
		userUC:       userUC,
		mealStreakUC: mealStreakUC,
		pushTokenUC:  pushTokenUC,
	}
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
// GetMyWeightHistory godoc
// @Summary List weight history for chart
// @Description Latest weight per day/week/month bucket in the given timezone
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string true "Start (RFC3339)"
// @Param to query string true "End (RFC3339)"
// @Param granularity query string true "day, week, or month"
// @Param timezone query string false "IANA timezone (default UTC)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /users/me/weight-history [get]
func (h *UserHandler) GetMyWeightHistory(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr == "" || toStr == "" {
		response.Error(c, errors.BadRequest("from and to query parameters are required (RFC3339)"))
		return
	}
	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		response.Error(c, errors.BadRequest("invalid from datetime"))
		return
	}
	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		response.Error(c, errors.BadRequest("invalid to datetime"))
		return
	}

	tz := c.Query("timezone")
	granularity := domainuser.WeightHistoryGranularity(c.Query("granularity"))
	if granularity == "" {
		response.Error(c, errors.BadRequest("granularity is required (day, week, month)"))
		return
	}

	points, err := h.userUC.ListMyWeightHistory(c.Request.Context(), userID, from, to, tz, granularity)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"points": points})
}

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

// GetMealStreak returns current and longest meal logging streak (recomputed from logs).
func (h *UserHandler) GetMealStreak(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	s, err := h.mealStreakUC.RecalculateAndPersist(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, s)
}

// RegisterPushToken stores an Expo push token for daily meal reminders.
func (h *UserHandler) RegisterPushToken(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var body mealdomain.RegisterPushTokenInput
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}
	if err := h.pushTokenUC.Register(c.Request.Context(), userID, body); err != nil {
		response.Error(c, err)
		return
	}
	c.Status(204)
}

// DeletePushToken removes a registered Expo push token.
func (h *UserHandler) DeletePushToken(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var body struct {
		ExpoPushToken string `json:"expo_push_token"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}
	if err := h.pushTokenUC.Delete(c.Request.Context(), userID, body.ExpoPushToken); err != nil {
		response.Error(c, err)
		return
	}
	c.Status(204)
}

// UpsertMyDailySteps godoc
// @Summary Upsert daily step count (Apple Health)
// @Description Upsert total steps for a local calendar day
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body useruc.UpsertDailyStepsInput true "Daily steps payload"
// @Success 204 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /users/me/steps/daily [post]
func (h *UserHandler) UpsertMyDailySteps(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var body useruc.UpsertDailyStepsInput
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	if err := h.userUC.UpsertMyDailySteps(c.Request.Context(), userID, body); err != nil {
		response.Error(c, err)
		return
	}

	c.Status(204)
}

// ListMyDailySteps godoc
// @Summary List daily step totals
// @Description List totals in [from,to] (YYYY-MM-DD)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string true "YYYY-MM-DD"
// @Param to query string true "YYYY-MM-DD"
// @Param source query string false "apple_health (default)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /users/me/steps/daily [get]
func (h *UserHandler) ListMyDailySteps(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr == "" || toStr == "" {
		response.Error(c, errors.BadRequest("from and to query parameters are required (YYYY-MM-DD)"))
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		response.Error(c, errors.BadRequest("invalid from date, expected YYYY-MM-DD"))
		return
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		response.Error(c, errors.BadRequest("invalid to date, expected YYYY-MM-DD"))
		return
	}

	source := c.Query("source")
	points, err := h.userUC.ListMyDailySteps(c.Request.Context(), userID, from, to, source)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"points": points})
}
