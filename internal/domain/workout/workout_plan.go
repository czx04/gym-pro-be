package workout

import (
	"time"

	"github.com/google/uuid"
)

// WorkoutPlan represents a user's workout plan
type WorkoutPlan struct {
	ID                    uuid.UUID `json:"id"`
	UserID                uuid.UUID `json:"user_id"`
	Title                 string    `json:"title"`
	Description           *string   `json:"description,omitempty"`
	DifficultyLevel       string    `json:"difficulty_level"`
	EstimatedDurationMins *int      `json:"estimated_duration_mins,omitempty"`
	EstimatedCalories     *int      `json:"estimated_calories,omitempty"`
	IsTemplate            bool      `json:"is_template"`
	IsPublic              bool      `json:"is_public"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	Exercises             []WorkoutPlanExercise `json:"exercises,omitempty"`
}

// WorkoutPlanExercise represents an exercise within a workout plan
type WorkoutPlanExercise struct {
	ID            uuid.UUID `json:"id"`
	WorkoutPlanID uuid.UUID `json:"workout_plan_id"`
	ExerciseID    uuid.UUID `json:"exercise_id"`
	Order         int       `json:"order"`
	Sets          *int      `json:"sets,omitempty"`
	Reps          *int      `json:"reps,omitempty"`
	DurationSecs  *int      `json:"duration_secs,omitempty"`
	RestSecs      *int      `json:"rest_secs,omitempty"`
	Notes         *string   `json:"notes,omitempty"`
	Exercise      *Exercise `json:"exercise,omitempty"`
}

// CreateWorkoutPlanInput represents input for creating a workout plan
type CreateWorkoutPlanInput struct {
	Title           string  `json:"title" validate:"required,min=3,max=200"`
	Description     *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	DifficultyLevel string  `json:"difficulty_level" validate:"required,oneof=beginner intermediate advanced"`
	IsTemplate      bool    `json:"is_template"`
	IsPublic        bool    `json:"is_public"`
}

// UpdateWorkoutPlanInput represents input for updating a workout plan
type UpdateWorkoutPlanInput struct {
	Title           *string `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Description     *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	DifficultyLevel *string `json:"difficulty_level,omitempty" validate:"omitempty,oneof=beginner intermediate advanced"`
	IsTemplate      *bool   `json:"is_template,omitempty"`
	IsPublic        *bool   `json:"is_public,omitempty"`
}

// AddExerciseToWorkoutInput represents input for adding exercise to workout
type AddExerciseToWorkoutInput struct {
	ExerciseID   uuid.UUID `json:"exercise_id" validate:"required"`
	Order        int       `json:"order" validate:"required,gte=1"`
	Sets         *int      `json:"sets,omitempty" validate:"omitempty,gte=1,lte=20"`
	Reps         *int      `json:"reps,omitempty" validate:"omitempty,gte=1,lte=100"`
	DurationSecs *int      `json:"duration_secs,omitempty" validate:"omitempty,gte=1"`
	RestSecs     *int      `json:"rest_secs,omitempty" validate:"omitempty,gte=0"`
	Notes        *string   `json:"notes,omitempty" validate:"omitempty,max=500"`
}

// UpdateExerciseInWorkoutInput represents input for updating exercise config
type UpdateExerciseInWorkoutInput struct {
	Order        *int    `json:"order,omitempty" validate:"omitempty,gte=1"`
	Sets         *int    `json:"sets,omitempty" validate:"omitempty,gte=1,lte=20"`
	Reps         *int    `json:"reps,omitempty" validate:"omitempty,gte=1,lte=100"`
	DurationSecs *int    `json:"duration_secs,omitempty" validate:"omitempty,gte=1"`
	RestSecs     *int    `json:"rest_secs,omitempty" validate:"omitempty,gte=0"`
	Notes        *string `json:"notes,omitempty" validate:"omitempty,max=500"`
}
