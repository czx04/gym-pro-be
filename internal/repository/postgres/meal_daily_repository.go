package postgres

import (
	"context"
	"time"

	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/internal/infrastructure/database"

	"github.com/google/uuid"
)

type mealDailyRepository struct {
	db *database.DB
}

func NewMealDailyRepository(db *database.DB) meal.MealDailyRepository {
	return &mealDailyRepository{db: db}
}

// InsertOrUpdate inserts a new meal daily log if it doesn't already exist for the user and date.
func (r *mealDailyRepository) InsertOrUpdate(ctx context.Context, md *meal.MealDaily) error {
	query := `
		INSERT INTO meal_daily (
			id, user_id, date, target_calories, target_protein_g, target_carbs_g, target_fat_g, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT ON CONSTRAINT unique_user_date DO NOTHING
	`
	_, err := r.db.Exec(ctx, query,
		md.ID, md.UserID, md.Date, md.TargetCalories, md.TargetProteinG, md.TargetCarbsG, md.TargetFatG, md.CreatedAt, md.UpdatedAt,
	)
	return err
}

// GetByDate retrieves a meal daily log for a user on a given exact date.
func (r *mealDailyRepository) GetByDate(ctx context.Context, userID uuid.UUID, date time.Time) (*meal.MealDaily, error) {
	query := `
		SELECT id, user_id, date, target_calories, target_protein_g, target_carbs_g, target_fat_g, created_at, updated_at
		FROM meal_daily
		WHERE user_id = $1 AND date = $2
	`
	var md meal.MealDaily
	err := r.db.QueryRow(ctx, query, userID, date).Scan(
		&md.ID, &md.UserID, &md.Date, &md.TargetCalories, &md.TargetProteinG, &md.TargetCarbsG, &md.TargetFatG, &md.CreatedAt, &md.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &md, nil
}

// GetLatestBeforeDate retrieves the newest meal daily log strictly before the given date.
func (r *mealDailyRepository) GetLatestBeforeDate(ctx context.Context, userID uuid.UUID, date time.Time) (*meal.MealDaily, error) {
	query := `
		SELECT id, user_id, date, target_calories, target_protein_g, target_carbs_g, target_fat_g, created_at, updated_at
		FROM meal_daily
		WHERE user_id = $1 AND date < $2
		ORDER BY date DESC
		LIMIT 1
	`
	var md meal.MealDaily
	err := r.db.QueryRow(ctx, query, userID, date).Scan(
		&md.ID, &md.UserID, &md.Date, &md.TargetCalories, &md.TargetProteinG, &md.TargetCarbsG, &md.TargetFatG, &md.CreatedAt, &md.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &md, nil
}
