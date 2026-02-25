package meal

import (
	"time"

	"github.com/google/uuid"
)

// Food represents a food item in the library
type Food struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	Brand           *string    `json:"brand,omitempty"`
	ServingSize     float64    `json:"serving_size"`
	Unit            string     `json:"unit"` // g, ml, cup, tbsp, etc.
	Calories        float64    `json:"calories"`
	ProteinG        float64    `json:"protein_g"`
	CarbsG          float64    `json:"carbs_g"`
	FatG            float64    `json:"fat_g"`
	FiberG          *float64   `json:"fiber_g,omitempty"`
	IsSystem        bool       `json:"is_system"` // System food vs user custom
	CreatedByUserID *uuid.UUID `json:"created_by_user_id,omitempty"`
	Category        *string    `json:"category,omitempty"` // protein, carb, vegetable, fruit, dairy, etc.
	Barcode         *string    `json:"barcode,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CreateFoodInput represents input for creating a food item
type CreateFoodInput struct {
	Name        string   `json:"name" validate:"required,min=2,max=200"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=500"`
	Brand       *string  `json:"brand,omitempty" validate:"omitempty,max=100"`
	ServingSize float64  `json:"serving_size" validate:"required,gt=0"`
	Unit        string   `json:"unit" validate:"required,oneof=g ml kg l cup tbsp tsp oz piece serving"`
	Calories    float64  `json:"calories" validate:"required,gte=0"`
	ProteinG    float64  `json:"protein_g" validate:"required,gte=0"`
	CarbsG      float64  `json:"carbs_g" validate:"required,gte=0"`
	FatG        float64  `json:"fat_g" validate:"required,gte=0"`
	FiberG      *float64 `json:"fiber_g,omitempty" validate:"omitempty,gte=0"`
	Category    *string  `json:"category,omitempty" validate:"omitempty,oneof=protein carb vegetable fruit dairy fat snack beverage other"`
	Barcode     *string  `json:"barcode,omitempty" validate:"omitempty,max=50"`
}

// UpdateFoodInput represents input for updating a food item
type UpdateFoodInput struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=500"`
	Brand       *string  `json:"brand,omitempty" validate:"omitempty,max=100"`
	ServingSize *float64 `json:"serving_size,omitempty" validate:"omitempty,gt=0"`
	Unit        *string  `json:"unit,omitempty" validate:"omitempty,oneof=g ml kg l cup tbsp tsp oz piece serving"`
	Calories    *float64 `json:"calories,omitempty" validate:"omitempty,gte=0"`
	ProteinG    *float64 `json:"protein_g,omitempty" validate:"omitempty,gte=0"`
	CarbsG      *float64 `json:"carbs_g,omitempty" validate:"omitempty,gte=0"`
	FatG        *float64 `json:"fat_g,omitempty" validate:"omitempty,gte=0"`
	FiberG      *float64 `json:"fiber_g,omitempty" validate:"omitempty,gte=0"`
	Category    *string  `json:"category,omitempty" validate:"omitempty,oneof=protein carb vegetable fruit dairy fat snack beverage other"`
}

// SearchFoodsFilter represents filters for searching foods
type SearchFoodsFilter struct {
	Query    *string
	Category *string
	IsSystem *bool
	Page     int
	PageSize int
}
