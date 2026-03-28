package handler

import (
	"time"

	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	domainmeal "gym-pro-2026-ptit/internal/domain/meal"
	mealuc "gym-pro-2026-ptit/internal/usecase/meal"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"

	"github.com/gin-gonic/gin"
)

// Keeps domainmeal in scope for swag (@Success domainmeal.DailyNutritionTargetResponse).
var _ = domainmeal.DailyNutritionTargetResponse{}

type MealDailyHandler struct {
	mealDailyUC *mealuc.MealDailyUseCases
}

func NewMealDailyHandler(mealDailyUC *mealuc.MealDailyUseCases) *MealDailyHandler {
	return &MealDailyHandler{mealDailyUC: mealDailyUC}
}

// GetMealDailyTargetByDate godoc
// @Summary Get user nutrition target for a specific date
// @Description Get user nutrition daily target by date
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date path string true "Date in YYYY-MM-DD format"
// @Success 200 {object} response.Response{data=domainmeal.DailyNutritionTargetResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /meal-daily/date/{date} [get]
func (h *MealDailyHandler) GetMealDailyTargetByDate(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	dateStr := c.Param("date")
	reqDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.Error(c, errors.BadRequest("invalid date format, must be YYYY-MM-DD"))
		return
	}

	target, err := h.mealDailyUC.GetMealDailyByDate(c.Request.Context(), userID, reqDate)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, target)
}
