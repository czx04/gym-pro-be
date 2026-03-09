package exercise

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/workout"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/pkg/validator"

	"github.com/google/uuid"
)

// ExerciseUseCases groups all exercise use cases with a single dependency set.
type ExerciseUseCases struct {
	exerciseRepo workout.ExerciseRepository
	validator    *validator.Validator
}

// NewExerciseUseCases creates the exercise use cases container.
func NewExerciseUseCases(
	exerciseRepo workout.ExerciseRepository,
	validator *validator.Validator,
) *ExerciseUseCases {
	return &ExerciseUseCases{
		exerciseRepo: exerciseRepo,
		validator:    validator,
	}
}

func (uc *ExerciseUseCases) ListExercises(ctx context.Context, page, pageSize int) ([]workout.Exercise, int64, error) {
	exercises, total, err := uc.exerciseRepo.List(ctx, page, pageSize)
	if err != nil {
		logger.Error("error listing exercises", "err", err)
		return nil, 0, err
	}
	return exercises, total, nil
}

func (uc *ExerciseUseCases) GetExercise(ctx context.Context, exerciseID uuid.UUID) (*workout.Exercise, error) {
	exercise, err := uc.exerciseRepo.GetByID(ctx, exerciseID)
	if err != nil {
		logger.Error("error getting exercise", "err", err)
		return nil, err
	}
	return exercise, nil
}

func (uc *ExerciseUseCases) FilterExercises(ctx context.Context, page, pageSize int, category, muscleGroup, equipment, difficultyLevel, query string) ([]workout.Exercise, int64, error) {
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
		logger.Error("error filtering exercises", "err", err)
		return nil, 0, err
	}
	return exercises, total, nil
}
