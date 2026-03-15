package meal

import (
	"context"
	"time"

	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/pkg/cloudinary"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"

	"github.com/google/uuid"
)

type FoodUseCases struct {
	foodRepo  meal.FoodRepository
	userRepo  user.Repository
	validator *validator.Validator
}

func NewFoodUseCases(foodRepo meal.FoodRepository, userRepo user.Repository, validator *validator.Validator) *FoodUseCases {
	return &FoodUseCases{
		foodRepo:  foodRepo,
		userRepo:  userRepo,
		validator: validator,
	}
}

func (uc *FoodUseCases) CreateFood(ctx context.Context, userID uuid.UUID, input meal.CreateFoodInput) (*meal.Food, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	// Look up user to determine if they are an admin
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.InternalServer("failed to retrieve user", err)
	}

	if input.Category != nil {
		if !meal.IsValidFoodCategory(*input.Category) {
			return nil, errors.BadRequest("invalid category: must be one of protein, carb, vegetable, fruit, dairy, fat, snack, beverage, other")
		}
	}

	if input.Barcode != nil && *input.Barcode != "" {
		existing, err := uc.foodRepo.GetByBarcode(ctx, *input.Barcode)
		if err == nil && existing != nil {
			return nil, errors.BadRequest("food with this barcode already exists")
		}
	}

	isSystem := u.IsAdmin()
	now := time.Now()

	food := &meal.Food{
		ID:              uuid.New(),
		Name:            input.Name,
		Description:     input.Description,
		Brand:           input.Brand,
		ImageUrl:        input.ImageUrl,
		Barcode:         input.Barcode,
		ServingSize:     input.ServingSize,
		Unit:            input.Unit,
		Calories:        input.Calories,
		ProteinG:        input.ProteinG,
		CarbsG:          input.CarbsG,
		FatG:            input.FatG,
		FiberG:          input.FiberG,
		IsSystem:        isSystem,
		CreatedByUserID: &userID,
		Category:        input.Category,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := uc.foodRepo.Create(ctx, food); err != nil {
		return nil, errors.DatabaseError("failed to create food", err)
	}

	return food, nil
}

func (uc *FoodUseCases) GetFood(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*meal.Food, error) {
	food, err := uc.foodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseError("failed to get food", err)
	}

	// Security check: normal user can only view system foods or their own foods
	if !food.IsSystem && (food.CreatedByUserID == nil || *food.CreatedByUserID != userID) {
		return nil, errors.Forbidden("you do not have permission to view this food")
	}

	return food, nil
}

func (uc *FoodUseCases) ListFoods(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]meal.Food, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	foods, total, err := uc.foodRepo.List(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, errors.DatabaseError("failed to list foods", err)
	}

	return foods, total, nil
}

func (uc *FoodUseCases) SearchFoods(ctx context.Context, userID uuid.UUID, filter meal.SearchFoodsFilter) ([]meal.Food, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	filter.UserID = &userID

	foods, total, err := uc.foodRepo.Search(ctx, filter)
	if err != nil {
		return nil, 0, errors.DatabaseError("failed to search foods", err)
	}

	return foods, total, nil
}

func (uc *FoodUseCases) UpdateFood(ctx context.Context, id uuid.UUID, userID uuid.UUID, input meal.UpdateFoodInput) (*meal.Food, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	if input.Category != nil {
		if !meal.IsValidFoodCategory(*input.Category) {
			return nil, errors.BadRequest("invalid category: must be one of protein, carb, vegetable, fruit, dairy, fat, snack, beverage, other")
		}
	}

	if input.Barcode != nil && *input.Barcode != "" {
		existing, err := uc.foodRepo.GetByBarcode(ctx, *input.Barcode)
		if err == nil && existing != nil && existing.ID != id {
			return nil, errors.BadRequest("another food with this barcode already exists")
		}
	}

	// First check if the user has permission to update this food
	food, err := uc.foodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseError("failed to get food", err)
	}

	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.InternalServer("failed to retrieve user", err)
	}

	// Only admins can edit system foods, users can edit their own non-system foods
	if food.IsSystem && !u.IsAdmin() {
		return nil, errors.Forbidden("only admins can update system foods")
	}
	if !food.IsSystem && (food.CreatedByUserID == nil || *food.CreatedByUserID != userID) {
		return nil, errors.Forbidden("you do not have permission to update this food")
	}

	// Delete the old image from Cloudinary if it existed.
	if input.ImageUrl != nil && food.ImageUrl != nil && *input.ImageUrl != *food.ImageUrl {
		go func(oldURL string) {
			_ = cloudinary.DeleteImage(context.Background(), oldURL)
		}(*food.ImageUrl)
	}

	err = uc.foodRepo.Update(ctx, id, input)
	if err != nil {
		return nil, errors.DatabaseError("failed to update food", err)
	}

	// Return updated food
	updatedFood, err := uc.foodRepo.GetByID(ctx, id)
	return updatedFood, err
}

func (uc *FoodUseCases) DeleteFood(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {

	food, err := uc.foodRepo.GetByID(ctx, id)
	if err != nil {
		return errors.DatabaseError("failed to get food for deletion check", err)
	}

	// User is an admin, they can delete system foods? (Optional rule: admins can't delete system either, keep logic simple)
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err == nil && u.IsAdmin() && food.IsSystem {
		return errors.Forbidden("system foods cannot be deleted")
	}

	err = uc.foodRepo.Delete(ctx, id, userID)
	if err != nil {
		return errors.DatabaseError("failed to delete food", err)
	}

	// Delete associated image on Cloudinary if it exists
	if food.ImageUrl != nil && *food.ImageUrl != "" {
		go func(oldURL string) {
			_ = cloudinary.DeleteImage(context.Background(), oldURL)
		}(*food.ImageUrl)
	}

	return nil
}
