package workout

import (
	"time"

	"github.com/google/uuid"
)

// Exercise represents a pre-defined exercise in the library
type Exercise struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Category          string    `json:"category"` // cardio, strength, flexibility, stretching
	MuscleGroups      []string  `json:"muscle_groups"`
	EquipmentNeeded   []string  `json:"equipment_needed"`
	DifficultyLevel   string    `json:"difficulty_level"` // beginner, intermediate, advanced
	CaloriesPerMinute *float64  `json:"calories_per_minute,omitempty"`
	VideoURL          *string   `json:"video_url,omitempty"`
	ThumbnailURL      *string   `json:"thumbnail_url,omitempty"`
	IsActive          bool      `json:"is_active"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty"` // Admin ID
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// SearchExercisesFilter represents filters for searching exercises
type SearchExercisesFilter struct {
	Category        *string
	MuscleGroup     *string
	Equipment       *string
	DifficultyLevel *string
	Query           *string
	Page            int
	PageSize        int
}
