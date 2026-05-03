package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/internal/infrastructure/database"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

// FoodRepository implementation
type foodRepository struct {
	db *database.DB
}

func NewFoodRepository(db *database.DB) meal.FoodRepository {
	return &foodRepository{db: db}
}

func (r *foodRepository) Create(ctx context.Context, food *meal.Food) error {
	query := `
		INSERT INTO foods (
			id, name, description, brand, image_url, barcode, serving_size, unit, calories, 
			protein_g, carbs_g, fat_g, fiber_g, is_system, created_by_user_id, 
			category, embedding, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`
	_, err := r.db.Exec(ctx, query,
		food.ID, food.Name, food.Description, food.Brand, food.ImageUrl, food.Barcode, food.ServingSize, food.Unit, food.Calories,
		food.ProteinG, food.CarbsG, food.FatG, food.FiberG, food.IsSystem, food.CreatedByUserID,
		food.Category, food.Embedding, food.CreatedAt, food.UpdatedAt,
	)

	fmt.Println(err)

	return err
}

func (r *foodRepository) GetByID(ctx context.Context, id uuid.UUID) (*meal.Food, error) {
	query := `
		SELECT id, name, description, brand, image_url, barcode, serving_size, unit, calories, 
			protein_g, carbs_g, fat_g, fiber_g, is_system, created_by_user_id, 
			category, created_at, updated_at
		FROM foods WHERE id = $1
	`
	var food meal.Food
	err := r.db.QueryRow(ctx, query, id).Scan(
		&food.ID, &food.Name, &food.Description, &food.Brand, &food.ImageUrl, &food.Barcode, &food.ServingSize, &food.Unit, &food.Calories,
		&food.ProteinG, &food.CarbsG, &food.FatG, &food.FiberG, &food.IsSystem, &food.CreatedByUserID,
		&food.Category, &food.CreatedAt, &food.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *foodRepository) GetByBarcode(ctx context.Context, barcode string) (*meal.Food, error) {
	query := `
		SELECT id, name, description, brand, image_url, barcode, serving_size, unit, calories, 
			protein_g, carbs_g, fat_g, fiber_g, is_system, created_by_user_id, 
			category, created_at, updated_at
		FROM foods WHERE barcode = $1
	`
	var food meal.Food
	err := r.db.QueryRow(ctx, query, barcode).Scan(
		&food.ID, &food.Name, &food.Description, &food.Brand, &food.ImageUrl, &food.Barcode, &food.ServingSize, &food.Unit, &food.Calories,
		&food.ProteinG, &food.CarbsG, &food.FatG, &food.FiberG, &food.IsSystem, &food.CreatedByUserID,
		&food.Category, &food.CreatedAt, &food.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *foodRepository) List(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]meal.Food, int64, error) {
	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM foods WHERE is_system = true OR created_by_user_id = $1`
	var total int64
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, description, brand, image_url, barcode, serving_size, unit, calories, 
			protein_g, carbs_g, fat_g, fiber_g, is_system, created_by_user_id, 
			category, created_at, updated_at
		FROM foods 
		WHERE is_system = true OR created_by_user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var foods []meal.Food
	for rows.Next() {
		var food meal.Food
		err := rows.Scan(
			&food.ID, &food.Name, &food.Description, &food.Brand, &food.ImageUrl, &food.Barcode, &food.ServingSize, &food.Unit, &food.Calories,
			&food.ProteinG, &food.CarbsG, &food.FatG, &food.FiberG, &food.IsSystem, &food.CreatedByUserID,
			&food.Category, &food.CreatedAt, &food.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		foods = append(foods, food)
	}
	return foods, total, nil
}

func (r *foodRepository) Search(ctx context.Context, filter meal.SearchFoodsFilter) ([]meal.Food, int64, error) {
	offset := (filter.Page - 1) * filter.PageSize

	baseWhere := `WHERE (is_system = true`
	args := []interface{}{}
	argCount := 1

	if filter.UserID != nil {
		baseWhere += ` OR created_by_user_id = $` + fmt.Sprint(argCount) + `)`
		args = append(args, *filter.UserID)
		argCount++
	} else {
		baseWhere += `)`
	}

	if filter.Query != nil {
		baseWhere += ` AND name ILIKE $` + fmt.Sprint(argCount)
		args = append(args, "%"+*filter.Query+"%")
		argCount++
	}

	if filter.Category != nil {
		baseWhere += ` AND category = $` + fmt.Sprint(argCount)
		args = append(args, *filter.Category)
		argCount++
	}

	if filter.IsSystem != nil {
		baseWhere += ` AND is_system = $` + fmt.Sprint(argCount)
		args = append(args, *filter.IsSystem)
		argCount++
	}

	countQuery := `SELECT COUNT(*) FROM foods ` + baseWhere
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, description, brand, image_url, barcode, serving_size, unit, calories, 
			protein_g, carbs_g, fat_g, fiber_g, is_system, created_by_user_id, 
			category, created_at, updated_at
		FROM foods 
		` + baseWhere + `
		ORDER BY created_at DESC
		LIMIT $` + fmt.Sprint(argCount) + ` OFFSET $` + fmt.Sprint(argCount+1)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var foods []meal.Food
	for rows.Next() {
		var food meal.Food
		err := rows.Scan(
			&food.ID, &food.Name, &food.Description, &food.Brand, &food.ImageUrl, &food.Barcode, &food.ServingSize, &food.Unit, &food.Calories,
			&food.ProteinG, &food.CarbsG, &food.FatG, &food.FiberG, &food.IsSystem, &food.CreatedByUserID,
			&food.Category, &food.CreatedAt, &food.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		foods = append(foods, food)
	}
	return foods, total, nil
}

func (r *foodRepository) Update(ctx context.Context, id uuid.UUID, input meal.UpdateFoodInput) error {
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
	if input.Brand != nil {
		setClauses = append(setClauses, fmt.Sprintf("brand = $%d", argID))
		args = append(args, *input.Brand)
		argID++
	}
	if input.ImageUrl != nil {
		setClauses = append(setClauses, fmt.Sprintf("image_url = $%d", argID))
		args = append(args, *input.ImageUrl)
		argID++
	}
	if input.Barcode != nil {
		setClauses = append(setClauses, fmt.Sprintf("barcode = $%d", argID))
		args = append(args, *input.Barcode)
		argID++
	}
	if input.ServingSize != nil {
		setClauses = append(setClauses, fmt.Sprintf("serving_size = $%d", argID))
		args = append(args, *input.ServingSize)
		argID++
	}
	if input.Unit != nil {
		setClauses = append(setClauses, fmt.Sprintf("unit = $%d", argID))
		args = append(args, *input.Unit)
		argID++
	}
	if input.Calories != nil {
		setClauses = append(setClauses, fmt.Sprintf("calories = $%d", argID))
		args = append(args, *input.Calories)
		argID++
	}
	if input.ProteinG != nil {
		setClauses = append(setClauses, fmt.Sprintf("protein_g = $%d", argID))
		args = append(args, *input.ProteinG)
		argID++
	}
	if input.CarbsG != nil {
		setClauses = append(setClauses, fmt.Sprintf("carbs_g = $%d", argID))
		args = append(args, *input.CarbsG)
		argID++
	}
	if input.FatG != nil {
		setClauses = append(setClauses, fmt.Sprintf("fat_g = $%d", argID))
		args = append(args, *input.FatG)
		argID++
	}
	if input.FiberG != nil {
		setClauses = append(setClauses, fmt.Sprintf("fiber_g = $%d", argID))
		args = append(args, *input.FiberG)
		argID++
	}
	if input.Category != nil {
		setClauses = append(setClauses, fmt.Sprintf("category = $%d", argID))
		args = append(args, *input.Category)
		argID++
	}

	if len(setClauses) == 1 {
		return nil // Nothing to update
	}

	query := fmt.Sprintf("UPDATE foods SET %s WHERE id = $1", strings.Join(setClauses, ", "))
	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *foodRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM foods WHERE id = $1 AND created_by_user_id = $2 AND is_system = false`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("food not found or you don't have permission to delete it")
	}
	return nil
}

func (r *foodRepository) GetAllFoods(ctx context.Context) ([]meal.Food, error) {
	query := `
		SELECT id, name, description, brand, image_url, barcode, serving_size, unit, calories, 
			protein_g, carbs_g, fat_g, fiber_g, is_system, created_by_user_id, 
			category, created_at, updated_at
		FROM foods 
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foods []meal.Food
	for rows.Next() {
		var food meal.Food
		err := rows.Scan(
			&food.ID, &food.Name, &food.Description, &food.Brand, &food.ImageUrl, &food.Barcode, &food.ServingSize, &food.Unit, &food.Calories,
			&food.ProteinG, &food.CarbsG, &food.FatG, &food.FiberG, &food.IsSystem, &food.CreatedByUserID,
			&food.Category, &food.CreatedAt, &food.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}
	return foods, nil
}

func (r *foodRepository) UpdateVector(ctx context.Context, id uuid.UUID, embedding []float32) error {
	vec := pgvector.NewVector(embedding)
	query := `UPDATE foods SET embedding = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, vec, id)
	return err
}

func (r *foodRepository) SearchByVector(ctx context.Context, userID uuid.UUID, vector []float32, limit int) ([]meal.Food, error) {
	vec := pgvector.NewVector(vector)
	// Using cosine distance operator `<=>`
	query := `
		SELECT id, name, description, brand, image_url, barcode, serving_size, unit, calories, 
			protein_g, carbs_g, fat_g, fiber_g, is_system, created_by_user_id, 
			category, created_at, updated_at
		FROM foods 
		WHERE embedding IS NOT NULL AND (is_system = true OR created_by_user_id = $1)
		ORDER BY embedding <=> $2
		LIMIT $3
	`
	rows, err := r.db.Query(ctx, query, userID, vec, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foods []meal.Food
	for rows.Next() {
		var food meal.Food
		err := rows.Scan(
			&food.ID, &food.Name, &food.Description, &food.Brand, &food.ImageUrl, &food.Barcode, &food.ServingSize, &food.Unit, &food.Calories,
			&food.ProteinG, &food.CarbsG, &food.FatG, &food.FiberG, &food.IsSystem, &food.CreatedByUserID,
			&food.Category, &food.CreatedAt, &food.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}
	return foods, nil
}
