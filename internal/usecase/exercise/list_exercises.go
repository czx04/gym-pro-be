package exercise

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/workout"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"
)

type ListExercisesUseCase struct {
	exerciseRepo workout.ExerciseRepository
	validator    *validator.Validator
}

func NewListExercisesUseCase(exerciseRepo workout.ExerciseRepository, validator *validator.Validator) *ListExercisesUseCase {
	return &ListExercisesUseCase{exerciseRepo: exerciseRepo, validator: validator}
}

func (uc *ListExercisesUseCase) Excute(ctx context.Context, page, pageSize int) ([]workout.Exercise, int64, error) {
	if err := uc.validator.Validate(ctx); err != nil {
		return nil, 0, errors.Validation(err.Error())
	}

	exercises, total, err := uc.exerciseRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return exercises, total, nil
}
