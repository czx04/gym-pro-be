package handler

import (
	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	"gym-pro-2026-ptit/internal/domain/meal"
	mealuc "gym-pro-2026-ptit/internal/usecase/meal"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"
	"strconv"

	"gym-pro-2026-ptit/pkg/cloudinary"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FoodHandler struct {
	foodUC *mealuc.FoodUseCases
}

func NewFoodHandler(foodUC *mealuc.FoodUseCases) *FoodHandler {
	return &FoodHandler{foodUC: foodUC}
}

// CreateFood godoc
// @Summary Create a food item
// @Description Create a new custom or system food based on user role
// @Tags foods
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Food Name"
// @Param description formData string false "Description"
// @Param brand formData string false "Brand"
// @Param barcode formData string false "Barcode"
// @Param serving_size formData number true "Serving Size"
// @Param unit formData string true "Unit"
// @Param calories formData number true "Calories"
// @Param protein_g formData number true "Protein (g)"
// @Param carbs_g formData number true "Carbs (g)"
// @Param fat_g formData number true "Fat (g)"
// @Param fiber_g formData number false "Fiber (g)"
// @Param category formData string false "Category"
// @Param image formData file false "Food Image (Max 10MB)"
// @Success 201 {object} response.Response{data=meal.Food}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /foods [post]
func (h *FoodHandler) CreateFood(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input meal.CreateFoodInput
	if err := c.ShouldBind(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request parameters"))
		return
	}

	// Handle file upload
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		if file.Size > 10*1024*1024 {
			response.Error(c, errors.BadRequest("image file size exceeds 10MB"))
			return
		}
		openedFile, err := file.Open()
		if err != nil {
			response.Error(c, errors.InternalServer("failed to open image file", err))
			return
		}
		defer openedFile.Close()

		imageURL, err := cloudinary.UploadImage(c.Request.Context(), openedFile)
		if err != nil {
			response.Error(c, errors.InternalServer("failed to upload image", err))
			return
		}
		input.ImageUrl = &imageURL
	}

	food, err := h.foodUC.CreateFood(c.Request.Context(), userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, food)
}

// GetFood godoc
// @Summary Get a food item
// @Description Retrieve a food item by ID
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Food ID"
// @Success 200 {object} response.Response{data=meal.Food}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /foods/{id} [get]
func (h *FoodHandler) GetFood(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	foodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid food ID"))
		return
	}

	food, err := h.foodUC.GetFood(c.Request.Context(), foodID, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, food)
}

// ListFoods godoc
// @Summary List food items
// @Description Retrieve food items with pagination (system foods + own custom foods)
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse{data=[]meal.Food}
// @Failure 401 {object} response.Response
// @Router /foods [get]
func (h *FoodHandler) ListFoods(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	foods, total, err := h.foodUC.ListFoods(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, foods, page, pageSize, total)
}

// SearchFoods godoc
// @Summary Search food items
// @Description Search food items by name, category, or system flag
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param query query string false "Search term"
// @Param category query string false "Category filter"
// @Param is_system query boolean false "System flag filter"
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse{data=[]meal.Food}
// @Failure 401 {object} response.Response
// @Router /foods/search [get]
func (h *FoodHandler) SearchFoods(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	filter := meal.SearchFoodsFilter{
		Page:     1,
		PageSize: 20,
	}

	if p, err := strconv.Atoi(c.Query("page")); err == nil {
		filter.Page = p
	}
	if ps, err := strconv.Atoi(c.Query("page_size")); err == nil {
		filter.PageSize = ps
	}

	if q := c.Query("query"); q != "" {
		filter.Query = &q
	}
	if cat := c.Query("category"); cat != "" {
		filter.Category = &cat
	}
	if isSysStr := c.Query("is_system"); isSysStr != "" {
		isSys, err := strconv.ParseBool(isSysStr)
		if err == nil {
			filter.IsSystem = &isSys
		}
	}

	foods, total, err := h.foodUC.SearchFoods(c.Request.Context(), userID, filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, foods, filter.Page, filter.PageSize, total)
}

// UpdateFood godoc
// @Summary Update a food item
// @Description Update a food item's details (must be owner or admin for system foods)
// @Tags foods
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Food ID"
// @Param name formData string false "Food Name"
// @Param description formData string false "Description"
// @Param brand formData string false "Brand"
// @Param barcode formData string false "Barcode"
// @Param serving_size formData number false "Serving Size"
// @Param unit formData string false "Unit"
// @Param calories formData number false "Calories"
// @Param protein_g formData number false "Protein (g)"
// @Param carbs_g formData number false "Carbs (g)"
// @Param fat_g formData number false "Fat (g)"
// @Param fiber_g formData number false "Fiber (g)"
// @Param category formData string false "Category"
// @Param image formData file false "Food Image (Max 10MB)"
// @Success 200 {object} response.Response{data=meal.Food}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /foods/{id} [put]
func (h *FoodHandler) UpdateFood(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	foodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid food ID"))
		return
	}

	var input meal.UpdateFoodInput
	if err := c.ShouldBind(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request parameters"))
		return
	}

	// Handle file upload
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		if file.Size > 10*1024*1024 {
			response.Error(c, errors.BadRequest("image file size exceeds 10MB"))
			return
		}
		openedFile, err := file.Open()
		if err != nil {
			response.Error(c, errors.InternalServer("failed to open image file", err))
			return
		}
		defer openedFile.Close()

		imageURL, err := cloudinary.UploadImage(c.Request.Context(), openedFile)
		if err != nil {
			response.Error(c, errors.InternalServer("failed to upload image", err))
			return
		}
		input.ImageUrl = &imageURL
	}

	food, err := h.foodUC.UpdateFood(c.Request.Context(), foodID, userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, food)
}

// DeleteFood godoc
// @Summary Delete a food item
// @Description Delete a custom food item (must be owner)
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Food ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /foods/{id} [delete]
func (h *FoodHandler) DeleteFood(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	foodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid food ID"))
		return
	}

	err = h.foodUC.DeleteFood(c.Request.Context(), foodID, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}
