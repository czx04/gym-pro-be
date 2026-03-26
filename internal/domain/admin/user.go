package admin

import (
	"time"

	"github.com/google/uuid"
)

type UserSummary struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	Name          string     `json:"name"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	Gender        *string    `json:"gender,omitempty"`
	FitnessGoal   *string    `json:"fitness_goal,omitempty"`
	ActivityLevel *string    `json:"activity_level,omitempty"`
	OAuthProvider *string    `json:"oauth_provider,omitempty"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
}

type UserDetail struct {
	ID                 uuid.UUID  `json:"id"`
	Email              string     `json:"email"`
	Name               string     `json:"name"`
	Bio                *string    `json:"bio,omitempty"`
	AvatarURL          *string    `json:"avatar_url,omitempty"`
	DateOfBirth        *time.Time `json:"date_of_birth,omitempty"`
	Gender             *string    `json:"gender,omitempty"`
	HeightCm           *float64   `json:"height_cm,omitempty"`
	WeightKg           *float64   `json:"weight_kg,omitempty"`
	FitnessGoal        *string    `json:"fitness_goal,omitempty"`
	ActivityLevel      *string    `json:"activity_level,omitempty"`
	DailyCalorieTarget *int       `json:"daily_calorie_target,omitempty"`
	ProteinTargetG     *int       `json:"protein_target_g,omitempty"`
	CarbsTargetG       *int       `json:"carbs_target_g,omitempty"`
	FatTargetG         *int       `json:"fat_target_g,omitempty"`
	OAuthProvider      *string    `json:"oauth_provider,omitempty"`
	IsActive           bool       `json:"is_active"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type UpdateUserStatusInput struct {
	IsActive bool `json:"is_active"`
}

type ListUsersFilter struct {
	Query    *string
	Gender   *string
	IsActive *bool
	Page     int
	PageSize int
}
