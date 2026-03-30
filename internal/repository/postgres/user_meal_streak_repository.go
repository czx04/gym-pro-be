package postgres

import (
	"context"
	"time"

	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/internal/infrastructure/database"

	"github.com/google/uuid"
)

type userMealStreakRepository struct {
	db *database.DB
}

func NewUserMealStreakRepository(db *database.DB) meal.UserMealStreakRepository {
	return &userMealStreakRepository{db: db}
}

func (r *userMealStreakRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*meal.MealStreak, error) {
	query := `
		SELECT user_id, current_streak, longest_streak, updated_at
		FROM user_meal_streaks
		WHERE user_id = $1
	`
	var s meal.MealStreak
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&s.UserID, &s.CurrentStreak, &s.LongestStreak, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *userMealStreakRepository) Upsert(ctx context.Context, streak *meal.MealStreak) error {
	now := time.Now()
	if streak.UpdatedAt.IsZero() {
		streak.UpdatedAt = now
	}
	query := `
		INSERT INTO user_meal_streaks (user_id, current_streak, longest_streak, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			current_streak = EXCLUDED.current_streak,
			longest_streak = EXCLUDED.longest_streak,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(ctx, query, streak.UserID, streak.CurrentStreak, streak.LongestStreak, streak.UpdatedAt)
	return err
}
