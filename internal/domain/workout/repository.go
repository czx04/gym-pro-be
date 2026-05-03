package workout

import (
	"context"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"time"

	"github.com/google/uuid"
)

// ExerciseRepository defines the interface for exercise data access
type ExerciseRepository interface {
	WithTx(tx *database.DB) ExerciseRepository
	// Create creates a new exercise (admin only)
	Create(ctx context.Context, exercise *Exercise) error

	// GetByID retrieves an exercise by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Exercise, error)

	// List retrieves exercises with pagination
	List(ctx context.Context, page, pageSize int) ([]Exercise, int64, error)

	// Search searches exercises with filters
	Search(ctx context.Context, filter SearchExercisesFilter) ([]Exercise, int64, error)

	// Update updates an exercise
	Update(ctx context.Context, exercise *Exercise) error

	// Delete soft deletes an exercise
	Delete(ctx context.Context, id uuid.UUID) error
}

// WorkoutPlanRepository defines the interface for workout plan data access
type WorkoutPlanRepository interface {
	WithTx(tx *database.DB) WorkoutPlanRepository
	// Create creates a new workout plan
	Create(ctx context.Context, plan *WorkoutPlan) error

	// GetByID retrieves a workout plan by ID with exercises
	GetByID(ctx context.Context, id uuid.UUID) (*WorkoutPlan, error)

	// GetByUserID retrieves workout plans for a user
	GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]WorkoutPlan, int64, error)

	// Update updates a workout plan
	Update(ctx context.Context, plan *WorkoutPlan) error

	// Delete deletes a workout plan
	Delete(ctx context.Context, id uuid.UUID) error

	// AddExercise adds an exercise to a workout plan
	AddExercise(ctx context.Context, planID uuid.UUID, exercises []*WorkoutPlanExercise) error

	// RemoveExercise removes an exercise from a workout plan
	RemoveExercise(ctx context.Context, planID uuid.UUID) error

	// GetExercises retrieves exercises in a workout plan
	GetExercises(ctx context.Context, planID uuid.UUID) ([]WorkoutPlanExercise, error)
}

// WorkoutSessionRepository defines the interface for workout session data access
type WorkoutSessionRepository interface {
	WithTx(tx *database.DB) WorkoutSessionRepository
	Create(ctx context.Context, session *WorkoutSession) error
	GetByID(ctx context.Context, id uuid.UUID) (*WorkoutSession, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]WorkoutSession, int64, error)
	Complete(ctx context.Context, id uuid.UUID, input CompleteWorkoutSessionInput) error
	AddExerciseLog(ctx context.Context, sessionID uuid.UUID, exercise *WorkoutSessionExercise) error
	GetExercises(ctx context.Context, sessionID uuid.UUID) ([]WorkoutSessionExercise, error)
	GetStats(ctx context.Context, userID uuid.UUID) (*WorkoutStats, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetScheduledDates(ctx context.Context, userID uuid.UUID, month, year int) ([]string, error)
	GetByDate(ctx context.Context, userID uuid.UUID, date string) ([]WorkoutSession, error)
	UpdateSetsBulk(ctx context.Context, sessionID uuid.UUID, sets []FinishSessionSetInput) error
	GetWeeklyAggregate(ctx context.Context, userID uuid.UUID, start, end time.Time) (*WeeklyWorkoutMetrics, error)

	GetExerciseStats(ctx context.Context, userID, exerciseID uuid.UUID) (*ExerciseStats, error)

	// GetProfileWorkoutStats returns lightweight stats for the profile screen.
	GetProfileWorkoutStats(ctx context.Context, userID uuid.UUID) (totalWorkouts int64, totalWorkoutDays int64, err error)
}
