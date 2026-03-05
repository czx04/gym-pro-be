package exercise

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/workout"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"
)

type FilterExerciseUseCase struct {
	exerciseRepo workout.ExerciseRepository
	validator    *validator.Validator
}

func NewFilterExerciseUseCase(exerciseRepo workout.ExerciseRepository, validator *validator.Validator) *FilterExerciseUseCase {
	return &FilterExerciseUseCase{exerciseRepo: exerciseRepo, validator: validator}
}

func (uc *FilterExerciseUseCase) Excute(ctx context.Context, page, pageSize int, category, muscleGroup, equipment, difficultyLevel, query string) ([]workout.Exercise, int64, error) {
	if err := uc.validator.Validate(ctx); err != nil {
		return nil, 0, errors.Validation(err.Error())
	}

	exercises, total, err := uc.exerciseRepo.Search(ctx, workout.SearchExercisesFilter{
		Category:        &category,
		MuscleGroup:     &muscleGroup,
		Equipment:       &equipment,
		DifficultyLevel: &difficultyLevel,
		Query:           &query,
		Page:            page,
		PageSize:        pageSize,
	})
	if err != nil {
		return nil, 0, err
	}

	return exercises, total, nil
}
