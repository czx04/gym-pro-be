package meal

import (
	"context"
	"fmt"
	"time"

	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/utils"
	"gym-pro-2026-ptit/pkg/validator"

	"github.com/google/uuid"
)

// MealLogUseCases encapsulates all meal-log business logic.
type MealLogUseCases struct {
	mealLogRepo   meal.MealLogRepository
	foodRepo      meal.FoodRepository
	recipeRepo    meal.RecipeRepository
	mealDailyUC   *MealDailyUseCases
	streakUC      *MealStreakUseCases
	validator     *validator.Validator
}

func NewMealLogUseCases(
	mealLogRepo meal.MealLogRepository,
	foodRepo meal.FoodRepository,
	recipeRepo meal.RecipeRepository,
	mealDailyUC *MealDailyUseCases,
	streakUC *MealStreakUseCases,
	validator *validator.Validator,
) *MealLogUseCases {
	return &MealLogUseCases{
		mealLogRepo:   mealLogRepo,
		foodRepo:      foodRepo,
		recipeRepo:    recipeRepo,
		mealDailyUC:   mealDailyUC,
		streakUC:      streakUC,
		validator:     validator,
	}
}

// CreateMealLog creates a new meal log with optional initial items.
func (uc *MealLogUseCases) CreateMealLog(ctx context.Context, userID uuid.UUID, input meal.CreateMealLogInput) (*meal.MealLog, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	existingLogs, err := uc.mealLogRepo.GetByDate(ctx, userID, input.LogDate)
	if err != nil {
		return nil, errors.DatabaseError("failed to check existing meal logs", err)
	}
	for _, l := range existingLogs {
		if l.MealTime == input.MealTime {
			return nil, errors.Conflict(fmt.Sprintf("meal log for %s already exists on this date", input.MealTime))
		}
	}

	now := time.Now()
	log := &meal.MealLog{
		ID:        uuid.New(),
		UserID:    userID,
		LogDate:   input.LogDate,
		MealTime:  input.MealTime,
		Notes:     input.Notes,
		Mood:      input.Mood,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if input.EnergyLevel != nil {
		log.EnergyLevel = input.EnergyLevel
	}

	if err := uc.mealLogRepo.Create(ctx, log); err != nil {
		return nil, errors.DatabaseError("failed to create meal log", err)
	}

	// Persist target to meal_daily for history tracking
	_ = uc.mealDailyUC.InsertOrUpdateByUserAndDate(ctx, userID, input.LogDate)

	// Add initial items if provided.
	for _, itemInput := range input.Items {
		if err := uc.addItemToLog(ctx, log.ID, itemInput); err != nil {
			return nil, err
		}
	}

	log, err = uc.mealLogRepo.GetByID(ctx, log.ID)
	if err != nil {
		return nil, err
	}

	uc.roundMealLog(log)
	uc.notifyStreakUpdate(ctx, userID)
	return log, nil
}

// GetMealLog retrieves a meal log by ID (only the owner can access it).
func (uc *MealLogUseCases) GetMealLog(ctx context.Context, id, userID uuid.UUID) (*meal.MealLog, error) {
	log, err := uc.mealLogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NotFound("meal log not found")
	}

	if log.UserID != userID {
		return nil, errors.Forbidden("you do not have permission to view this meal log")
	}

	uc.roundMealLog(log)
	return log, nil
}

// GetMealLogsByDate returns all meal logs for a user on a specific date, with nutrition summary.
func (uc *MealLogUseCases) GetMealLogsByDate(ctx context.Context, userID uuid.UUID, date time.Time) (*meal.DailyMealResponse, error) {
	logs, err := uc.mealLogRepo.GetByDate(ctx, userID, date)
	if err != nil {
		return nil, errors.DatabaseError("failed to get meal logs", err)
	}

	summary, err := uc.mealLogRepo.GetDailySummary(ctx, userID, date)
	if err != nil {
		return nil, errors.DatabaseError("failed to get daily summary", err)
	}

	uc.roundDailySummary(summary)
	for i := range logs {
		uc.roundMealLog(&logs[i])
	}

	return &meal.DailyMealResponse{
		Date:    date,
		MealLog: logs,
		Summary: summary,
	}, nil
}

// UpdateMealLog updates notes / mood / energy_level and optionally replaces all items.
func (uc *MealLogUseCases) UpdateMealLog(ctx context.Context, id, userID uuid.UUID, input meal.UpdateMealLogInput) (*meal.MealLog, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	existing, err := uc.mealLogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NotFound("meal log not found")
	}
	if existing.UserID != userID {
		return nil, errors.Forbidden("you do not have permission to edit this meal log")
	}

	newLogDate := existing.LogDate
	if input.LogDate != nil {
		newLogDate = *input.LogDate
	}
	newMealTime := existing.MealTime
	if input.MealTime != nil {
		newMealTime = *input.MealTime
	}

	if input.LogDate != nil || input.MealTime != nil {
		// Only check if it changed to avoid conflict with itself
		if !newLogDate.Equal(existing.LogDate) || newMealTime != existing.MealTime {
			existingLogs, err := uc.mealLogRepo.GetByDate(ctx, userID, newLogDate)
			if err != nil {
				return nil, errors.DatabaseError("failed to check existing meal logs", err)
			}
			for _, l := range existingLogs {
				if l.MealTime == newMealTime && l.ID != id {
					return nil, errors.Conflict(fmt.Sprintf("meal log for %s already exists on this date", newMealTime))
				}
			}
		}
	}

	if err := uc.mealLogRepo.Update(ctx, id, input); err != nil {
		return nil, errors.DatabaseError("failed to update meal log", err)
	}

	// If a new items list is provided, replace all existing items.
	if input.Items != nil {
		for _, existingItem := range existing.Items {
			if err := uc.mealLogRepo.RemoveItem(ctx, existingItem.ID); err != nil {
				return nil, errors.DatabaseError("failed to clear existing items", err)
			}
		}
		for _, itemInput := range input.Items {
			if err := uc.addItemToLog(ctx, id, itemInput); err != nil {
				return nil, err
			}
		}
	}

	log, err := uc.mealLogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	uc.roundMealLog(log)
	uc.notifyStreakUpdate(ctx, userID)
	return log, nil
}

// DeleteMealLog deletes a meal log (items cascade).
func (uc *MealLogUseCases) DeleteMealLog(ctx context.Context, id, userID uuid.UUID) error {
	existing, err := uc.mealLogRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NotFound("meal log not found")
	}
	if existing.UserID != userID {
		return errors.Forbidden("you can only delete your own meal logs")
	}

	if err := uc.mealLogRepo.Delete(ctx, id); err != nil {
		return errors.DatabaseError("failed to delete meal log", err)
	}

	uc.notifyStreakUpdate(ctx, userID)
	return nil
}

func (uc *MealLogUseCases) notifyStreakUpdate(ctx context.Context, userID uuid.UUID) {
	if uc.streakUC == nil {
		return
	}
	_, _ = uc.streakUC.RecalculatePersistAndNotify(ctx, userID)
}

// addItemToLog validates, calculates nutrition, and persists a meal log item.
func (uc *MealLogUseCases) addItemToLog(ctx context.Context, logID uuid.UUID, input meal.AddItemToMealLogInput) error {
	if err := uc.validator.Validate(input); err != nil {
		return errors.Validation(err.Error())
	}

	item := &meal.MealLogItem{
		ID:        uuid.New(),
		MealLogID: logID,
		ItemType:  input.ItemType,
		FoodID:    input.FoodID,
		RecipeID:  input.RecipeID,
		Quantity:  input.Quantity,
		Unit:      input.Unit,
		Order:     input.Order,
	}

	switch input.ItemType {
	case "food":
		if input.FoodID == nil {
			return errors.BadRequest("food_id is required for item_type 'food'")
		}
		food, err := uc.foodRepo.GetByID(ctx, *input.FoodID)
		if err != nil || food == nil {
			return errors.BadRequest("food not found")
		}
		ratio := input.Quantity * *input.ServingSize / food.ServingSize
		item.ServingSize = input.ServingSize
		item.Unit = food.Unit
		item.Calories = food.Calories * ratio
		item.ProteinG = food.ProteinG * ratio
		item.CarbsG = food.CarbsG * ratio
		item.FatG = food.FatG * ratio

	case "recipe":
		if input.RecipeID == nil {
			return errors.BadRequest("recipe_id is required for item_type 'recipe'")
		}
		recipe, err := uc.recipeRepo.GetByID(ctx, *input.RecipeID)
		if err != nil || recipe == nil {
			return errors.BadRequest("recipe not found")
		}
		item.ServingSize = input.ServingSize

		ratio := input.Quantity * *input.ServingSize
		item.Calories = recipe.PerServingCalories * ratio
		item.ProteinG = recipe.PerServingProteinG * ratio
		item.CarbsG = recipe.PerServingCarbsG * ratio
		item.FatG = recipe.PerServingFatG * ratio

	default:
		return errors.BadRequest("invalid item_type: must be 'food' or 'recipe'")
	}

	if err := uc.mealLogRepo.AddItem(ctx, logID, item); err != nil {
		return errors.DatabaseError("failed to add item to meal log", err)
	}

	return nil
}

func (uc *MealLogUseCases) roundMealLog(log *meal.MealLog) {
	if log == nil {
		return
	}
	log.TotalCalories = utils.RoundToTwo(log.TotalCalories)
	log.TotalProteinG = utils.RoundToTwo(log.TotalProteinG)
	log.TotalCarbsG = utils.RoundToTwo(log.TotalCarbsG)
	log.TotalFatG = utils.RoundToTwo(log.TotalFatG)

	for i := range log.Items {
		item := &log.Items[i]
		item.Calories = utils.RoundToTwo(item.Calories)
		item.ProteinG = utils.RoundToTwo(item.ProteinG)
		item.CarbsG = utils.RoundToTwo(item.CarbsG)
		item.FatG = utils.RoundToTwo(item.FatG)

		if item.Food != nil {
			item.Food.Calories = utils.RoundToTwo(item.Food.Calories)
			item.Food.ProteinG = utils.RoundToTwo(item.Food.ProteinG)
			item.Food.CarbsG = utils.RoundToTwo(item.Food.CarbsG)
			item.Food.FatG = utils.RoundToTwo(item.Food.FatG)
		}

		if item.Recipe != nil {
			item.Recipe.TotalCalories = utils.RoundToTwo(item.Recipe.TotalCalories)
			item.Recipe.TotalProteinG = utils.RoundToTwo(item.Recipe.TotalProteinG)
			item.Recipe.TotalCarbsG = utils.RoundToTwo(item.Recipe.TotalCarbsG)
			item.Recipe.TotalFatG = utils.RoundToTwo(item.Recipe.TotalFatG)
			item.Recipe.PerServingCalories = utils.RoundToTwo(item.Recipe.PerServingCalories)
			item.Recipe.PerServingProteinG = utils.RoundToTwo(item.Recipe.PerServingProteinG)
			item.Recipe.PerServingCarbsG = utils.RoundToTwo(item.Recipe.PerServingCarbsG)
			item.Recipe.PerServingFatG = utils.RoundToTwo(item.Recipe.PerServingFatG)
		}
	}
}

func (uc *MealLogUseCases) roundDailySummary(s *meal.DailyNutritionSummary) {
	if s == nil {
		return
	}
	s.TotalCalories = utils.RoundToTwo(s.TotalCalories)
	s.TotalProteinG = utils.RoundToTwo(s.TotalProteinG)
	s.TotalCarbsG = utils.RoundToTwo(s.TotalCarbsG)
	s.TotalFatG = utils.RoundToTwo(s.TotalFatG)
}

// GetNutritionStats fetches nutrition stats for a single date range
func (uc *MealLogUseCases) GetNutritionStats(ctx context.Context, userID uuid.UUID, input meal.GetNutritionStatsRequest) (*meal.NutritionStats, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	start, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, errors.BadRequest("invalid start_date format, expected YYYY-MM-DD")
	}
	end, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, errors.BadRequest("invalid end_date format, expected YYYY-MM-DD")
	}

	if start.After(end) {
		return nil, errors.BadRequest("start date must be before end date")
	}

	stats, err := uc.mealLogRepo.GetStats(ctx, userID, start, end, "custom_period")
	if err != nil {
		return nil, errors.DatabaseError("failed to get stats", err)
	}

	// Round stats
	uc.roundNutritionStats(stats)

	return stats, nil
}

const maxLoggedDatesRangeDays = 400

// ListLoggedDates returns YYYY-MM-DD strings for each day in the range that has at least one meal log.
func (uc *MealLogUseCases) ListLoggedDates(ctx context.Context, userID uuid.UUID, input meal.ListLoggedDatesQuery) ([]string, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	start, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, errors.BadRequest("invalid start_date format, expected YYYY-MM-DD")
	}
	end, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, errors.BadRequest("invalid end_date format, expected YYYY-MM-DD")
	}
	if start.After(end) {
		return nil, errors.BadRequest("start date must be on or before end date")
	}
	if end.Sub(start) > maxLoggedDatesRangeDays*24*time.Hour {
		return nil, errors.BadRequest("date range too large (max 400 days)")
	}

	dates, err := uc.mealLogRepo.ListDistinctLogDates(ctx, userID, start, end)
	if err != nil {
		return nil, errors.DatabaseError("failed to list logged dates", err)
	}

	out := make([]string, 0, len(dates))
	for _, d := range dates {
		out = append(out, d.UTC().Format("2006-01-02"))
	}
	return out, nil
}

func (uc *MealLogUseCases) roundNutritionStats(s *meal.NutritionStats) {
	if s == nil {
		return
	}
	s.AverageCalories = utils.RoundToTwo(s.AverageCalories)
	s.AverageProteinG = utils.RoundToTwo(s.AverageProteinG)
	s.AverageCarbsG = utils.RoundToTwo(s.AverageCarbsG)
	s.AverageFatG = utils.RoundToTwo(s.AverageFatG)
	s.TotalCalories = utils.RoundToTwo(s.TotalCalories)
	s.TotalProteinG = utils.RoundToTwo(s.TotalProteinG)
	s.TotalCarbsG = utils.RoundToTwo(s.TotalCarbsG)
	s.TotalFatG = utils.RoundToTwo(s.TotalFatG)
	s.AverageAdherencePercent = utils.RoundToTwo(s.AverageAdherencePercent)
}
