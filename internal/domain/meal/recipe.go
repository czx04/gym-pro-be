package meal

import (
	"time"

	"github.com/google/uuid"
)

// Recipe represents a recipe containing multiple foods
type Recipe struct {
	ID                  uuid.UUID    `json:"id"`
	UserID              uuid.UUID    `json:"user_id"`
	Name                string       `json:"name"`
	Description         *string      `json:"description,omitempty"`
	PrepTimeMins        *int         `json:"prep_time_mins,omitempty"`
	CookTimeMins        *int         `json:"cook_time_mins,omitempty"`
	Servings            int          `json:"servings"`
	Instructions        *string      `json:"instructions,omitempty"`
	ImageURL            *string      `json:"image_url,omitempty"`
	TotalCalories       float64      `json:"total_calories"`
	TotalProteinG       float64      `json:"total_protein_g"`
	TotalCarbsG         float64      `json:"total_carbs_g"`
	TotalFatG           float64      `json:"total_fat_g"`
	PerServingCalories  float64      `json:"per_serving_calories"`
	PerServingProteinG  float64      `json:"per_serving_protein_g"`
	PerServingCarbsG    float64      `json:"per_serving_carbs_g"`
	PerServingFatG      float64      `json:"per_serving_fat_g"`
	IsPublic            bool         `json:"is_public"`
	Visibility          string       `json:"visibility"` // public, private, friends
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
	Foods               []RecipeFood `json:"foods,omitempty"`
}

// RecipeFood represents a food in a recipe with quantity
type RecipeFood struct {
	ID         uuid.UUID `json:"id"`
	RecipeID   uuid.UUID `json:"recipe_id"`
	FoodID     uuid.UUID `json:"food_id"`
	Quantity   float64   `json:"quantity"`
	Unit       string    `json:"unit"`
	Calories   float64   `json:"calories"`
	ProteinG   float64   `json:"protein_g"`
	CarbsG     float64   `json:"carbs_g"`
	FatG       float64   `json:"fat_g"`
	Food       *Food     `json:"food,omitempty"`
}

// CreateRecipeInput represents input for creating a recipe
type CreateRecipeInput struct {
	Name         string  `json:"name" validate:"required,min=2,max=200"`
	Description  *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	PrepTimeMins *int    `json:"prep_time_mins,omitempty" validate:"omitempty,gte=0,lte=1440"`
	CookTimeMins *int    `json:"cook_time_mins,omitempty" validate:"omitempty,gte=0,lte=1440"`
	Servings     int     `json:"servings" validate:"required,gte=1,lte=100"`
	Instructions *string `json:"instructions,omitempty" validate:"omitempty,max=5000"`
	ImageURL     *string `json:"image_url,omitempty" validate:"omitempty,url"`
	IsPublic     bool    `json:"is_public"`
	Visibility   string  `json:"visibility" validate:"required,oneof=public private friends"`
}

// UpdateRecipeInput represents input for updating a recipe
type UpdateRecipeInput struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Description  *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	PrepTimeMins *int    `json:"prep_time_mins,omitempty" validate:"omitempty,gte=0,lte=1440"`
	CookTimeMins *int    `json:"cook_time_mins,omitempty" validate:"omitempty,gte=0,lte=1440"`
	Servings     *int    `json:"servings,omitempty" validate:"omitempty,gte=1,lte=100"`
	Instructions *string `json:"instructions,omitempty" validate:"omitempty,max=5000"`
	ImageURL     *string `json:"image_url,omitempty" validate:"omitempty,url"`
	IsPublic     *bool   `json:"is_public,omitempty"`
	Visibility   *string `json:"visibility,omitempty" validate:"omitempty,oneof=public private friends"`
}

// AddFoodToRecipeInput represents input for adding food to recipe
type AddFoodToRecipeInput struct {
	FoodID   uuid.UUID `json:"food_id" validate:"required"`
	Quantity float64   `json:"quantity" validate:"required,gt=0"`
	Unit     string    `json:"unit" validate:"required"`
}

// UpdateFoodInRecipeInput represents input for updating food quantity
type UpdateFoodInRecipeInput struct {
	Quantity *float64 `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Unit     *string  `json:"unit,omitempty"`
}
