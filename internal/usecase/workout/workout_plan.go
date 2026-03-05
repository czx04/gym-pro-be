package workout

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/workout"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"
	"time"

	"github.com/google/uuid"
)

// CreateWorkoutPlanUseCase handles creating workout plans
type CreateWorkoutPlanUseCase struct {
	workoutPlanRepo workout.WorkoutPlanRepository
	validator       *validator.Validator
}

// NewCreateWorkoutPlanUseCase creates a new use case
func NewCreateWorkoutPlanUseCase(
	workoutPlanRepo workout.WorkoutPlanRepository,
	validator *validator.Validator,
) *CreateWorkoutPlanUseCase {
	return &CreateWorkoutPlanUseCase{
		workoutPlanRepo: workoutPlanRepo,
		validator:       validator,
	}
}

// Execute creates a new workout plan
func (uc *CreateWorkoutPlanUseCase) Execute(ctx context.Context, userID uuid.UUID, input workout.CreateWorkoutPlanInput) (*workout.WorkoutPlan, error) {
	// TODO: 1. Validate input
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	// TODO: 2. Create workout plan entity
	plan := &workout.WorkoutPlan{
		ID:              uuid.New(),
		UserID:          userID,
		Title:           input.Title,
		Description:     input.Description,
		DifficultyLevel: input.DifficultyLevel,
		IsTemplate:      input.IsTemplate,
		IsPublic:        input.IsPublic,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// TODO: 3. Save to database
	if err := uc.workoutPlanRepo.Create(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// AddExerciseToWorkoutUseCase handles adding exercises to workout plans
type AddExerciseToWorkoutUseCase struct {
	workoutPlanRepo workout.WorkoutPlanRepository
	exerciseRepo    workout.ExerciseRepository
	validator       *validator.Validator
}

// NewAddExerciseToWorkoutUseCase creates a new use case
func NewAddExerciseToWorkoutUseCase(
	workoutPlanRepo workout.WorkoutPlanRepository,
	exerciseRepo workout.ExerciseRepository,
	validator *validator.Validator,
) *AddExerciseToWorkoutUseCase {
	return &AddExerciseToWorkoutUseCase{
		workoutPlanRepo: workoutPlanRepo,
		exerciseRepo:    exerciseRepo,
		validator:       validator,
	}
}

// Execute adds an exercise to a workout plan
func (uc *AddExerciseToWorkoutUseCase) Execute(ctx context.Context, planID uuid.UUID, input workout.AddExerciseToWorkoutInput) error {
	if err := uc.validator.Validate(input); err != nil {
		return errors.Validation(err.Error())
	}

	_, err := uc.workoutPlanRepo.GetByID(ctx, planID)
	if err != nil {
		return err
	}

	_, err = uc.exerciseRepo.GetByID(ctx, input.ExerciseID)
	if err != nil {
		return errors.NotFound("exercise")
	}
	planExercise := &workout.WorkoutPlanExercise{
		ID:            uuid.New(),
		WorkoutPlanID: planID,
		ExerciseID:    input.ExerciseID,
		Order:         input.Order,
		Sets:          input.Sets,
		Reps:          input.Reps,
		DurationSecs:  input.DurationSecs,
		RestSecs:      input.RestSecs,
		Notes:         input.Notes,
	}

	if err := uc.workoutPlanRepo.AddExercise(ctx, planID, planExercise); err != nil {
		return err
	}

	return nil
}

// TODO: Implement more use cases:
// - GetWorkoutPlanUseCase
// - UpdateWorkoutPlanUseCase
// - DeleteWorkoutPlanUseCase
// - RemoveExerciseFromWorkoutUseCase
// - UpdateExerciseInWorkoutUseCase
