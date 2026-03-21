package handler

import (
	"time"

	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	"gym-pro-2026-ptit/internal/domain/meal"
	mealuc "gym-pro-2026-ptit/internal/usecase/meal"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MealLogHandler handles all meal-log HTTP endpoints.
type MealLogHandler struct {
	mealLogUC *mealuc.MealLogUseCases
}

func NewMealLogHandler(mealLogUC *mealuc.MealLogUseCases) *MealLogHandler {
	return &MealLogHandler{mealLogUC: mealLogUC}
}

// CreateMealLog godoc
// @Summary Create a meal log
// @Description Create a new meal log entry. You can include initial items (food or recipe) in the same request.
// @Tags meal-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body meal.CreateMealLogInput true "Meal log details"
// @Success 201 {object} response.Response{data=meal.MealLog}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /meal-logs [post]
func (h *MealLogHandler) CreateMealLog(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input meal.CreateMealLogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	log, err := h.mealLogUC.CreateMealLog(c.Request.Context(), userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, log)
}

// GetMealLog godoc
// @Summary Get meal log detail
// @Description Retrieve a single meal log entry with all food/recipe items
// @Tags meal-logs
// @Produce json
// @Security BearerAuth
// @Param id path string true "Meal Log ID"
// @Success 200 {object} response.Response{data=meal.MealLog}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /meal-logs/{id} [get]
func (h *MealLogHandler) GetMealLog(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	logID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid meal log ID"))
		return
	}

	log, err := h.mealLogUC.GetMealLog(c.Request.Context(), logID, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, log)
}

// GetMealLogsByDate godoc
// @Summary Get meal logs by date
// @Description Retrieve all meal logs for a specific date (YYYY-MM-DD) along with the daily nutrition summary
// @Tags meal-logs
// @Produce json
// @Security BearerAuth
// @Param date path string true "Date in YYYY-MM-DD format"
// @Success 200 {object} response.Response{data=meal.DailyMealResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /meal-logs/date/{date} [get]
func (h *MealLogHandler) GetMealLogsByDate(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.Error(c, errors.BadRequest("invalid date format, expected YYYY-MM-DD"))
		return
	}

	result, err := h.mealLogUC.GetMealLogsByDate(c.Request.Context(), userID, date)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// UpdateMealLog godoc
// @Summary Update a meal log
// @Description Update notes, mood, energy level and/or replace all items of a meal log
// @Tags meal-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Meal Log ID"
// @Param body body meal.UpdateMealLogInput true "Update data"
// @Success 200 {object} response.Response{data=meal.MealLog}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /meal-logs/{id} [put]
func (h *MealLogHandler) UpdateMealLog(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	logID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid meal log ID"))
		return
	}

	var input meal.UpdateMealLogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	log, err := h.mealLogUC.UpdateMealLog(c.Request.Context(), logID, userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, log)
}

// DeleteMealLog godoc
// @Summary Delete a meal log
// @Description Delete a meal log and all its items
// @Tags meal-logs
// @Produce json
// @Security BearerAuth
// @Param id path string true "Meal Log ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /meal-logs/{id} [delete]
func (h *MealLogHandler) DeleteMealLog(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	logID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid meal log ID"))
		return
	}

	if err := h.mealLogUC.DeleteMealLog(c.Request.Context(), logID, userID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// GetNutritionStats godoc
// @Summary Get nutrition statistics
// @Description Get nutrition statistics for a date period
// @Tags meal-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} response.Response{data=meal.NutritionStats}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /meal-logs/stats [get]
func (h *MealLogHandler) GetNutritionStats(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input meal.GetNutritionStatsRequest
	if err := c.ShouldBindQuery(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid query parameters"))
		return
	}

	result, err := h.mealLogUC.GetNutritionStats(c.Request.Context(), userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}
