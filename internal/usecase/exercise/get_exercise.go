package exercise

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/workout"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"

	"github.com/google/uuid"
)

type GetExerciseUseCase struct {
	exerciseRepo workout.ExerciseRepository
	validator    *validator.Validator
}

func NewGetExerciseUseCase(exerciseRepo workout.ExerciseRepository, validator *validator.Validator) *GetExerciseUseCase {
	return &GetExerciseUseCase{exerciseRepo: exerciseRepo, validator: validator}
}

func (uc *GetExerciseUseCase) Excute(ctx context.Context, exerciseID uuid.UUID) (*workout.Exercise, error) {
	if err := uc.validator.Validate(ctx); err != nil {
		return nil, errors.Validation(err.Error())
	}

	exercise, err := uc.exerciseRepo.GetByID(ctx, exerciseID)
	if err != nil {
		return nil, err
	}

	return exercise, nil
}
