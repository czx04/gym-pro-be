package handler

import (
	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	"gym-pro-2026-ptit/internal/domain/admin"
	adminuc "gym-pro-2026-ptit/internal/usecase/admin"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	adminUC *adminuc.AdminUseCases
}

func NewAdminHandler(adminUC *adminuc.AdminUseCases) *AdminHandler {
	return &AdminHandler{adminUC: adminUC}
}

// GetOverviewStats godoc
// @Summary Get admin overview stats
// @Description Returns summary statistics for admin dashboard (total users, exercises, foods)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=admin.OverviewStats}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /admin/stats [get]
func (h *AdminHandler) GetOverviewStats(c *gin.Context) {
	stats, err := h.adminUC.GetOverviewStats(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, stats)
}

// ListUsers godoc
// @Summary List all users (admin)
// @Description Returns a paginated list of all users with optional filters
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param query query string false "Search by name or email"
// @Param gender query string false "Filter by gender"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} response.PaginatedResponse{data=[]admin.UserSummary}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page := c.GetInt("page")
	pageSize := c.GetInt("page_size")
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	filter := admin.ListUsersFilter{
		Page:     page,
		PageSize: pageSize,
	}
	if q := c.Query("query"); q != "" {
		filter.Query = &q
	}
	if g := c.Query("gender"); g != "" {
		filter.Gender = &g
	}
	if v := c.Query("is_active"); v != "" {
		isActive := v == "true"
		filter.IsActive = &isActive
	}

	users, total, err := h.adminUC.ListUsers(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Paginated(c, users, page, pageSize, total)
}

// GetUser godoc
// @Summary Get user detail (admin)
// @Description Returns detailed information of a user by ID
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} response.Response{data=admin.UserDetail}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/users/{id} [get]
func (h *AdminHandler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid user ID"))
		return
	}

	u, err := h.adminUC.GetUser(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, u)
}

// UpdateUserStatus godoc
// @Summary Update user active status (admin)
// @Description Enable or disable a user account
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Param request body admin.UpdateUserStatusInput true "Status update"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/users/{id}/status [patch]
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid user ID"))
		return
	}

	var input admin.UpdateUserStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	if err := h.adminUC.UpdateUserStatus(c.Request.Context(), userID, input); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"message": "user status updated"})
}

// DeleteUser godoc
// @Summary Delete a user (admin)
// @Description Permanently delete a user account
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Success 204 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid user ID"))
		return
	}

	if err := h.adminUC.DeleteUser(c.Request.Context(), userID); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

// ListExercises godoc
// @Summary List all exercises (admin)
// @Description Returns a paginated list of all exercises with optional filters
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param query query string false "Search by name"
// @Param category query string false "Filter by category (cardio, strength, flexibility, stretching)"
// @Param muscle_group query string false "Filter by muscle group"
// @Param difficulty_level query string false "Filter by difficulty (beginner, intermediate, advanced)"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} response.PaginatedResponse{data=[]admin.AdminExercise}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /admin/exercises [get]
func (h *AdminHandler) ListExercises(c *gin.Context) {
	page := c.GetInt("page")
	pageSize := c.GetInt("page_size")
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	filter := admin.ListExercisesFilter{
		Page:     page,
		PageSize: pageSize,
	}
	if q := c.Query("query"); q != "" {
		filter.Query = &q
	}
	if v := c.Query("category"); v != "" {
		filter.Category = &v
	}
	if v := c.Query("muscle_group"); v != "" {
		filter.MuscleGroup = &v
	}
	if v := c.Query("difficulty_level"); v != "" {
		filter.DifficultyLevel = &v
	}
	if v := c.Query("is_active"); v != "" {
		isActive := v == "true"
		filter.IsActive = &isActive
	}

	exercises, total, err := h.adminUC.ListExercises(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Paginated(c, exercises, page, pageSize, total)
}

// GetExercise godoc
// @Summary Get exercise detail (admin)
// @Description Returns detailed information of an exercise by ID
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "Exercise ID (UUID)"
// @Success 200 {object} response.Response{data=admin.AdminExercise}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/exercises/{id} [get]
func (h *AdminHandler) GetExercise(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid exercise ID"))
		return
	}

	e, err := h.adminUC.GetExercise(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, e)
}

// CreateExercise godoc
// @Summary Create a new exercise (admin)
// @Description Admin creates a new exercise in the system
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body admin.CreateExerciseInput true "Exercise data"
// @Success 201 {object} response.Response{data=admin.AdminExercise}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /admin/exercises [post]
func (h *AdminHandler) CreateExercise(c *gin.Context) {
	adminID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input admin.CreateExerciseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	e, err := h.adminUC.CreateExercise(c.Request.Context(), adminID, input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, e)
}

// UpdateExercise godoc
// @Summary Update an exercise (admin)
// @Description Admin updates an existing exercise
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Exercise ID (UUID)"
// @Param request body admin.UpdateExerciseInput true "Exercise update data"
// @Success 200 {object} response.Response{data=admin.AdminExercise}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/exercises/{id} [put]
func (h *AdminHandler) UpdateExercise(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid exercise ID"))
		return
	}

	var input admin.UpdateExerciseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	e, err := h.adminUC.UpdateExercise(c.Request.Context(), id, input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, e)
}

// DeleteExercise godoc
// @Summary Delete an exercise (admin)
// @Description Admin deletes an exercise from the system
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "Exercise ID (UUID)"
// @Success 204 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/exercises/{id} [delete]
func (h *AdminHandler) DeleteExercise(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid exercise ID"))
		return
	}

	if err := h.adminUC.DeleteExercise(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

// ListFoods godoc
// @Summary List all foods (admin)
// @Description Returns a paginated list of all foods with optional filters
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param query query string false "Search by name"
// @Param category query string false "Filter by category"
// @Param is_system query bool false "Filter by system food"
// @Success 200 {object} response.PaginatedResponse{data=[]admin.AdminFood}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /admin/foods [get]
func (h *AdminHandler) ListFoods(c *gin.Context) {
	page := c.GetInt("page")
	pageSize := c.GetInt("page_size")
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	filter := admin.ListFoodsFilter{
		Page:     page,
		PageSize: pageSize,
	}
	if q := c.Query("query"); q != "" {
		filter.Query = &q
	}
	if v := c.Query("category"); v != "" {
		filter.Category = &v
	}
	if v := c.Query("is_system"); v != "" {
		isSystem := v == "true"
		filter.IsSystem = &isSystem
	}

	foods, total, err := h.adminUC.ListFoods(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Paginated(c, foods, page, pageSize, total)
}

// GetFood godoc
// @Summary Get food detail (admin)
// @Description Returns detailed information of a food item by ID
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "Food ID (UUID)"
// @Success 200 {object} response.Response{data=admin.AdminFood}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/foods/{id} [get]
func (h *AdminHandler) GetFood(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid food ID"))
		return
	}

	f, err := h.adminUC.GetFood(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, f)
}

// CreateSystemFood godoc
// @Summary Create a system food (admin)
// @Description Admin creates a new system-level food item
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body admin.CreateSystemFoodInput true "Food data"
// @Success 201 {object} response.Response{data=admin.AdminFood}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /admin/foods [post]
func (h *AdminHandler) CreateSystemFood(c *gin.Context) {
	adminID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input admin.CreateSystemFoodInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	f, err := h.adminUC.CreateSystemFood(c.Request.Context(), adminID, input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, f)
}

// UpdateFood godoc
// @Summary Update a food item (admin)
// @Description Admin updates an existing food item
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Food ID (UUID)"
// @Param request body admin.AdminUpdateFoodInput true "Food update data"
// @Success 200 {object} response.Response{data=admin.AdminFood}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/foods/{id} [put]
func (h *AdminHandler) UpdateFood(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid food ID"))
		return
	}

	var input admin.AdminUpdateFoodInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	f, err := h.adminUC.UpdateFood(c.Request.Context(), id, input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, f)
}

// DeleteFood godoc
// @Summary Delete a food item (admin)
// @Description Admin deletes a food item from the system
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "Food ID (UUID)"
// @Success 204 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/foods/{id} [delete]
func (h *AdminHandler) DeleteFood(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid food ID"))
		return
	}

	if err := h.adminUC.DeleteFood(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}
