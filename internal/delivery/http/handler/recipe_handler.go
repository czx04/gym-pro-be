package handler

import (
	"strconv"

	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	"gym-pro-2026-ptit/internal/domain/meal"
	mealuc "gym-pro-2026-ptit/internal/usecase/meal"
	"gym-pro-2026-ptit/pkg/cloudinary"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecipeHandler struct {
	recipeUC *mealuc.RecipeUseCases
}

func NewRecipeHandler(recipeUC *mealuc.RecipeUseCases) *RecipeHandler {
	return &RecipeHandler{recipeUC: recipeUC}
}

// CreateRecipe godoc
// @Summary Create a recipe
// @Description Create a new recipe with optional foods and image
// @Tags recipes
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Recipe Name"
// @Param description formData string false "Description"
// @Param prep_time_mins formData int false "Prep Time (mins)"
// @Param cook_time_mins formData int false "Cook Time (mins)"
// @Param servings formData int true "Servings"
// @Param instructions formData string false "Instructions"
// @Param is_public formData boolean false "Is Public"
// @Param visibility formData string false "Visibility (public, private, friends)"
// @Param foods formData string false "Foods JSON string array: [{\"food_id\":\"uuid\", \"quantity\":100}]"
// @Param image formData file false "Recipe Image (Max 10MB)"
// @Success 201 {object} response.Response{data=meal.Recipe}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /recipes [post]
func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input meal.CreateRecipeInput
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
		input.ImageURL = &imageURL
	}

	recipe, err := h.recipeUC.CreateRecipe(c.Request.Context(), userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, recipe)
}

// GetRecipe godoc
// @Summary Get a recipe
// @Description Retrieve a recipe by ID with dynamically calculated nutrition
// @Tags recipes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Recipe ID"
// @Success 200 {object} response.Response{data=meal.Recipe}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /recipes/{id} [get]
func (h *RecipeHandler) GetRecipe(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	recipeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid recipe ID"))
		return
	}

	recipe, err := h.recipeUC.GetRecipe(c.Request.Context(), recipeID, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, recipe)
}

// ListRecipes godoc
// @Summary List recipes
// @Description Retrieve user's recipes and public recipes. Can filter by recipe name using 'query' param.
// @Tags recipes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param query query string false "Search query by name"
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse{data=[]meal.Recipe}
// @Failure 401 {object} response.Response
// @Router /recipes [get]
func (h *RecipeHandler) ListRecipes(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	query := c.Query("query")

	recipes, total, err := h.recipeUC.ListRecipes(c.Request.Context(), userID, page, pageSize, query)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, recipes, page, pageSize, total)
}

// UpdateRecipe godoc
// @Summary Update a recipe
// @Description Update a recipe details or foods
// @Tags recipes
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Recipe ID"
// @Param name formData string false "Recipe Name"
// @Param description formData string false "Description"
// @Param prep_time_mins formData int false "Prep Time (mins)"
// @Param cook_time_mins formData int false "Cook Time (mins)"
// @Param servings formData int false "Servings"
// @Param instructions formData string false "Instructions"
// @Param is_public formData boolean false "Is Public"
// @Param visibility formData string false "Visibility (public, private, friends)"
// @Param foods formData string false "Foods JSON string array to replace old foods: [{\"food_id\":\"uuid\", \"quantity\":100}]"
// @Param image formData file false "Recipe Image (Max 10MB)"
// @Success 200 {object} response.Response{data=meal.Recipe}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /recipes/{id} [put]
func (h *RecipeHandler) UpdateRecipe(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	recipeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid recipe ID"))
		return
	}

	var input meal.UpdateRecipeInput
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
		input.ImageURL = &imageURL
	}

	recipe, err := h.recipeUC.UpdateRecipe(c.Request.Context(), recipeID, userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, recipe)
}

// DeleteRecipe godoc
// @Summary Delete a recipe
// @Description Delete a user's recipe
// @Tags recipes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Recipe ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /recipes/{id} [delete]
func (h *RecipeHandler) DeleteRecipe(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	recipeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid recipe ID"))
		return
	}

	err = h.recipeUC.DeleteRecipe(c.Request.Context(), recipeID, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}
