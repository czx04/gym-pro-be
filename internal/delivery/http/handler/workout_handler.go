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
	createPlanUC    *workoutuc.CreateWorkoutPlanUseCase
	addExerciseUC   *workoutuc.AddExerciseToWorkoutUseCase
	// TODO: Add more use cases as you implement them
}

// NewWorkoutHandler creates a new workout handler
func NewWorkoutHandler(
	createPlanUC *workoutuc.CreateWorkoutPlanUseCase,
	addExerciseUC *workoutuc.AddExerciseToWorkoutUseCase,
) *WorkoutHandler {
	return &WorkoutHandler{
		createPlanUC:  createPlanUC,
		addExerciseUC: addExerciseUC,
	}
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
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input workout.CreateWorkoutPlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	plan, err := h.createPlanUC.Execute(c.Request.Context(), userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, plan, "Workout plan created successfully")
}

// AddExerciseToWorkout godoc
// @Summary Add exercise to workout plan
// @Description Add an exercise to a workout plan with configuration
// @Tags workout-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Workout Plan ID" format(uuid)
// @Param request body workout.AddExerciseToWorkoutInput true "Exercise configuration"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /workout-plans/{id}/exercises [post]
func (h *WorkoutHandler) AddExerciseToWorkout(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid workout plan ID"))
		return
	}

	var input workout.AddExerciseToWorkoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	if err := h.addExerciseUC.Execute(c.Request.Context(), planID, input); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil, "Exercise added to workout plan")
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
	// TODO: Implement list workout plans
	response.Error(c, errors.InternalServer("not implemented", nil))
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
	// TODO: Implement get workout plan
	response.Error(c, errors.InternalServer("not implemented", nil))
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
	// TODO: Implement update workout plan
	response.Error(c, errors.InternalServer("not implemented", nil))
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
	// TODO: Implement delete workout plan
	response.Error(c, errors.InternalServer("not implemented", nil))
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
