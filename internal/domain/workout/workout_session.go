package workout

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// WorkoutSession represents an actual workout execution
type WorkoutSession struct {
	ID                 uuid.UUID  `json:"id"`
	WorkoutScheduleID  *uuid.UUID `json:"workout_schedule_id,omitempty"`
	UserID             uuid.UUID  `json:"user_id"`
	WorkoutPlanID      uuid.UUID  `json:"workout_plan_id"`
	StartedAt          time.Time  `json:"started_at"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	DurationMins       *int       `json:"duration_mins,omitempty"`
	TotalCaloriesBurned *int      `json:"total_calories_burned,omitempty"`
	Notes              *string    `json:"notes,omitempty"`
	Mood               *string    `json:"mood,omitempty"` // happy, neutral, tired, energetic
	DifficultyRating   *int       `json:"difficulty_rating,omitempty"` // 1-5
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	Exercises          []WorkoutSessionExercise `json:"exercises,omitempty"`
	WorkoutPlan        *WorkoutPlan `json:"workout_plan,omitempty"`
}

// WorkoutSessionExercise represents exercise tracking in a session
type WorkoutSessionExercise struct {
	ID                uuid.UUID       `json:"id"`
	WorkoutSessionID  uuid.UUID       `json:"workout_session_id"`
	ExerciseID        uuid.UUID       `json:"exercise_id"`
	Order             int             `json:"order"`
	TargetSets        *int            `json:"target_sets,omitempty"`
	TargetReps        *int            `json:"target_reps,omitempty"`
	ActualSetsCompleted json.RawMessage `json:"actual_sets_completed,omitempty"` // JSON array
	DurationSecs      *int            `json:"duration_secs,omitempty"`
	Notes             *string         `json:"notes,omitempty"`
	Skipped           bool            `json:"skipped"`
	Exercise          *Exercise       `json:"exercise,omitempty"`
}

// SetData represents data for a single set
type SetData struct {
	Set      int      `json:"set"`
	Reps     *int     `json:"reps,omitempty"`
	WeightKg *float64 `json:"weight_kg,omitempty"`
	Duration *int     `json:"duration,omitempty"` // seconds
}

// StartWorkoutSessionInput represents input for starting a session
type StartWorkoutSessionInput struct {
	WorkoutScheduleID *uuid.UUID `json:"workout_schedule_id,omitempty"`
	WorkoutPlanID     uuid.UUID  `json:"workout_plan_id" validate:"required"`
	StartedAt         *time.Time `json:"started_at,omitempty"`
}

// LogExerciseSetInput represents input for logging a set
type LogExerciseSetInput struct {
	ExerciseID uuid.UUID `json:"exercise_id" validate:"required"`
	SetData    []SetData `json:"set_data" validate:"required,min=1"`
	Notes      *string   `json:"notes,omitempty" validate:"omitempty,max=500"`
}

// CompleteWorkoutSessionInput represents input for completing a session
type CompleteWorkoutSessionInput struct {
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	Notes            *string    `json:"notes,omitempty" validate:"omitempty,max=1000"`
	Mood             *string    `json:"mood,omitempty" validate:"omitempty,oneof=happy neutral tired energetic"`
	DifficultyRating *int       `json:"difficulty_rating,omitempty" validate:"omitempty,gte=1,lte=5"`
}

// WorkoutStats represents workout statistics
type WorkoutStats struct {
	TotalWorkouts       int     `json:"total_workouts"`
	TotalDurationMins   int     `json:"total_duration_mins"`
	TotalCaloriesBurned int     `json:"total_calories_burned"`
	AverageDuration     float64 `json:"average_duration"`
	CompletionRate      float64 `json:"completion_rate"`
}
