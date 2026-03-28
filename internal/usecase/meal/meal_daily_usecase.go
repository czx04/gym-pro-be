package meal

import (
	"context"
	"time"

	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/pkg/errors"

	"github.com/google/uuid"
)

// MealDailyUseCases encapsulates business logic for meal daily targets.
type MealDailyUseCases struct {
	mealDailyRepo meal.MealDailyRepository
	userRepo      user.Repository
}

func NewMealDailyUseCases(
	mealDailyRepo meal.MealDailyRepository,
	userRepo user.Repository,
) *MealDailyUseCases {
	return &MealDailyUseCases{
		mealDailyRepo: mealDailyRepo,
		userRepo:      userRepo,
	}
}

// InsertOrUpdateByUserAndDate fetches the user's current nutrition targets and inserts them
// into meal_daily if no record exists for that date.
func (uc *MealDailyUseCases) InsertOrUpdateByUserAndDate(ctx context.Context, userID uuid.UUID, logDate time.Time) error {
	now := time.Now()
	if currentUser, err := uc.userRepo.GetByID(ctx, userID); err == nil {
		var targetCalories float64
		if currentUser.DailyCalorieTarget != nil {
			targetCalories = float64(*currentUser.DailyCalorieTarget)
		}
		var targetProteinG float64
		if currentUser.ProteinTargetG != nil {
			targetProteinG = float64(*currentUser.ProteinTargetG)
		}
		var targetCarbsG float64
		if currentUser.CarbsTargetG != nil {
			targetCarbsG = float64(*currentUser.CarbsTargetG)
		}
		var targetFatG float64
		if currentUser.FatTargetG != nil {
			targetFatG = float64(*currentUser.FatTargetG)
		}

		md := &meal.MealDaily{
			ID:             uuid.New(),
			UserID:         userID,
			Date:           logDate.Truncate(24 * time.Hour),
			TargetCalories: targetCalories,
			TargetProteinG: targetProteinG,
			TargetCarbsG:   targetCarbsG,
			TargetFatG:     targetFatG,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		// InsertOrUpdate with ON CONFLICT DO NOTHING (ignore if exists)
		return uc.mealDailyRepo.InsertOrUpdate(ctx, md)
	}
	return nil
}

// GetMealDailyByDate retrieves a user's daily meal target for a specific date.
// Logic:
// 1. Check if there is an exact match for the requested date.
// 2. If not, get the most recent target before that date.
// 3. If no history is found, fallback to the current user's targets.
func (uc *MealDailyUseCases) GetMealDailyByDate(ctx context.Context, userID uuid.UUID, date time.Time) (*meal.DailyNutritionTargetResponse, error) {
	reqDate := date.Truncate(24 * time.Hour)

	// 1. Try to get exact date match
	md, err := uc.mealDailyRepo.GetByDate(ctx, userID, reqDate)
	if err == nil && md != nil {
		return &meal.DailyNutritionTargetResponse{
			DailyCalorieTarget: &md.TargetCalories,
			ProteinTargetG:     &md.TargetProteinG,
			CarbsTargetG:       &md.TargetCarbsG,
			FatTargetG:         &md.TargetFatG,
		}, nil
	}

	// 2. If not found, get the most recent before date
	latestMd, err := uc.mealDailyRepo.GetLatestBeforeDate(ctx, userID, reqDate)
	if err == nil && latestMd != nil {
		return &meal.DailyNutritionTargetResponse{
			DailyCalorieTarget: &latestMd.TargetCalories,
			ProteinTargetG:     &latestMd.TargetProteinG,
			CarbsTargetG:       &latestMd.TargetCarbsG,
			FatTargetG:         &latestMd.TargetFatG,
		}, nil
	}

	// 3. If no meal daily found before, fallback to current user targets
	currentUser, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.DatabaseError("failed to get user", err)
	}
	
	res := &meal.DailyNutritionTargetResponse{}
	if currentUser.DailyCalorieTarget != nil {
		val := float64(*currentUser.DailyCalorieTarget)
		res.DailyCalorieTarget = &val
	}
	if currentUser.ProteinTargetG != nil {
		val := float64(*currentUser.ProteinTargetG)
		res.ProteinTargetG = &val
	}
	if currentUser.CarbsTargetG != nil {
		val := float64(*currentUser.CarbsTargetG)
		res.CarbsTargetG = &val
	}
	if currentUser.FatTargetG != nil {
		val := float64(*currentUser.FatTargetG)
		res.FatTargetG = &val
	}
	
	return res, nil
}
