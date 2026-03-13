package handler

import (
	"strconv"

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

// GetScheduledDates godoc
// @Summary Get dates in month that have scheduled workouts
// @Tags workout-sessions
// @Produce json
// @Security BearerAuth
// @Param month query int true "Month (1-12)"
// @Param year query int true "Year"
// @Success 200 {object} response.Response{data=[]string}
// @Router /workout-sessions/scheduled-dates [get]
func (h *WorkoutHandler) GetScheduledDates(c *gin.Context) {
	user, err := middleware.GetUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	month, _ := strconv.Atoi(c.DefaultQuery("month", "0"))
	year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))
	if month < 1 || month > 12 || year < 2000 || year > 2100 {
		response.Error(c, errors.BadRequest("month (1-12) and year required"))
		return
	}
	dates, err := h.workoutUC.GetScheduledDates(c.Request.Context(), user.ID, month, year)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, dates)
}

// GetSessionsByDate godoc
// @Summary Get workout sessions for a date
// @Tags workout-sessions
// @Produce json
// @Security BearerAuth
// @Param date query string true "Date YYYY-MM-DD"
// @Success 200 {object} response.Response{data=[]workout.WorkoutSession}
// @Router /workout-sessions [get]
func (h *WorkoutHandler) GetSessionsByDate(c *gin.Context) {
	user, err := middleware.GetUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	date := c.Query("date")
	if date == "" {
		response.Error(c, errors.BadRequest("date (YYYY-MM-DD) required"))
		return
	}
	list, err := h.workoutUC.GetSessionsByDate(c.Request.Context(), user.ID, date)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, list)
}

// GetSessionByID godoc
// @Summary Get session detail for tracking screen
// @Tags workout-sessions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Success 200 {object} response.Response{data=workout.WorkoutSession}
// @Router /workout-sessions/{id} [get]
func (h *WorkoutHandler) GetSessionByID(c *gin.Context) {
	user, err := middleware.GetUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	session, err := h.workoutUC.GetSessionByID(c.Request.Context(), user.ID, c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, session)
}

// CreateWorkoutSession godoc
// @Summary Create / schedule a workout session
// @Tags workout-sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body workout.CreateWorkoutSessionInput true "Body"
// @Success 201 {object} response.Response{data=workout.WorkoutSession}
// @Router /workout-sessions [post]
func (h *WorkoutHandler) CreateWorkoutSession(c *gin.Context) {
	user, err := middleware.GetUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var input workout.CreateWorkoutSessionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}
	session, err := h.workoutUC.CreateWorkoutSession(c.Request.Context(), user.ID, input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, session)
}

// UpdateWorkoutSession godoc
// @Summary Update session (e.g. status in_progress, startedAt)
// @Tags workout-sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Param body body workout.UpdateWorkoutSessionInput true "Body"
// @Success 200 {object} response.Response{data=workout.WorkoutSession}
// @Router /workout-sessions/{id} [patch]
func (h *WorkoutHandler) UpdateWorkoutSession(c *gin.Context) {
	user, err := middleware.GetUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var input workout.UpdateWorkoutSessionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}
	session, err := h.workoutUC.UpdateWorkoutSession(c.Request.Context(), user.ID, c.Param("id"), input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, session)
}

// UpdateSessionSet godoc
// @Summary Update one set (reps, weight_kg, completed) - Complete Set
// @Tags workout-sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Param setId path string true "Set ID"
// @Param body body workout.UpdateSessionSetInput true "Body"
// @Success 200 {object} response.Response
// @Router /workout-sessions/{id}/exercise-sets/{setId} [patch]
func (h *WorkoutHandler) UpdateSessionSet(c *gin.Context) {
	user, err := middleware.GetUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var input workout.UpdateSessionSetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}
	err = h.workoutUC.UpdateSessionSet(c.Request.Context(), user.ID, c.Param("id"), c.Param("setId"), input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}

// FinishWorkoutSession godoc
// @Summary Finish session (status completed, durationSecs)
// @Tags workout-sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Param body body workout.CompleteWorkoutSessionInput true "Body"
// @Success 200 {object} response.Response{data=workout.WorkoutSession}
// @Router /workout-sessions/{id}/finish [patch]
func (h *WorkoutHandler) FinishWorkoutSession(c *gin.Context) {
	user, err := middleware.GetUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var input workout.CompleteWorkoutSessionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}
	session, err := h.workoutUC.FinishWorkoutSession(c.Request.Context(), user.ID, c.Param("id"), input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, session)
}
