package admin

import (
	"time"

	"github.com/google/uuid"
)

type AdminFood struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	Brand           *string    `json:"brand,omitempty"`
	ImageUrl        *string    `json:"image_url,omitempty"`
	Barcode         *string    `json:"barcode,omitempty"`
	ServingSize     float64    `json:"serving_size"`
	Unit            string     `json:"unit"`
	Calories        float64    `json:"calories"`
	ProteinG        float64    `json:"protein_g"`
	CarbsG          float64    `json:"carbs_g"`
	FatG            float64    `json:"fat_g"`
	FiberG          *float64   `json:"fiber_g,omitempty"`
	IsSystem        bool       `json:"is_system"`
	CreatedByUserID *uuid.UUID `json:"created_by_user_id,omitempty"`
	Category        *string    `json:"category,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type CreateSystemFoodInput struct {
	Name        string   `json:"name" validate:"required,min=2,max=200"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=500"`
	Brand       *string  `json:"brand,omitempty" validate:"omitempty,max=100"`
	ImageUrl    *string  `json:"image_url,omitempty" validate:"omitempty,url"`
	Barcode     *string  `json:"barcode,omitempty" validate:"omitempty,max=50"`
	ServingSize float64  `json:"serving_size" validate:"required,gt=0"`
	Unit        string   `json:"unit" validate:"required,oneof=g ml kg l cup tbsp tsp oz piece serving"`
	Calories    float64  `json:"calories" validate:"required,gte=0"`
	ProteinG    float64  `json:"protein_g" validate:"required,gte=0"`
	CarbsG      float64  `json:"carbs_g" validate:"required,gte=0"`
	FatG        float64  `json:"fat_g" validate:"required,gte=0"`
	FiberG      *float64 `json:"fiber_g,omitempty" validate:"omitempty,gte=0"`
	Category    *string  `json:"category,omitempty" validate:"omitempty,oneof=protein carb vegetable fruit dairy fat snack beverage other"`
}

type AdminUpdateFoodInput struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=500"`
	Brand       *string  `json:"brand,omitempty" validate:"omitempty,max=100"`
	ImageUrl    *string  `json:"image_url,omitempty" validate:"omitempty,url"`
	Barcode     *string  `json:"barcode,omitempty" validate:"omitempty,max=50"`
	ServingSize *float64 `json:"serving_size,omitempty" validate:"omitempty,gt=0"`
	Unit        *string  `json:"unit,omitempty" validate:"omitempty,oneof=g ml kg l cup tbsp tsp oz piece serving"`
	Calories    *float64 `json:"calories,omitempty" validate:"omitempty,gte=0"`
	ProteinG    *float64 `json:"protein_g,omitempty" validate:"omitempty,gte=0"`
	CarbsG      *float64 `json:"carbs_g,omitempty" validate:"omitempty,gte=0"`
	FatG        *float64 `json:"fat_g,omitempty" validate:"omitempty,gte=0"`
	FiberG      *float64 `json:"fiber_g,omitempty" validate:"omitempty,gte=0"`
	Category    *string  `json:"category,omitempty" validate:"omitempty,oneof=protein carb vegetable fruit dairy fat snack beverage other"`
	IsSystem    *bool    `json:"is_system,omitempty"`
}

type ListFoodsFilter struct {
	Query    *string
	Category *string
	IsSystem *bool
	Page     int
	PageSize int
}
