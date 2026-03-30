package meal

import (
	"context"
	"time"

	"gym-pro-2026-ptit/internal/delivery/http/websocket"
	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/mealstreak"

	"github.com/google/uuid"
)

// MealStreakUseCases computes and caches meal logging streaks (Asia/Ho_Chi_Minh calendar).
type MealStreakUseCases struct {
	mealLogRepo meal.MealLogRepository
	streakRepo  meal.UserMealStreakRepository
	hub         *websocket.Hub
}

func NewMealStreakUseCases(
	mealLogRepo meal.MealLogRepository,
	streakRepo meal.UserMealStreakRepository,
	hub *websocket.Hub,
) *MealStreakUseCases {
	return &MealStreakUseCases{
		mealLogRepo: mealLogRepo,
		streakRepo:  streakRepo,
		hub:         hub,
	}
}

// RecalculateAndPersist recomputes streaks from meal_logs and upserts user_meal_streaks.
func (uc *MealStreakUseCases) RecalculateAndPersist(ctx context.Context, userID uuid.UUID) (*meal.MealStreak, error) {
	dates, err := uc.mealLogRepo.ListAllDistinctLogDateStrings(ctx, userID)
	if err != nil {
		return nil, errors.DatabaseError("failed to list meal log dates", err)
	}

	now := time.Now()
	current := mealstreak.CurrentStreak(dates, now)
	longest := mealstreak.LongestStreak(dates)

	row := &meal.MealStreak{
		UserID:        userID,
		CurrentStreak: current,
		LongestStreak: longest,
		UpdatedAt:     now,
	}
	if err := uc.streakRepo.Upsert(ctx, row); err != nil {
		return nil, errors.DatabaseError("failed to save meal streak", err)
	}
	return row, nil
}

// RecalculatePersistAndNotify runs RecalculateAndPersist and pushes a WebSocket event (best-effort).
func (uc *MealStreakUseCases) RecalculatePersistAndNotify(ctx context.Context, userID uuid.UUID) (*meal.MealStreak, error) {
	s, err := uc.RecalculateAndPersist(ctx, userID)
	if err != nil {
		return nil, err
	}
	if uc.hub != nil {
		_ = uc.hub.SendJSONToUser(userID, map[string]any{
			"type":            "meal_streak_updated",
			"current_streak":  s.CurrentStreak,
			"longest_streak":  s.LongestStreak,
		})
	}
	return s, nil
}
