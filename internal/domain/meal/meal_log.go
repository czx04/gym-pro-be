package meal

import (
	"time"

	"github.com/google/uuid"
)

// MealLog represents a meal consumption log
type MealLog struct {
	ID            uuid.UUID      `json:"id"`
	UserID        uuid.UUID      `json:"user_id"`
	LogDate       time.Time      `json:"log_date"` // Date only
	MealTime      string         `json:"meal_time"` // breakfast, lunch, dinner, snack, other
	TotalCalories float64        `json:"total_calories"`
	TotalProteinG float64        `json:"total_protein_g"`
	TotalCarbsG   float64        `json:"total_carbs_g"`
	TotalFatG     float64        `json:"total_fat_g"`
	Notes         *string        `json:"notes,omitempty"`
	Mood          *string        `json:"mood,omitempty"` // great, good, okay, tired
	EnergyLevel   *int           `json:"energy_level,omitempty"` // 1-5
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Items         []MealLogItem  `json:"items,omitempty"`
}

// MealLogItem represents a food or recipe in a meal log
type MealLogItem struct {
	ID          uuid.UUID  `json:"id"`
	MealLogID   uuid.UUID  `json:"meal_log_id"`
	ItemType    string     `json:"item_type"` // food, recipe
	FoodID      *uuid.UUID `json:"food_id,omitempty"`
	RecipeID    *uuid.UUID `json:"recipe_id,omitempty"`
	Quantity    float64    `json:"quantity"`
	Unit        string     `json:"unit"`
	ServingSize *float64   `json:"serving_size,omitempty"` // For recipes
	Calories    float64    `json:"calories"`
	ProteinG    float64    `json:"protein_g"`
	CarbsG      float64    `json:"carbs_g"`
	FatG        float64    `json:"fat_g"`
	Order       int        `json:"order"`
	Food        *Food      `json:"food,omitempty"`
	Recipe      *Recipe    `json:"recipe,omitempty"`
}

// CreateMealLogInput represents input for creating a meal log
type CreateMealLogInput struct {
	LogDate     time.Time              `json:"log_date" validate:"required"`
	MealTime    string                 `json:"meal_time" validate:"required,oneof=breakfast lunch dinner snack other"`
	Notes       *string                `json:"notes,omitempty" validate:"omitempty,max=1000"`
	Mood        *string                `json:"mood,omitempty" validate:"omitempty,oneof=great good okay tired"`
	EnergyLevel *int                   `json:"energy_level,omitempty" validate:"omitempty,gte=1,lte=5"`
	Items       []AddItemToMealLogInput `json:"items,omitempty"`
}

// UpdateMealLogInput represents input for updating a meal log
type UpdateMealLogInput struct {
	Notes       *string                 `json:"notes,omitempty" validate:"omitempty,max=1000"`
	Mood        *string                 `json:"mood,omitempty" validate:"omitempty,oneof=great good okay tired"`
	EnergyLevel *int                    `json:"energy_level,omitempty" validate:"omitempty,gte=1,lte=5"`
	// Items, when non-nil, replaces ALL existing items in the meal log.
	Items       []AddItemToMealLogInput `json:"items,omitempty"`
}

// DailyMealResponse represents the response for getting meal logs by date
type DailyMealResponse struct {
	Date    time.Time              `json:"date"`
	MealLog []MealLog              `json:"meal_logs"`
	Summary *DailyNutritionSummary `json:"summary"`
}

// AddItemToMealLogInput represents input for adding item to meal log
type AddItemToMealLogInput struct {
	ItemType    string     `json:"item_type" validate:"required,oneof=food recipe"`
	FoodID      *uuid.UUID `json:"food_id,omitempty"`
	RecipeID    *uuid.UUID `json:"recipe_id,omitempty"`
	Quantity    float64    `json:"quantity" validate:"required,gt=0"`
	Unit        string     `json:"unit" validate:"required"`
	ServingSize *float64   `json:"serving_size,omitempty" validate:"omitempty,gt=0"` // For recipes
	Order       int        `json:"order" validate:"required,gte=1"`
}

// UpdateItemInMealLogInput represents input for updating item quantity
type UpdateItemInMealLogInput struct {
	Quantity    *float64 `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Unit        *string  `json:"unit,omitempty"`
	ServingSize *float64 `json:"serving_size,omitempty" validate:"omitempty,gt=0"`
	Order       *int     `json:"order,omitempty" validate:"omitempty,gte=1"`
}

// GetMealLogsFilter represents filters for getting meal logs
type GetMealLogsFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	MealTime  *string
	Page      int
	PageSize  int
}

// DailyNutritionSummary represents daily nutrition summary
type DailyNutritionSummary struct {
	Date              time.Time `json:"date"`
	TotalCalories     float64   `json:"total_calories"`
	TotalProteinG     float64   `json:"total_protein_g"`
	TotalCarbsG       float64   `json:"total_carbs_g"`
	TotalFatG         float64   `json:"total_fat_g"`
	CalorieTarget     int       `json:"calorie_target"`
	AdherencePercent  float64   `json:"adherence_percent"`
	MealsLogged       int       `json:"meals_logged"`
}

// NutritionStats represents nutrition statistics
type NutritionStats struct {
	Period                string  `json:"period"` // daily, weekly, monthly
	AverageCalories       float64 `json:"average_calories"`
	AverageProteinG       float64 `json:"average_protein_g"`
	AverageCarbsG         float64 `json:"average_carbs_g"`
	AverageFatG           float64 `json:"average_fat_g"`
	AverageAdherencePercent float64 `json:"average_adherence_percent"`
	TotalMealsLogged      int     `json:"total_meals_logged"`
	DaysTracked           int     `json:"days_tracked"`
}
