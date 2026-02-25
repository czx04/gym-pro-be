package meal

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// FoodRepository defines the interface for food data access
type FoodRepository interface {
	// Create creates a new food item
	Create(ctx context.Context, food *Food) error
	
	// GetByID retrieves a food by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Food, error)
	
	// List retrieves foods with pagination
	List(ctx context.Context, page, pageSize int) ([]Food, int64, error)
	
	// Search searches foods with filters
	Search(ctx context.Context, filter SearchFoodsFilter) ([]Food, int64, error)
	
	// Update updates a food item
	Update(ctx context.Context, id uuid.UUID, input UpdateFoodInput) error
	
	// Delete deletes a food item (only custom foods)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	
	// GetByBarcode retrieves a food by barcode
	GetByBarcode(ctx context.Context, barcode string) (*Food, error)
}

// RecipeRepository defines the interface for recipe data access
type RecipeRepository interface {
	// Create creates a new recipe
	Create(ctx context.Context, recipe *Recipe) error
	
	// GetByID retrieves a recipe by ID with foods
	GetByID(ctx context.Context, id uuid.UUID) (*Recipe, error)
	
	// GetByUserID retrieves recipes for a user
	GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]Recipe, int64, error)
	
	// Update updates a recipe
	Update(ctx context.Context, id uuid.UUID, input UpdateRecipeInput) error
	
	// Delete deletes a recipe
	Delete(ctx context.Context, id uuid.UUID) error
	
	// AddFood adds a food to a recipe
	AddFood(ctx context.Context, recipeID uuid.UUID, food *RecipeFood) error
	
	// UpdateFood updates a food in a recipe
	UpdateFood(ctx context.Context, recipeFoodID uuid.UUID, input UpdateFoodInRecipeInput) error
	
	// RemoveFood removes a food from a recipe
	RemoveFood(ctx context.Context, recipeFoodID uuid.UUID) error
	
	// GetFoods retrieves foods in a recipe
	GetFoods(ctx context.Context, recipeID uuid.UUID) ([]RecipeFood, error)
	
	// RecalculateNutrition recalculates recipe nutrition based on foods
	RecalculateNutrition(ctx context.Context, recipeID uuid.UUID) error
}

// MealLogRepository defines the interface for meal log data access
type MealLogRepository interface {
	// Create creates a new meal log
	Create(ctx context.Context, log *MealLog) error
	
	// GetByID retrieves a meal log by ID with items
	GetByID(ctx context.Context, id uuid.UUID) (*MealLog, error)
	
	// GetByUserID retrieves meal logs for a user with filters
	GetByUserID(ctx context.Context, userID uuid.UUID, filter GetMealLogsFilter) ([]MealLog, int64, error)
	
	// GetByDate retrieves meal logs for a specific date
	GetByDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]MealLog, error)
	
	// Update updates a meal log
	Update(ctx context.Context, id uuid.UUID, input UpdateMealLogInput) error
	
	// Delete deletes a meal log
	Delete(ctx context.Context, id uuid.UUID) error
	
	// AddItem adds an item to a meal log
	AddItem(ctx context.Context, logID uuid.UUID, item *MealLogItem) error
	
	// UpdateItem updates an item in a meal log
	UpdateItem(ctx context.Context, itemID uuid.UUID, input UpdateItemInMealLogInput) error
	
	// RemoveItem removes an item from a meal log
	RemoveItem(ctx context.Context, itemID uuid.UUID) error
	
	// GetItems retrieves items in a meal log
	GetItems(ctx context.Context, logID uuid.UUID) ([]MealLogItem, error)
	
	// RecalculateNutrition recalculates meal log nutrition based on items
	RecalculateNutrition(ctx context.Context, logID uuid.UUID) error
	
	// GetDailySummary retrieves daily nutrition summary
	GetDailySummary(ctx context.Context, userID uuid.UUID, date time.Time) (*DailyNutritionSummary, error)
	
	// GetStats retrieves nutrition statistics for a period
	GetStats(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, period string) (*NutritionStats, error)
}
