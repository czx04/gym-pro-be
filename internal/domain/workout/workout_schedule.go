package workout

import (
	"time"

	"github.com/google/uuid"
)

// WorkoutSchedule represents a scheduled workout
type WorkoutSchedule struct {
	ID             uuid.UUID  `json:"id"`
	WorkoutPlanID  uuid.UUID  `json:"workout_plan_id"`
	UserID         uuid.UUID  `json:"user_id"`
	ScheduledDate  time.Time  `json:"scheduled_date"`
	ScheduledTime  *string    `json:"scheduled_time,omitempty"` // HH:MM format
	RecurrenceRule *string    `json:"recurrence_rule,omitempty"` // e.g., "WEEKLY:MON,WED,FRI"
	IsCompleted    bool       `json:"is_completed"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	WorkoutPlan    *WorkoutPlan `json:"workout_plan,omitempty"`
}

// CreateWorkoutScheduleInput represents input for scheduling a workout
type CreateWorkoutScheduleInput struct {
	WorkoutPlanID  uuid.UUID `json:"workout_plan_id" validate:"required"`
	ScheduledDate  time.Time `json:"scheduled_date" validate:"required"`
	ScheduledTime  *string   `json:"scheduled_time,omitempty" validate:"omitempty,len=5"` // HH:MM
	RecurrenceRule *string   `json:"recurrence_rule,omitempty" validate:"omitempty,max=100"`
}

// BulkScheduleInput represents input for bulk scheduling
type BulkScheduleInput struct {
	WorkoutPlanID uuid.UUID   `json:"workout_plan_id" validate:"required"`
	Dates         []time.Time `json:"dates" validate:"required,min=1"`
	ScheduledTime *string     `json:"scheduled_time,omitempty" validate:"omitempty,len=5"`
}

// UpdateScheduleInput represents input for updating a schedule
type UpdateScheduleInput struct {
	ScheduledDate *time.Time `json:"scheduled_date,omitempty"`
	ScheduledTime *string    `json:"scheduled_time,omitempty" validate:"omitempty,len=5"`
}

// GetScheduleFilter represents filters for getting schedules
type GetScheduleFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Completed *bool
}
