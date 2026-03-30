package postgres

import (
	"context"
	"time"

	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/internal/infrastructure/database"

	"github.com/google/uuid"
)

type pushTokenRepository struct {
	db *database.DB
}

func NewPushTokenRepository(db *database.DB) meal.PushTokenRepository {
	return &pushTokenRepository{db: db}
}

func (r *pushTokenRepository) Upsert(ctx context.Context, userID uuid.UUID, expoPushToken, platform string) error {
	query := `
		INSERT INTO user_push_tokens (id, user_id, expo_push_token, platform, updated_at)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4)
		ON CONFLICT (expo_push_token) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			platform = EXCLUDED.platform,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(ctx, query, userID, expoPushToken, nullIfEmpty(platform), time.Now())
	return err
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func (r *pushTokenRepository) DeleteByUserAndToken(ctx context.Context, userID uuid.UUID, expoPushToken string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM user_push_tokens WHERE user_id = $1 AND expo_push_token = $2`,
		userID, expoPushToken,
	)
	return err
}

func (r *pushTokenRepository) ListTokensUsersWithoutMealOnDate(ctx context.Context, date time.Time) ([]meal.PushTokenRow, error) {
	query := `
		SELECT u.user_id, u.expo_push_token
		FROM user_push_tokens u
		WHERE NOT EXISTS (
			SELECT 1 FROM meal_logs m
			WHERE m.user_id = u.user_id AND m.log_date = $1::date
		)
	`
	rows, err := r.db.Query(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []meal.PushTokenRow
	for rows.Next() {
		var row meal.PushTokenRow
		if err := rows.Scan(&row.UserID, &row.ExpoPushToken); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}
