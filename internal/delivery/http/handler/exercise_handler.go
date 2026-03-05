package handler

import (
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"

	"gym-pro-2026-ptit/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	exerciseuc "gym-pro-2026-ptit/internal/usecase/exercise"
)

type ExerciseHandler struct {
	listExercisesUC  *exerciseuc.ListExercisesUseCase
	getExerciseUC    *exerciseuc.GetExerciseUseCase
	searchExerciseUC *exerciseuc.FilterExerciseUseCase
	log              logger.Logger
}

func NewExerciseHandler(
	listExercisesUC *exerciseuc.ListExercisesUseCase,
	getExerciseUC *exerciseuc.GetExerciseUseCase,
	searchExerciseUC *exerciseuc.FilterExerciseUseCase,
	log logger.Logger,
) *ExerciseHandler {
	return &ExerciseHandler{listExercisesUC: listExercisesUC, getExerciseUC: getExerciseUC, searchExerciseUC: searchExerciseUC, log: log}
}

func (h *ExerciseHandler) ListExercises(c *gin.Context) {
	page, pageSize := c.GetInt("page"), c.GetInt("page_size")
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	isFilter := false

	category := c.Query("category")
	muscleGroup := c.Query("muscle_group")
	equipment := c.Query("equipment")
	difficultyLevel := c.Query("difficulty_level")
	query := c.Query("query")
	if category != "" || muscleGroup != "" || equipment != "" || difficultyLevel != "" || query != "" {
		h.log.Info("filtering exercises with category", zap.String("category", category))
		h.log.Info("filtering exercises with muscleGroup", zap.String("muscleGroup", muscleGroup))
		h.log.Info("filtering exercises with equipment", zap.String("equipment", equipment))
		h.log.Info("filtering exercises with difficultyLevel", zap.String("difficultyLevel", difficultyLevel))
		h.log.Info("filtering exercises with query", zap.String("query", query))
		isFilter = true
	}
	if isFilter {
		exercises, total, err := h.searchExerciseUC.Excute(c.Request.Context(), page, pageSize, category, muscleGroup, equipment, difficultyLevel, query)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.Paginated(c, exercises, page, pageSize, total)
		return
	} else {
		exercises, total, err := h.listExercisesUC.Excute(c.Request.Context(), page, pageSize)
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
	exercise, err := h.getExerciseUC.Excute(c.Request.Context(), exerciseID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, exercise)
}
