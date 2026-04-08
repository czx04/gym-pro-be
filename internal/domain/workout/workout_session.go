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
	RestSecs                 *int       `json:"rest_secs,omitempty"`
	Completed                bool       `json:"completed"`
	CompletedAt              *time.Time `json:"completed_at,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

type CompleteWorkoutSessionInput struct {
	StartedAt        *time.Time              `json:"started_at,omitempty"`
	CompletedAt      *time.Time              `json:"completed_at,omitempty"`
	Notes            *string                 `json:"notes,omitempty" validate:"omitempty,max=1000"`
	Mood             *string                 `json:"mood,omitempty" validate:"omitempty,oneof=happy neutral tired energetic"`
	DifficultyRating *int                    `json:"difficulty_rating,omitempty" validate:"omitempty,gte=1,lte=5"`
	Sets             []FinishSessionSetInput `json:"sets,omitempty" validate:"omitempty,dive"`
}

type FinishSessionSetInput struct {
	SetID     string   `json:"set_id" validate:"required,uuid4"`
	Reps      *int     `json:"reps,omitempty"`
	WeightKg  *float64 `json:"weight_kg,omitempty"`
	RestSecs  *int     `json:"rest_secs,omitempty" validate:"omitempty,gte=0,lte=1800"`
	Completed *bool    `json:"completed,omitempty"`
}

type CreateWorkoutSessionInput struct {
	WorkoutPlanID uuid.UUID `json:"workout_plan_id" validate:"required"`
	ScheduledDate string    `json:"scheduled_date" validate:"required"` // YYYY-MM-DD
	StartNow      bool      `json:"start_now,omitempty"`                // true = set status in_progress, started_at = now
}

type GetWeeklySummaryRequest struct {
	StartDate string `form:"start_date" validate:"required"`
	EndDate   string `form:"end_date" validate:"required"`
}

type WeeklyWorkoutMetrics struct {
	TotalWorkouts       int     `json:"total_workouts"`
	CompletedWorkouts   int     `json:"completed_workouts"`
	TotalDurationMins   int     `json:"total_duration_mins"`
	TotalCaloriesBurned int     `json:"total_calories_burned"`
	TotalSetsCompleted  int     `json:"total_sets_completed"`
	TotalRepsCompleted  int     `json:"total_reps_completed"`
	TotalVolumeKg       float64 `json:"total_volume_kg"`
	AvgWeightKg         float64 `json:"avg_weight_kg"`
	AvgRestSecs         float64 `json:"avg_rest_secs"`
	RestSamples         int     `json:"rest_samples"`
	AvgMoodScore        float64 `json:"avg_mood_score"`
	AvgDifficulty       float64 `json:"avg_difficulty"`
	CompletionRate      float64 `json:"completion_rate"`
}

type TrendDelta struct {
	Current  float64 `json:"current"`
	Previous float64 `json:"previous"`
	Delta    float64 `json:"delta"`
	Trend    string  `json:"trend"`
}

type WeeklyInsight struct {
	Code      string `json:"code"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
	Evidence  string `json:"evidence,omitempty"`
	MetricKey string `json:"metric_key,omitempty"`
}

type WeeklyWorkoutSummary struct {
	StartDate            string               `json:"start_date"`
	EndDate              string               `json:"end_date"`
	PreviousStartDate    string               `json:"previous_start_date"`
	PreviousEndDate      string               `json:"previous_end_date"`
	Current              WeeklyWorkoutMetrics `json:"current"`
	Previous             WeeklyWorkoutMetrics `json:"previous"`
	StrengthTrend        TrendDelta           `json:"strength_trend"`
	RestTrend            TrendDelta           `json:"rest_trend"`
	MoodTrend            TrendDelta           `json:"mood_trend"`
	BodyWeightTrend      TrendDelta           `json:"body_weight_trend"`
	Insights             []WeeklyInsight      `json:"insights"`
	Recommendations      []string             `json:"recommendations"`
	RecommendationSource string               `json:"recommendation_source"`
	AISummary            string               `json:"ai_summary,omitempty"`
	AIModel              string               `json:"ai_model,omitempty"`
}

type WorkoutStats struct {
	TotalWorkouts       int     `json:"total_workouts"`
	TotalDurationMins   int     `json:"total_duration_mins"`
	TotalCaloriesBurned int     `json:"total_calories_burned"`
	AverageDuration     float64 `json:"average_duration"`
	CompletionRate      float64 `json:"completion_rate"`
}
