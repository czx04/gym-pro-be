package workout

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	SessionStatusScheduled  = "scheduled"
	SessionStatusInProgress = "in_progress"
	SessionStatusCompleted  = "completed"
)

type WorkoutSession struct {
	ID                  uuid.UUID                `json:"id"`
	WorkoutScheduleID   *uuid.UUID               `json:"workout_schedule_id,omitempty"`
	UserID              uuid.UUID                `json:"user_id"`
	WorkoutPlanID       uuid.UUID                `json:"workout_plan_id"`
	ScheduledDate       *string                  `json:"scheduled_date,omitempty"` // YYYY-MM-DD
	Status              string                   `json:"status"`                   // scheduled | in_progress | completed
	StartedAt           *time.Time               `json:"started_at,omitempty"`
	CompletedAt         *time.Time               `json:"completed_at,omitempty"`
	DurationMins        *int                     `json:"duration_mins,omitempty"`
	DurationSecs        *int                     `json:"duration_secs,omitempty"`
	TotalCaloriesBurned *int                     `json:"total_calories_burned,omitempty"`
	Notes               *string                  `json:"notes,omitempty"`
	Mood                *string                  `json:"mood,omitempty"`
	DifficultyRating    *int                     `json:"difficulty_rating,omitempty"`
	CreatedAt           time.Time                `json:"created_at"`
	UpdatedAt           time.Time                `json:"updated_at"`
	Exercises           []WorkoutSessionExercise `json:"exercises,omitempty"`
	WorkoutPlan         *WorkoutPlan             `json:"workout_plan,omitempty"`
	Title               string                   `json:"title,omitempty"` // from plan for list view
}

// WorkoutSessionExercise represents exercise tracking in a session
type WorkoutSessionExercise struct {
	ID                  uuid.UUID           `json:"id"`
	WorkoutSessionID    uuid.UUID           `json:"workout_session_id"`
	ExerciseID          uuid.UUID           `json:"exercise_id"`
	Order               int                 `json:"order"`
	TargetSets          *int                `json:"target_sets,omitempty"`
	TargetReps          *int                `json:"target_reps,omitempty"`
	ActualSetsCompleted json.RawMessage     `json:"actual_sets_completed,omitempty"`
	DurationSecs        *int                `json:"duration_secs,omitempty"`
	Notes               *string             `json:"notes,omitempty"`
	Skipped             bool                `json:"skipped"`
	Exercise            *Exercise           `json:"exercise,omitempty"`
	Sets                []WorkoutSessionSet `json:"sets,omitempty"`
}

type WorkoutSessionSet struct {
	ID                       uuid.UUID  `json:"id"`
	WorkoutSessionExerciseID uuid.UUID  `json:"workout_session_exercise_id"`
	SetIndex                 int        `json:"set_index"`
	Reps                     *int       `json:"reps,omitempty"`
	WeightKg                 *float64   `json:"weight_kg,omitempty"`
	Completed                bool       `json:"completed"`
	CompletedAt              *time.Time `json:"completed_at,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

type SetData struct {
	Set      int      `json:"set"`
	Reps     *int     `json:"reps,omitempty"`
	WeightKg *float64 `json:"weight_kg,omitempty"`
	Duration *int     `json:"duration,omitempty"`
}

type StartWorkoutSessionInput struct {
	WorkoutScheduleID *uuid.UUID `json:"workout_schedule_id,omitempty"`
	WorkoutPlanID     uuid.UUID  `json:"workout_plan_id" validate:"required"`
	StartedAt         *time.Time `json:"started_at,omitempty"`
}

type LogExerciseSetInput struct {
	ExerciseID uuid.UUID `json:"exercise_id" validate:"required"`
	SetData    []SetData `json:"set_data" validate:"required,min=1"`
	Notes      *string   `json:"notes,omitempty" validate:"omitempty,max=500"`
}

type CompleteWorkoutSessionInput struct {
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	DurationSecs     *int       `json:"duration_secs,omitempty" validate:"omitempty,gte=0"`
	Notes            *string    `json:"notes,omitempty" validate:"omitempty,max=1000"`
	Mood             *string    `json:"mood,omitempty" validate:"omitempty,oneof=happy neutral tired energetic"`
	DifficultyRating *int       `json:"difficulty_rating,omitempty" validate:"omitempty,gte=1,lte=5"`
}

type CreateWorkoutSessionInput struct {
	WorkoutPlanID uuid.UUID `json:"workout_plan_id" validate:"required"`
	ScheduledDate string    `json:"scheduled_date" validate:"required"` // YYYY-MM-DD
	StartNow      bool      `json:"start_now,omitempty"`                // true = set status in_progress, started_at = now
}

type UpdateWorkoutSessionInput struct {
	Status    *string    `json:"status,omitempty" validate:"omitempty,oneof=scheduled in_progress completed"`
	StartedAt *time.Time `json:"started_at,omitempty"`
}

type UpdateSessionSetInput struct {
	Reps      *int     `json:"reps,omitempty"`
	WeightKg  *float64 `json:"weight_kg,omitempty"`
	Completed *bool    `json:"completed,omitempty"`
}

type WorkoutStats struct {
	TotalWorkouts       int     `json:"total_workouts"`
	TotalDurationMins   int     `json:"total_duration_mins"`
	TotalCaloriesBurned int     `json:"total_calories_burned"`
	AverageDuration     float64 `json:"average_duration"`
	CompletionRate      float64 `json:"completion_rate"`
}
