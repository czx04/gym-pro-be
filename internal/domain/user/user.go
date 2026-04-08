package user

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID  `json:"id"`
	Email              string     `json:"email"`
	PasswordHash       string     `json:"-"`
	OAuthProvider      *string    `json:"oauth_provider,omitempty"`
	OAuthID            *string    `json:"oauth_id,omitempty"`
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
	PrivacySettings    *string    `json:"privacy_settings,omitempty"` // JSON string
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type CreateUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
}

type UpdateProfileInput struct {
	Name               *string    `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Bio                *string    `json:"bio,omitempty" validate:"omitempty,max=500"`
	AvatarURL          *string    `json:"avatar_url,omitempty" validate:"omitempty,url"`
	DateOfBirth        *time.Time `json:"date_of_birth,omitempty"`
	Gender             *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	HeightCm           *float64   `json:"height_cm,omitempty" validate:"omitempty,gt=0,lte=300"`
	WeightKg           *float64   `json:"weight_kg,omitempty" validate:"omitempty,gt=0,lte=500"`
	FitnessGoal        *string    `json:"fitness_goal,omitempty" validate:"omitempty,oneof=lose_weight maintain gain_muscle improve_endurance"`
	ActivityLevel      *string    `json:"activity_level,omitempty" validate:"omitempty,oneof=sedentary light moderate active very_active"`
	DailyCalorieTarget *int       `json:"daily_calorie_target,omitempty" validate:"omitempty,gte=500,lte=10000"`
	ProteinTargetG     *int       `json:"protein_target_g,omitempty" validate:"omitempty,gt=0,lte=500"`
	CarbsTargetG       *int       `json:"carbs_target_g,omitempty" validate:"omitempty,gt=0,lte=1000"`
	FatTargetG         *int       `json:"fat_target_g,omitempty" validate:"omitempty,gt=0,lte=300"`
}

type UserNutritionTarget struct {
	DailyCalorieTarget *int `json:"daily_calorie_target,omitempty"`
	ProteinTargetG     *int `json:"protein_target_g,omitempty"`
	CarbsTargetG       *int `json:"carbs_target_g,omitempty"`
	FatTargetG         *int `json:"fat_target_g,omitempty"`
}

type WeightHistory struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	WeightKg   float64   `json:"weight_kg"`
	MeasuredAt time.Time `json:"measured_at"`
	Source     string    `json:"source"`
	CreatedAt  time.Time `json:"created_at"`
}

// WeightHistoryGranularity controls how weight samples are bucketed for charts.
type WeightHistoryGranularity string

const (
	WeightHistoryGranularityDay   WeightHistoryGranularity = "day"
	WeightHistoryGranularityWeek  WeightHistoryGranularity = "week"
	WeightHistoryGranularityMonth WeightHistoryGranularity = "month"
)

// WeightHistoryPoint is one chart point: latest measurement within each calendar bucket (in timezone).
type WeightHistoryPoint struct {
	PeriodStart time.Time `json:"period_start"`
	WeightKg    float64   `json:"weight_kg"`
	MeasuredAt  time.Time `json:"measured_at"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type OAuthUserInfo struct {
	Provider  string
	ID        string
	Email     string
	Name      string
	AvatarURL *string
}

func (u *User) IsAdmin() bool {
	return strings.Contains(u.Email, "@gym-pro.com")
}
