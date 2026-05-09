package meal

import (
	"context"
	"encoding/json"
	"time"

	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/pkg/cloudinary"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/utils"
	"gym-pro-2026-ptit/pkg/validator"

	"github.com/google/uuid"
)

type RecipeUseCases struct {
	recipeRepo meal.RecipeRepository
	foodRepo   meal.FoodRepository
	userRepo   user.Repository
	validator  *validator.Validator
}

func NewRecipeUseCases(
	recipeRepo meal.RecipeRepository,
	foodRepo meal.FoodRepository,
	userRepo user.Repository,
	validator *validator.Validator,
) *RecipeUseCases {
	return &RecipeUseCases{
		recipeRepo: recipeRepo,
		foodRepo:   foodRepo,
		userRepo:   userRepo,
		validator:  validator,
	}
}

func (uc *RecipeUseCases) CreateRecipe(ctx context.Context, userID uuid.UUID, input meal.CreateRecipeInput) (*meal.Recipe, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	var foodInputs []meal.AddFoodToRecipeInput
	if input.Foods != "" {
		if err := json.Unmarshal([]byte(input.Foods), &foodInputs); err != nil {
			return nil, errors.BadRequest("invalid foods JSON format")
		}
	}

	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.InternalServer("failed to retrieve user", err)
	}

	if input.Visibility == "" {
		input.Visibility = "public"
	}

	now := time.Now()
	recipe := &meal.Recipe{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         input.Name,
		Description:  input.Description,
		PrepTimeMins: input.PrepTimeMins,
		CookTimeMins: input.CookTimeMins,
		Servings:     input.Servings,
		Instructions: input.Instructions,
		ImageURL:     input.ImageURL,
		IsPublic:     input.IsPublic,
		Visibility:   input.Visibility,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.recipeRepo.Create(ctx, recipe); err != nil {
		return nil, errors.DatabaseError("failed to create recipe", err)
	}

	// Add foods and calculate macros
	if err := uc.processRecipeFoods(ctx, recipe.ID, foodInputs); err != nil {
		return nil, err
	}

	// Fetch fresh to return with computed fields
	createdRecipe, err := uc.recipeRepo.GetByID(ctx, recipe.ID)
	if err != nil {
		return nil, err
	}

	uc.roundRecipe(createdRecipe)
	return createdRecipe, nil
}

func (uc *RecipeUseCases) GetRecipe(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*meal.Recipe, error) {
	recipe, err := uc.recipeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseError("failed to get recipe", err)
	}
	if recipe == nil {
		return nil, errors.NotFound("recipe not found")
	}

	// Basic visibility check
	if recipe.UserID != userID && !recipe.IsPublic {
		return nil, errors.Forbidden("you do not have permission to view this recipe")
	}

	uc.roundRecipe(recipe)
	return recipe, nil
}

func (uc *RecipeUseCases) ListRecipes(ctx context.Context, userID uuid.UUID, page, pageSize int, query string) ([]meal.Recipe, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	recipes, total, err := uc.recipeRepo.GetByUserID(ctx, userID, page, pageSize, query)
	if err != nil {
		return nil, 0, errors.DatabaseError("failed to list recipes", err)
	}

	for i := range recipes {
		uc.roundRecipe(&recipes[i])
	}

	return recipes, total, nil
}

func (uc *RecipeUseCases) UpdateRecipe(ctx context.Context, id uuid.UUID, userID uuid.UUID, input meal.UpdateRecipeInput) (*meal.Recipe, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	// Check permission
	recipe, err := uc.recipeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseError("failed to get recipe", err)
	}
	if recipe == nil {
		return nil, errors.NotFound("recipe not found")
	}
	if recipe.UserID != userID {
		return nil, errors.Forbidden("you do not have permission to edit this recipe")
	}

	// Delete old image if new one is provided
	if input.ImageURL != nil && recipe.ImageURL != nil && *input.ImageURL != *recipe.ImageURL {
		go func(oldURL string) {
			_ = cloudinary.DeleteImage(context.Background(), oldURL)
		}(*recipe.ImageURL)
	}

	err = uc.recipeRepo.Update(ctx, id, input)
	if err != nil {
		return nil, errors.DatabaseError("failed to update recipe", err)
	}

	if input.Foods != nil {
		var foodInputs []meal.AddFoodToRecipeInput
		if *input.Foods != "" {
			if err := json.Unmarshal([]byte(*input.Foods), &foodInputs); err != nil {
				return nil, errors.BadRequest("invalid foods JSON format")
			}
		}

		if err := uc.recipeRepo.ClearFoods(ctx, id); err != nil {
			return nil, errors.DatabaseError("failed to clear old foods", err)
		}

		if err := uc.processRecipeFoods(ctx, id, foodInputs); err != nil {
			return nil, err
		}
	}

	recipe, err = uc.recipeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	uc.roundRecipe(recipe)
	return recipe, nil
}

func (uc *RecipeUseCases) DeleteRecipe(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	recipe, err := uc.recipeRepo.GetByID(ctx, id)
	if err != nil {
		return errors.DatabaseError("failed to fetch recipe for deletion check", err)
	}
	if recipe == nil {
		return errors.NotFound("recipe not found")
	}

	if recipe.UserID != userID {
		return errors.Forbidden("you can only delete your own recipes")
	}

	err = uc.recipeRepo.Delete(ctx, id)
	if err != nil {
		return errors.DatabaseError("failed to delete recipe", err)
	}

	if recipe.ImageURL != nil && *recipe.ImageURL != "" {
		go func(oldURL string) {
			_ = cloudinary.DeleteImage(context.Background(), oldURL)
		}(*recipe.ImageURL)
	}

	return nil
}

func (uc *RecipeUseCases) processRecipeFoods(ctx context.Context, recipeID uuid.UUID, foodInputs []meal.AddFoodToRecipeInput) error {
	for _, fi := range foodInputs {
		foodData, err := uc.foodRepo.GetByID(ctx, fi.FoodID)
		if err != nil || foodData == nil {
			return errors.BadRequest("one or more foods do not exist")
		}

		ratio := fi.Quantity / foodData.ServingSize

		rf := &meal.RecipeFood{
			ID:       uuid.New(),
			RecipeID: recipeID,
			FoodID:   fi.FoodID,
			Quantity: fi.Quantity,
			Unit:     foodData.Unit, // Inherit unit from base food
			Calories: foodData.Calories * ratio,
			ProteinG: foodData.ProteinG * ratio,
			CarbsG:   foodData.CarbsG * ratio,
			FatG:     foodData.FatG * ratio,
		}

		if err := uc.recipeRepo.AddFood(ctx, recipeID, rf); err != nil {
			return errors.DatabaseError("failed to append food to recipe", err)
		}
	}
	return nil
}

func (uc *RecipeUseCases) roundRecipe(recipe *meal.Recipe) {
	if recipe == nil {
		return
	}
	recipe.TotalCalories = utils.RoundToTwo(recipe.TotalCalories)
	recipe.TotalProteinG = utils.RoundToTwo(recipe.TotalProteinG)
	recipe.TotalCarbsG = utils.RoundToTwo(recipe.TotalCarbsG)
	recipe.TotalFatG = utils.RoundToTwo(recipe.TotalFatG)
	recipe.PerServingCalories = utils.RoundToTwo(recipe.PerServingCalories)
	recipe.PerServingProteinG = utils.RoundToTwo(recipe.PerServingProteinG)
	recipe.PerServingCarbsG = utils.RoundToTwo(recipe.PerServingCarbsG)
	recipe.PerServingFatG = utils.RoundToTwo(recipe.PerServingFatG)

	for i := range recipe.Foods {
		rf := &recipe.Foods[i]
		rf.Calories = utils.RoundToTwo(rf.Calories)
		rf.ProteinG = utils.RoundToTwo(rf.ProteinG)
		rf.CarbsG = utils.RoundToTwo(rf.CarbsG)
		rf.FatG = utils.RoundToTwo(rf.FatG)

		if rf.Food != nil {
			rf.Food.Calories = utils.RoundToTwo(rf.Food.Calories)
			rf.Food.ProteinG = utils.RoundToTwo(rf.Food.ProteinG)
			rf.Food.CarbsG = utils.RoundToTwo(rf.Food.CarbsG)
			rf.Food.FatG = utils.RoundToTwo(rf.Food.FatG)
		}
	}
}
