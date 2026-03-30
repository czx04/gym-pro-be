package meal

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

// FoodCategory represents the category of a food item
type FoodCategory string

const (
	CategoryProtein   FoodCategory = "protein"
	CategoryCarb      FoodCategory = "carb"
	CategoryVegetable FoodCategory = "vegetable"
	CategoryFruit     FoodCategory = "fruit"
	CategoryDairy     FoodCategory = "dairy"
	CategoryFat       FoodCategory = "fat"
	CategorySnack     FoodCategory = "snack"
	CategoryBeverage  FoodCategory = "beverage"
	CategoryOther     FoodCategory = "other"
)

// IsValidFoodCategory checks if the given category string is valid
func IsValidFoodCategory(c string) bool {
	switch FoodCategory(c) {
	case CategoryProtein, CategoryCarb, CategoryVegetable, CategoryFruit, CategoryDairy, CategoryFat, CategorySnack, CategoryBeverage, CategoryOther:
		return true
	}
	return false
}

// Food represents a food item in the library
type Food struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	Brand           *string    `json:"brand,omitempty"`
	ImageUrl        *string    `json:"image_url,omitempty"`
	Barcode         *string    `json:"barcode,omitempty"`
	ServingSize     float64    `json:"serving_size"`
	Unit            string     `json:"unit"` // g, ml, cup, tbsp, etc.
	Calories        float64    `json:"calories"`
	ProteinG        float64    `json:"protein_g"`
	CarbsG          float64    `json:"carbs_g"`
	FatG            float64    `json:"fat_g"`
	FiberG          *float64   `json:"fiber_g,omitempty"`
	IsSystem        bool       `json:"is_system"` // System food vs user custom
	CreatedByUserID *uuid.UUID       `json:"created_by_user_id,omitempty"`
	Category        *string          `json:"category,omitempty"` // protein, carb, vegetable, fruit, dairy, etc.
	Embedding       *pgvector.Vector `json:"-"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// CreateFoodInput represents input for creating a food item
type CreateFoodInput struct {
	Name        string   `json:"name" form:"name" validate:"required,min=2,max=200"`
	Description *string  `json:"description,omitempty" form:"description" validate:"omitempty,max=500"`
	Brand       *string  `json:"brand,omitempty" form:"brand" validate:"omitempty,max=100"`
	ImageUrl    *string  `json:"image_url,omitempty" form:"image_url" validate:"omitempty,url"`
	Barcode     *string  `json:"barcode,omitempty" form:"barcode" validate:"omitempty,max=50"`
	ServingSize float64  `json:"serving_size" form:"serving_size" validate:"required,gt=0"`
	Unit        string   `json:"unit" form:"unit" validate:"required,oneof=g ml kg l cup tbsp tsp oz piece serving"`
	Calories    float64  `json:"calories" form:"calories" validate:"required,gte=0"`
	ProteinG    float64  `json:"protein_g" form:"protein_g" validate:"required,gte=0"`
	CarbsG      float64  `json:"carbs_g" form:"carbs_g" validate:"required,gte=0"`
	FatG        float64  `json:"fat_g" form:"fat_g" validate:"required,gte=0"`
	FiberG      *float64 `json:"fiber_g,omitempty" form:"fiber_g" validate:"omitempty,gte=0"`
	Category    *string  `json:"category,omitempty" form:"category" validate:"omitempty,oneof=protein carb vegetable fruit dairy fat snack beverage other"`
}

// UpdateFoodInput represents input for updating a food item
type UpdateFoodInput struct {
	Name        *string  `json:"name,omitempty" form:"name" validate:"omitempty,min=2,max=200"`
	Description *string  `json:"description,omitempty" form:"description" validate:"omitempty,max=500"`
	Brand       *string  `json:"brand,omitempty" form:"brand" validate:"omitempty,max=100"`
	ImageUrl    *string  `json:"image_url,omitempty" form:"image_url" validate:"omitempty,url"`
	Barcode     *string  `json:"barcode,omitempty" form:"barcode" validate:"omitempty,max=50"`
	ServingSize *float64 `json:"serving_size,omitempty" form:"serving_size" validate:"omitempty,gt=0"`
	Unit        *string  `json:"unit,omitempty" form:"unit" validate:"omitempty,oneof=g ml kg l cup tbsp tsp oz piece serving"`
	Calories    *float64 `json:"calories,omitempty" form:"calories" validate:"omitempty,gte=0"`
	ProteinG    *float64 `json:"protein_g,omitempty" form:"protein_g" validate:"omitempty,gte=0"`
	CarbsG      *float64 `json:"carbs_g,omitempty" form:"carbs_g" validate:"omitempty,gte=0"`
	FatG        *float64 `json:"fat_g,omitempty" form:"fat_g" validate:"omitempty,gte=0"`
	FiberG      *float64 `json:"fiber_g,omitempty" form:"fiber_g" validate:"omitempty,gte=0"`
	Category    *string  `json:"category,omitempty" form:"category" validate:"omitempty,oneof=protein carb vegetable fruit dairy fat snack beverage other"`
}

// SearchFoodsFilter represents filters for searching foods
type SearchFoodsFilter struct {
	Query    *string
	Category *string
	IsSystem *bool
	UserID   *uuid.UUID // added to filter foods for a specific user + system
	Page     int
	PageSize int
}
