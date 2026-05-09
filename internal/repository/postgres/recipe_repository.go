package postgres

import (
	"context"
	"fmt"
	"strings"

	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/internal/infrastructure/database"

	"github.com/google/uuid"
)

type recipeRepository struct {
	db *database.DB
}

func NewRecipeRepository(db *database.DB) meal.RecipeRepository {
	return &recipeRepository{db: db}
}

func (r *recipeRepository) Create(ctx context.Context, recipe *meal.Recipe) error {
	query := `
		INSERT INTO recipes (
			id, user_id, name, description, prep_time_mins, cook_time_mins, 
			servings, instructions, image_url, is_public, visibility, 
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(ctx, query,
		recipe.ID, recipe.UserID, recipe.Name, recipe.Description, recipe.PrepTimeMins,
		recipe.CookTimeMins, recipe.Servings, recipe.Instructions, recipe.ImageURL,
		recipe.IsPublic, recipe.Visibility, recipe.CreatedAt, recipe.UpdatedAt,
	)
	return err
}

func (r *recipeRepository) GetByID(ctx context.Context, id uuid.UUID) (*meal.Recipe, error) {
	query := `
		SELECT r.id, r.user_id, r.name, r.description, r.prep_time_mins, r.cook_time_mins, 
			r.servings, r.instructions, r.image_url, r.is_public, r.visibility, r.created_at, r.updated_at,
			COALESCE(SUM(rf.calories), 0) as total_calories,
			COALESCE(SUM(rf.protein_g), 0) as total_protein_g,
			COALESCE(SUM(rf.carbs_g), 0) as total_carbs_g,
			COALESCE(SUM(rf.fat_g), 0) as total_fat_g
		FROM recipes r
		LEFT JOIN recipe_foods rf ON r.id = rf.recipe_id
		WHERE r.id = $1
		GROUP BY r.id
	`
	var rec meal.Recipe
	err := r.db.QueryRow(ctx, query, id).Scan(
		&rec.ID, &rec.UserID, &rec.Name, &rec.Description, &rec.PrepTimeMins, &rec.CookTimeMins,
		&rec.Servings, &rec.Instructions, &rec.ImageURL, &rec.IsPublic, &rec.Visibility,
		&rec.CreatedAt, &rec.UpdatedAt,
		&rec.TotalCalories, &rec.TotalProteinG, &rec.TotalCarbsG, &rec.TotalFatG,
	)
	if err != nil {
		return nil, err
	}

	if rec.Servings > 0 {
		rec.PerServingCalories = rec.TotalCalories / float64(rec.Servings)
		rec.PerServingProteinG = rec.TotalProteinG / float64(rec.Servings)
		rec.PerServingCarbsG = rec.TotalCarbsG / float64(rec.Servings)
		rec.PerServingFatG = rec.TotalFatG / float64(rec.Servings)
	}

	foods, err := r.GetFoods(ctx, id)
	if err == nil {
		rec.Foods = foods
	}

	return &rec, nil
}

func (r *recipeRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int, searchQuery string) ([]meal.Recipe, int64, error) {
	offset := (page - 1) * pageSize

	baseWhere := `WHERE (user_id = $1 OR is_public = true)`
	args := []interface{}{userID}
	argCount := 2

	if searchQuery != "" {
		baseWhere += ` AND name ILIKE $` + fmt.Sprint(argCount)
		args = append(args, "%"+searchQuery+"%")
		argCount++
	}

	countQuery := `SELECT COUNT(*) FROM recipes ` + baseWhere
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT r.id, r.user_id, r.name, r.description, r.prep_time_mins, r.cook_time_mins, 
			r.servings, r.instructions, r.image_url, r.is_public, r.visibility, r.created_at, r.updated_at,
			COALESCE(SUM(rf.calories), 0) as total_calories,
			COALESCE(SUM(rf.protein_g), 0) as total_protein_g,
			COALESCE(SUM(rf.carbs_g), 0) as total_carbs_g,
			COALESCE(SUM(rf.fat_g), 0) as total_fat_g
		FROM recipes r
		LEFT JOIN recipe_foods rf ON r.id = rf.recipe_id
		` + baseWhere + `
		GROUP BY r.id
		ORDER BY r.created_at DESC
		LIMIT $` + fmt.Sprint(argCount) + ` OFFSET $` + fmt.Sprint(argCount+1)

	args = append(args, pageSize, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var recipes []meal.Recipe
	for rows.Next() {
		var rec meal.Recipe
		err := rows.Scan(
			&rec.ID, &rec.UserID, &rec.Name, &rec.Description, &rec.PrepTimeMins, &rec.CookTimeMins,
			&rec.Servings, &rec.Instructions, &rec.ImageURL, &rec.IsPublic, &rec.Visibility,
			&rec.CreatedAt, &rec.UpdatedAt,
			&rec.TotalCalories, &rec.TotalProteinG, &rec.TotalCarbsG, &rec.TotalFatG,
		)
		if err != nil {
			return nil, 0, err
		}

		if rec.Servings > 0 {
			rec.PerServingCalories = rec.TotalCalories / float64(rec.Servings)
			rec.PerServingProteinG = rec.TotalProteinG / float64(rec.Servings)
			rec.PerServingCarbsG = rec.TotalCarbsG / float64(rec.Servings)
			rec.PerServingFatG = rec.TotalFatG / float64(rec.Servings)
		}

		recipes = append(recipes, rec)
	}

	return recipes, total, nil
}

func (r *recipeRepository) Update(ctx context.Context, id uuid.UUID, input meal.UpdateRecipeInput) error {
	setClauses := []string{"updated_at = CURRENT_TIMESTAMP"}
	args := []interface{}{id}
	argID := 2

	if input.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argID))
		args = append(args, *input.Name)
		argID++
	}
	if input.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argID))
		args = append(args, *input.Description)
		argID++
	}
	if input.PrepTimeMins != nil {
		setClauses = append(setClauses, fmt.Sprintf("prep_time_mins = $%d", argID))
		args = append(args, *input.PrepTimeMins)
		argID++
	}
	if input.CookTimeMins != nil {
		setClauses = append(setClauses, fmt.Sprintf("cook_time_mins = $%d", argID))
		args = append(args, *input.CookTimeMins)
		argID++
	}
	if input.Servings != nil {
		setClauses = append(setClauses, fmt.Sprintf("servings = $%d", argID))
		args = append(args, *input.Servings)
		argID++
	}
	if input.Instructions != nil {
		setClauses = append(setClauses, fmt.Sprintf("instructions = $%d", argID))
		args = append(args, *input.Instructions)
		argID++
	}
	if input.ImageURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("image_url = $%d", argID))
		args = append(args, *input.ImageURL)
		argID++
	}
	if input.IsPublic != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_public = $%d", argID))
		args = append(args, *input.IsPublic)
		argID++
	}
	if input.Visibility != nil {
		setClauses = append(setClauses, fmt.Sprintf("visibility = $%d", argID))
		args = append(args, *input.Visibility)
	}

	if len(setClauses) == 1 {
		return nil // Nothing to update
	}

	query := fmt.Sprintf("UPDATE recipes SET %s WHERE id = $1", strings.Join(setClauses, ", "))
	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *recipeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM recipes WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *recipeRepository) AddFood(ctx context.Context, recipeID uuid.UUID, food *meal.RecipeFood) error {
	query := `
		INSERT INTO recipe_foods (
			id, recipe_id, food_id, quantity, unit, 
			calories, protein_g, carbs_g, fat_g
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		food.ID, recipeID, food.FoodID, food.Quantity, food.Unit,
		food.Calories, food.ProteinG, food.CarbsG, food.FatG,
	)
	return err
}

func (r *recipeRepository) UpdateFood(ctx context.Context, recipeFoodID uuid.UUID, input meal.UpdateFoodInRecipeInput) error {
	setClauses := []string{}
	args := []interface{}{recipeFoodID}
	argID := 2

	if input.Quantity != nil {
		setClauses = append(setClauses, fmt.Sprintf("quantity = $%d", argID))
		args = append(args, *input.Quantity)
	}

	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE recipe_foods SET %s WHERE id = $1", strings.Join(setClauses, ", "))
	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *recipeRepository) RemoveFood(ctx context.Context, recipeFoodID uuid.UUID) error {
	query := `DELETE FROM recipe_foods WHERE id = $1`
	_, err := r.db.Exec(ctx, query, recipeFoodID)
	return err
}

func (r *recipeRepository) ClearFoods(ctx context.Context, recipeID uuid.UUID) error {
	query := `DELETE FROM recipe_foods WHERE recipe_id = $1`
	_, err := r.db.Exec(ctx, query, recipeID)
	return err
}

func (r *recipeRepository) GetFoods(ctx context.Context, recipeID uuid.UUID) ([]meal.RecipeFood, error) {
	query := `
		SELECT rf.id, rf.recipe_id, rf.food_id, rf.quantity, rf.unit, 
			rf.calories, rf.protein_g, rf.carbs_g, rf.fat_g,
			f.id, f.name, f.brand, f.image_url, f.serving_size, f.unit, f.calories, 
			f.protein_g, f.carbs_g, f.fat_g, f.category
		FROM recipe_foods rf
		JOIN foods f ON rf.food_id = f.id
		WHERE rf.recipe_id = $1
	`

	rows, err := r.db.Query(ctx, query, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foods []meal.RecipeFood
	for rows.Next() {
		var rf meal.RecipeFood
		var f meal.Food
		err := rows.Scan(
			&rf.ID, &rf.RecipeID, &rf.FoodID, &rf.Quantity, &rf.Unit,
			&rf.Calories, &rf.ProteinG, &rf.CarbsG, &rf.FatG,
			&f.ID, &f.Name, &f.Brand, &f.ImageUrl, &f.ServingSize, &f.Unit, &f.Calories,
			&f.ProteinG, &f.CarbsG, &f.FatG, &f.Category,
		)
		if err != nil {
			return nil, err
		}
		rf.Food = &f
		foods = append(foods, rf)
	}
	return foods, nil
}

func (r *recipeRepository) RecalculateNutrition(ctx context.Context, recipeID uuid.UUID) error {
	// Not storing totals dynamically, skipping actual updates
	return nil
}
