package meal

import (
	"time"

	"github.com/google/uuid"
)

// MealDaily represents daily nutrition targets log for a user.
type MealDaily struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Date           time.Time `json:"date"`
	TargetCalories float64   `json:"target_calories"`
	TargetProteinG float64   `json:"target_protein_g"`
	TargetCarbsG   float64   `json:"target_carbs_g"`
	TargetFatG     float64   `json:"target_fat_g"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// DailyNutritionTargetResponse represents the float-based target response
type DailyNutritionTargetResponse struct {
	DailyCalorieTarget *float64 `json:"daily_calorie_target,omitempty"`
	ProteinTargetG     *float64 `json:"protein_target_g,omitempty"`
	CarbsTargetG       *float64 `json:"carbs_target_g,omitempty"`
	FatTargetG         *float64 `json:"fat_target_g,omitempty"`
}
