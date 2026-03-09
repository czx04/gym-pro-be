package workout

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/domain/workout"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"
	"time"

	"github.com/google/uuid"
)

// WorkoutUseCases groups all workout use cases with a single dependency set.
type WorkoutUseCases struct {
	db              *database.DB
	workoutPlanRepo workout.WorkoutPlanRepository
	exerciseRepo    workout.ExerciseRepository
	validator       *validator.Validator
}

// NewWorkoutUseCases creates the workout use cases container.
func NewWorkoutUseCases(
	db *database.DB,
	workoutPlanRepo workout.WorkoutPlanRepository,
	exerciseRepo workout.ExerciseRepository,
	validator *validator.Validator,
) *WorkoutUseCases {
	return &WorkoutUseCases{
		db:              db,
		workoutPlanRepo: workoutPlanRepo,
		exerciseRepo:    exerciseRepo,
		validator:       validator,
	}
}

func (uc *WorkoutUseCases) CreateWorkoutPlan(ctx context.Context, u *user.User, input workout.CreateWorkoutPlanInput) (*workout.WorkoutPlan, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}
	isTemplate := u.IsAdmin()

	tx, err := uc.db.Begin(ctx)
	if err != nil {
		return nil, errors.DatabaseError("begin transaction", err)
	}
	defer tx.Rollback(ctx)

	planRepo := uc.workoutPlanRepo.WithTx(tx)

	plan := &workout.WorkoutPlan{
		ID:              uuid.New(),
		UserID:          u.ID,
		Title:           input.Title,
		Description:     input.Description,
		DifficultyLevel: input.DifficultyLevel,
		IsTemplate:      isTemplate,
		IsPublic:        input.IsPublic,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := planRepo.Create(ctx, plan); err != nil {
		return nil, err
	}

	exercises := make([]*workout.WorkoutPlanExercise, len(input.Exercises))
	for i, exercise := range input.Exercises {
		exercises[i] = &workout.WorkoutPlanExercise{
			WorkoutPlanID: plan.ID,
			ExerciseID:    exercise.ExerciseID,
			Order:         exercise.Order,
			Sets:          exercise.Sets,
			Reps:          exercise.Reps,
			DurationSecs:  exercise.DurationSecs,
			RestSecs:      exercise.RestSecs,
			Notes:         exercise.Notes,
		}
	}

	if len(exercises) > 0 {
		if err := planRepo.AddExercise(ctx, plan.ID, exercises); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.DatabaseError("commit transaction", err)
	}
	return plan, nil
}

func (uc *WorkoutUseCases) ListWorkoutPlans(ctx context.Context, user user.User, page, pageSize int) ([]workout.WorkoutPlan, int64, error) {
	return uc.workoutPlanRepo.GetByUserID(ctx, user.ID, page, pageSize)
}

func (uc *WorkoutUseCases) GetWorkoutPlan(ctx context.Context, userID uuid.UUID, planID string) (*workout.WorkoutPlan, error) {
	uuidPlanID, err := uuid.Parse(planID)
	if err != nil {
		logger.Error("error parsing plan ID", "err", err, "planID", planID)
		return nil, errors.BadRequest("invalid plan ID")
	}
	plan, err := uc.workoutPlanRepo.GetByID(ctx, uuidPlanID)
	if err != nil {
		logger.Error("error getting workout plan by ID", "err", err, "planID", planID)
		return nil, errors.DatabaseError("get workout plan by ID", err)
	}
	if plan.UserID != userID {
		logger.Error("user is not allowed to get this workout plan", "userID", userID, "planID", planID)
		return nil, errors.Forbidden("you are not allowed to get this workout plan")
	}

	exercises, err := uc.workoutPlanRepo.GetExercises(ctx, plan.ID)
	if err != nil {
		logger.Error("error getting exercises by plan id", "err", err, "planID", planID)
		return nil, errors.DatabaseError("get exercises by plan id", err)
	}
	plan.Exercises = exercises
	return plan, nil
}
func (uc *WorkoutUseCases) DeleteWorkoutPlan(ctx context.Context, userID uuid.UUID, planID string) error {
	uuidPlanID, err := uuid.Parse(planID)
	if err != nil {
		logger.Error("error parsing plan ID", "err", err, "planID", planID)
		return errors.BadRequest("invalid plan ID")
	}
	plan, err := uc.workoutPlanRepo.GetByID(ctx, uuidPlanID)
	if err != nil {
		logger.Error("error getting workout plan by ID", "err", err, "planID", planID)
		return errors.DatabaseError("get workout plan by ID", err)
	}
	if plan.UserID != userID {
		logger.Error("user is not allowed to delete this workout plan", "userID", userID, "planID", planID)
		return errors.Forbidden("you are not allowed to delete this workout plan")
	}
	return uc.workoutPlanRepo.Delete(ctx, plan.ID)
}

func (uc *WorkoutUseCases) UpdateWorkoutPlan(ctx context.Context, userID uuid.UUID, input workout.UpdateWorkoutPlanInput) (*workout.WorkoutPlan, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	plan, err := uc.workoutPlanRepo.GetByID(ctx, input.ID)
	if err != nil {
		logger.Error("error getting workout plan by ID", "err", err, "planID", input.ID)
		return nil, errors.DatabaseError("get workout plan by ID", err)
	}
	if plan.UserID != userID {
		logger.Error("user is not allowed to update this workout plan", "userID", userID, "planID", input.ID)
		return nil, errors.Forbidden("you are not allowed to update this workout plan")
	}

	uc.buildWorkoutPlanFromUpdateInput(plan, input)

	db, err := uc.db.Begin(ctx)
	if err != nil {
		return nil, errors.DatabaseError("begin transaction", err)
	}
	defer db.Rollback(ctx)

	planRepo := uc.workoutPlanRepo.WithTx(db)

	if err := planRepo.Update(ctx, plan); err != nil {
		return nil, errors.DatabaseError("update workout plan", err)
	}

	if input.IsUpdateExercises && len(input.Exercises) > 0 {
		if err := planRepo.RemoveExercise(ctx, plan.ID); err != nil {
			return nil, errors.DatabaseError("remove exercise from workout plan", err)
		}
		exercises := make([]*workout.WorkoutPlanExercise, len(input.Exercises))
		for i, exercise := range input.Exercises {
			exercises[i] = &workout.WorkoutPlanExercise{
				WorkoutPlanID: plan.ID,
				ExerciseID:    exercise.ID,
				Order:         exercise.Order,
				Sets:          exercise.Sets,
				Reps:          exercise.Reps,
				DurationSecs:  exercise.DurationSecs,
				RestSecs:      exercise.RestSecs,
				Notes:         exercise.Notes,
			}
		}
		if err := planRepo.AddExercise(ctx, plan.ID, exercises); err != nil {
			return nil, errors.DatabaseError("update exercise in workout plan", err)
		}
	}

	if err := db.Commit(ctx); err != nil {
		return nil, errors.DatabaseError("commit transaction", err)
	}
	return plan, nil
}

func (uc *WorkoutUseCases) buildWorkoutPlanFromUpdateInput(currentPlan *workout.WorkoutPlan, input workout.UpdateWorkoutPlanInput) {
	if input.Title != nil {
		currentPlan.Title = *input.Title
	}
	if input.Description != nil {
		currentPlan.Description = input.Description
	}
	if input.DifficultyLevel != nil {
		currentPlan.DifficultyLevel = *input.DifficultyLevel
	}
	if input.IsTemplate != nil {
		currentPlan.IsTemplate = *input.IsTemplate
	}
	if input.IsPublic != nil {
		currentPlan.IsPublic = *input.IsPublic
	}
	currentPlan.UpdatedAt = time.Now()
}
