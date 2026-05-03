package handler

import (
	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"
	"strconv"

	exerciseuc "gym-pro-2026-ptit/internal/usecase/exercise"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExerciseHandler struct {
	exerciseUC *exerciseuc.ExerciseUseCases
}

func NewExerciseHandler(exerciseUC *exerciseuc.ExerciseUseCases) *ExerciseHandler {
	return &ExerciseHandler{exerciseUC: exerciseUC}
}

func (h *ExerciseHandler) ListExercises(c *gin.Context) {
	page, pageSize := 1, 20
	if c.Query("page") != "" {
		page, _ = strconv.Atoi(c.Query("page"))
	}
	if c.Query("page_size") != "" {
		pageSize, _ = strconv.Atoi(c.Query("page_size"))
	}

	isFilter := false

	category := c.Query("category")
	muscleGroup := c.Query("muscle_group")
	equipment := c.Query("equipment")
	difficultyLevel := c.Query("difficulty_level")
	query := c.Query("query")
	if category != "" || muscleGroup != "" || equipment != "" || difficultyLevel != "" || query != "" {
		logger.Info("filtering exercises", "category", category, "muscle_group", muscleGroup, "equipment", equipment, "difficulty_level", difficultyLevel, "query", query)
		isFilter = true
	}
	if isFilter {
		exercises, total, err := h.exerciseUC.FilterExercises(c.Request.Context(), page, pageSize, category, muscleGroup, equipment, difficultyLevel, query)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.Paginated(c, exercises, page, pageSize, total)
		return
	} else {
		exercises, total, err := h.exerciseUC.ListExercises(c.Request.Context(), page, pageSize)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.Paginated(c, exercises, page, pageSize, total)
	}
}

func (h *ExerciseHandler) GetExercise(c *gin.Context) {
	exerciseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid exercise ID"))
		return
	}
	exercise, err := h.exerciseUC.GetExercise(c.Request.Context(), exerciseID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, exercise)
}

func (h *ExerciseHandler) GetExerciseStats(c *gin.Context) {
	exerciseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid exercise ID"))
		return
	}
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	stats, err := h.exerciseUC.GetExerciseStats(c.Request.Context(), userID, exerciseID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, stats)
}
