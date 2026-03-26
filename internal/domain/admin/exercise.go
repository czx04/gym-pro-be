package admin

import (
	"time"

	"github.com/google/uuid"
)

type AdminExercise struct {
	ID                uuid.UUID  `json:"id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	Category          string     `json:"category"`
	MuscleGroups      []string   `json:"muscle_groups"`
	EquipmentNeeded   []string   `json:"equipment_needed"`
	DifficultyLevel   string     `json:"difficulty_level"`
	CaloriesPerMinute *float64   `json:"calories_per_minute,omitempty"`
	VideoURL          *string    `json:"video_url,omitempty"`
	ThumbnailURL      *string    `json:"thumbnail_url,omitempty"`
	IsActive          bool       `json:"is_active"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type CreateExerciseInput struct {
	Name              string   `json:"name" validate:"required,min=2,max=200"`
	Description       string   `json:"description" validate:"required,min=5,max=1000"`
	Category          string   `json:"category" validate:"required,oneof=cardio strength flexibility stretching"`
	MuscleGroups      []string `json:"muscle_groups" validate:"required,min=1"`
	EquipmentNeeded   []string `json:"equipment_needed" validate:"required"`
	DifficultyLevel   string   `json:"difficulty_level" validate:"required,oneof=beginner intermediate advanced"`
	CaloriesPerMinute *float64 `json:"calories_per_minute,omitempty" validate:"omitempty,gt=0"`
	VideoURL          *string  `json:"video_url,omitempty" validate:"omitempty,url"`
	ThumbnailURL      *string  `json:"thumbnail_url,omitempty" validate:"omitempty,url"`
	IsActive          bool     `json:"is_active"`
}

type UpdateExerciseInput struct {
	Name              *string  `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Description       *string  `json:"description,omitempty" validate:"omitempty,min=5,max=1000"`
	Category          *string  `json:"category,omitempty" validate:"omitempty,oneof=cardio strength flexibility stretching"`
	MuscleGroups      []string `json:"muscle_groups,omitempty"`
	EquipmentNeeded   []string `json:"equipment_needed,omitempty"`
	DifficultyLevel   *string  `json:"difficulty_level,omitempty" validate:"omitempty,oneof=beginner intermediate advanced"`
	CaloriesPerMinute *float64 `json:"calories_per_minute,omitempty" validate:"omitempty,gt=0"`
	VideoURL          *string  `json:"video_url,omitempty" validate:"omitempty,url"`
	ThumbnailURL      *string  `json:"thumbnail_url,omitempty" validate:"omitempty,url"`
	IsActive          *bool    `json:"is_active,omitempty"`
}

type ListExercisesFilter struct {
	Query           *string
	Category        *string
	MuscleGroup     *string
	DifficultyLevel *string
	IsActive        *bool
	Page            int
	PageSize        int
}
