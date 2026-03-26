package admin

import (
	"context"

	"gym-pro-2026-ptit/internal/domain/admin"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"

	"github.com/google/uuid"
)

type AdminUseCases struct {
	userRepo     admin.UserRepository
	exerciseRepo admin.ExerciseRepository
	foodRepo     admin.FoodRepository
	validator    *validator.Validator
}

func NewAdminUseCases(
	userRepo admin.UserRepository,
	exerciseRepo admin.ExerciseRepository,
	foodRepo admin.FoodRepository,
	v *validator.Validator,
) *AdminUseCases {
	return &AdminUseCases{
		userRepo:     userRepo,
		exerciseRepo: exerciseRepo,
		foodRepo:     foodRepo,
		validator:    v,
	}
}

func (uc *AdminUseCases) GetOverviewStats(ctx context.Context) (*admin.OverviewStats, error) {
	stats, err := uc.userRepo.GetOverviewStats(ctx)
	if err != nil {
		logger.Error("error getting overview stats", "err", err)
		return nil, err
	}
	return stats, nil
}

func (uc *AdminUseCases) ListUsers(ctx context.Context, filter admin.ListUsersFilter) ([]admin.UserSummary, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	users, total, err := uc.userRepo.ListUsers(ctx, filter)
	if err != nil {
		logger.Error("error listing users", "err", err)
		return nil, 0, err
	}
	return users, total, nil
}

func (uc *AdminUseCases) GetUser(ctx context.Context, id uuid.UUID) (*admin.UserDetail, error) {
	u, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		logger.Error("error getting user", "id", id, "err", err)
		return nil, err
	}
	return u, nil
}

func (uc *AdminUseCases) UpdateUserStatus(ctx context.Context, id uuid.UUID, input admin.UpdateUserStatusInput) error {
	if err := uc.userRepo.UpdateUserStatus(ctx, id, input.IsActive); err != nil {
		logger.Error("error updating user status", "id", id, "err", err)
		return err
	}
	return nil
}

func (uc *AdminUseCases) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := uc.userRepo.DeleteUser(ctx, id); err != nil {
		logger.Error("error deleting user", "id", id, "err", err)
		return err
	}
	return nil
}

func (uc *AdminUseCases) ListExercises(ctx context.Context, filter admin.ListExercisesFilter) ([]admin.AdminExercise, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	exercises, total, err := uc.exerciseRepo.ListExercises(ctx, filter)
	if err != nil {
		logger.Error("error listing exercises (admin)", "err", err)
		return nil, 0, err
	}
	return exercises, total, nil
}

func (uc *AdminUseCases) GetExercise(ctx context.Context, id uuid.UUID) (*admin.AdminExercise, error) {
	e, err := uc.exerciseRepo.GetExerciseByID(ctx, id)
	if err != nil {
		logger.Error("error getting exercise (admin)", "id", id, "err", err)
		return nil, err
	}
	return e, nil
}

func (uc *AdminUseCases) CreateExercise(ctx context.Context, adminID uuid.UUID, input admin.CreateExerciseInput) (*admin.AdminExercise, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.BadRequest(err.Error())
	}
	e, err := uc.exerciseRepo.CreateExercise(ctx, input, adminID)
	if err != nil {
		logger.Error("error creating exercise (admin)", "err", err)
		return nil, err
	}
	return e, nil
}

func (uc *AdminUseCases) UpdateExercise(ctx context.Context, id uuid.UUID, input admin.UpdateExerciseInput) (*admin.AdminExercise, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.BadRequest(err.Error())
	}
	e, err := uc.exerciseRepo.UpdateExercise(ctx, id, input)
	if err != nil {
		logger.Error("error updating exercise (admin)", "id", id, "err", err)
		return nil, err
	}
	return e, nil
}

func (uc *AdminUseCases) DeleteExercise(ctx context.Context, id uuid.UUID) error {
	if err := uc.exerciseRepo.DeleteExercise(ctx, id); err != nil {
		logger.Error("error deleting exercise (admin)", "id", id, "err", err)
		return err
	}
	return nil
}

func (uc *AdminUseCases) ListFoods(ctx context.Context, filter admin.ListFoodsFilter) ([]admin.AdminFood, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	foods, total, err := uc.foodRepo.ListFoods(ctx, filter)
	if err != nil {
		logger.Error("error listing foods (admin)", "err", err)
		return nil, 0, err
	}
	return foods, total, nil
}

func (uc *AdminUseCases) GetFood(ctx context.Context, id uuid.UUID) (*admin.AdminFood, error) {
	f, err := uc.foodRepo.GetFoodByID(ctx, id)
	if err != nil {
		logger.Error("error getting food (admin)", "id", id, "err", err)
		return nil, err
	}
	return f, nil
}

func (uc *AdminUseCases) CreateSystemFood(ctx context.Context, adminID uuid.UUID, input admin.CreateSystemFoodInput) (*admin.AdminFood, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.BadRequest(err.Error())
	}
	f, err := uc.foodRepo.CreateSystemFood(ctx, input, adminID)
	if err != nil {
		logger.Error("error creating system food", "err", err)
		return nil, err
	}
	return f, nil
}

func (uc *AdminUseCases) UpdateFood(ctx context.Context, id uuid.UUID, input admin.AdminUpdateFoodInput) (*admin.AdminFood, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.BadRequest(err.Error())
	}
	f, err := uc.foodRepo.UpdateFood(ctx, id, input)
	if err != nil {
		logger.Error("error updating food (admin)", "id", id, "err", err)
		return nil, err
	}
	return f, nil
}

func (uc *AdminUseCases) DeleteFood(ctx context.Context, id uuid.UUID) error {
	if err := uc.foodRepo.DeleteFood(ctx, id); err != nil {
		logger.Error("error deleting food (admin)", "id", id, "err", err)
		return err
	}
	return nil
}
