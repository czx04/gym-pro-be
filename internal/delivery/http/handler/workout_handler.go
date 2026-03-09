package handler

import (
	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	"gym-pro-2026-ptit/internal/domain/workout"
	workoutuc "gym-pro-2026-ptit/internal/usecase/workout"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WorkoutHandler handles workout-related requests
type WorkoutHandler struct {
	workoutUC *workoutuc.WorkoutUseCases
}

// NewWorkoutHandler creates a new workout handler
func NewWorkoutHandler(workoutUC *workoutuc.WorkoutUseCases) *WorkoutHandler {
	return &WorkoutHandler{workoutUC: workoutUC}
}

// CreateWorkoutPlan godoc
// @Summary Create a workout plan
// @Description Create a new workout plan for the authenticated user
// @Tags workout-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body workout.CreateWorkoutPlanInput true "Workout plan data"
// @Success 201 {object} response.Response{data=workout.WorkoutPlan}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /workout-plans [post]
func (h *WorkoutHandler) CreateWorkoutPlan(c *gin.Context) {
	user, err := middleware.GetUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input workout.CreateWorkoutPlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	plan, err := h.workoutUC.CreateWorkoutPlan(c.Request.Context(), user, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, plan)
}

// ListWorkoutPlans godoc
// @Summary List workout plans
// @Description Get a paginated list of user's workout plans
// @Tags workout-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} response.PaginatedResponse{data=[]workout.WorkoutPlan}
// @Failure 401 {object} response.Response
// @Router /workout-plans [get]
func (h *WorkoutHandler) ListWorkoutPlans(c *gin.Context) {
	user, err := middleware.GetUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	page, pageSize := c.GetInt("page"), c.GetInt("page_size")
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	plans, total, err := h.workoutUC.ListWorkoutPlans(c.Request.Context(), *user, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Paginated(c, plans, page, pageSize, total)
}

// GetWorkoutPlan godoc
// @Summary Get workout plan
// @Description Get workout plan details with exercises
// @Tags workout-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Workout Plan ID" format(uuid)
// @Success 200 {object} response.Response{data=workout.WorkoutPlan}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /workout-plans/{id} [get]
func (h *WorkoutHandler) GetWorkoutPlan(c *gin.Context) {
	user, err := middleware.GetUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	plan, err := h.workoutUC.GetWorkoutPlan(c.Request.Context(), user.ID, c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, plan)
}

// UpdateWorkoutPlan godoc
// @Summary Update workout plan
// @Description Update workout plan details
// @Tags workout-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Workout Plan ID" format(uuid)
// @Param request body workout.UpdateWorkoutPlanInput true "Update data"
// @Success 200 {object} response.Response{data=workout.WorkoutPlan}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /workout-plans/{id} [put]
func (h *WorkoutHandler) UpdateWorkoutPlan(c *gin.Context) {
	user, err := middleware.GetUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	planID := c.Param("id")
	isUpdateExercises := c.Query("is_update_exercises") == "true"

	uuidPlanID, err := uuid.Parse(planID)
	if err != nil {
		response.Error(c, errors.BadRequest("invalid plan ID"))
		return
	}
	var input workout.UpdateWorkoutPlanInput
	input.ID = uuidPlanID
	input.IsUpdateExercises = isUpdateExercises
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}
	plan, err := h.workoutUC.UpdateWorkoutPlan(c.Request.Context(), user.ID, input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, plan)
}

// DeleteWorkoutPlan godoc
// @Summary Delete workout plan
// @Description Delete a workout plan
// @Tags workout-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Workout Plan ID" format(uuid)
// @Success 204 "No Content"
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /workout-plans/{id} [delete]
func (h *WorkoutHandler) DeleteWorkoutPlan(c *gin.Context) {
	user, err := middleware.GetUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	err = h.workoutUC.DeleteWorkoutPlan(c.Request.Context(), user.ID, c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// TODO: Implement more handlers:
// Exercise handlers:
// - ListExercises
// - GetExercise
// - SearchExercises
//
// Schedule handlers:
// - ScheduleWorkout
// - BulkScheduleWorkout
// - ListSchedules
// - GetCalendarView
// - UpdateSchedule
// - DeleteSchedule
//
// Session handlers:
// - StartWorkoutSession
// - LogExerciseSet
// - CompleteSession
// - GetSessionHistory
// - GetSessionDetails
// - GetWorkoutStats
