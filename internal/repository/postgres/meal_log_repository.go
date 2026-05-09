package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/internal/infrastructure/database"

	"github.com/google/uuid"
)

type mealLogRepository struct {
	db *database.DB
}

func NewMealLogRepository(db *database.DB) meal.MealLogRepository {
	return &mealLogRepository{db: db}
}

// Create inserts a new meal log (nutrition totals start at 0).
func (r *mealLogRepository) Create(ctx context.Context, log *meal.MealLog) error {
	query := `
		INSERT INTO meal_logs (
			id, user_id, log_date, meal_time,
			total_calories, total_protein_g, total_carbs_g, total_fat_g,
			notes, mood, energy_level, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(ctx, query,
		log.ID, log.UserID, log.LogDate, log.MealTime,
		log.TotalCalories, log.TotalProteinG, log.TotalCarbsG, log.TotalFatG,
		log.Notes, log.Mood, log.EnergyLevel, log.CreatedAt, log.UpdatedAt,
	)
	return err
}

// GetByID retrieves a meal log with all its items (including food/recipe details).
func (r *mealLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*meal.MealLog, error) {
	query := `
		SELECT id, user_id, log_date, meal_time,
			total_calories, total_protein_g, total_carbs_g, total_fat_g,
			notes, mood, energy_level, created_at, updated_at
		FROM meal_logs
		WHERE id = $1
	`
	var log meal.MealLog
	err := r.db.QueryRow(ctx, query, id).Scan(
		&log.ID, &log.UserID, &log.LogDate, &log.MealTime,
		&log.TotalCalories, &log.TotalProteinG, &log.TotalCarbsG, &log.TotalFatG,
		&log.Notes, &log.Mood, &log.EnergyLevel, &log.CreatedAt, &log.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	items, err := r.GetItems(ctx, id)
	if err != nil {
		return nil, err
	}
	log.Items = items

	return &log, nil
}

// GetByUserID retrieves a user's meal logs with optional date-range / meal-time / pagination filters.
func (r *mealLogRepository) GetByUserID(ctx context.Context, userID uuid.UUID, filter meal.GetMealLogsFilter) ([]meal.MealLog, int64, error) {
	args := []interface{}{userID}
	argID := 2
	conditions := []string{"user_id = $1"}

	if filter.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("log_date >= $%d", argID))
		args = append(args, filter.StartDate)
		argID++
	}
	if filter.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("log_date <= $%d", argID))
		args = append(args, filter.EndDate)
		argID++
	}
	if filter.MealTime != nil {
		conditions = append(conditions, fmt.Sprintf("meal_time = $%d", argID))
		args = append(args, *filter.MealTime)
		argID++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int64
	countQuery := "SELECT COUNT(*) FROM meal_logs " + where
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	dataQuery := fmt.Sprintf(`
		SELECT id, user_id, log_date, meal_time,
			total_calories, total_protein_g, total_carbs_g, total_fat_g,
			notes, mood, energy_level, created_at, updated_at
		FROM meal_logs
		%s
		ORDER BY log_date DESC, meal_time
		LIMIT $%d OFFSET $%d
	`, where, argID, argID+1)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []meal.MealLog
	for rows.Next() {
		var l meal.MealLog
		if err := rows.Scan(
			&l.ID, &l.UserID, &l.LogDate, &l.MealTime,
			&l.TotalCalories, &l.TotalProteinG, &l.TotalCarbsG, &l.TotalFatG,
			&l.Notes, &l.Mood, &l.EnergyLevel, &l.CreatedAt, &l.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}

	return logs, total, nil
}

// GetByDate retrieves all meal logs for a user on a specific date, ordered by meal_time.
func (r *mealLogRepository) GetByDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]meal.MealLog, error) {
	query := `
		SELECT id, user_id, log_date, meal_time,
			total_calories, total_protein_g, total_carbs_g, total_fat_g,
			notes, mood, energy_level, created_at, updated_at
		FROM meal_logs
		WHERE user_id = $1 AND log_date = $2
		ORDER BY meal_time
	`
	rows, err := r.db.Query(ctx, query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []meal.MealLog
	for rows.Next() {
		var l meal.MealLog
		if err := rows.Scan(
			&l.ID, &l.UserID, &l.LogDate, &l.MealTime,
			&l.TotalCalories, &l.TotalProteinG, &l.TotalCarbsG, &l.TotalFatG,
			&l.Notes, &l.Mood, &l.EnergyLevel, &l.CreatedAt, &l.UpdatedAt,
		); err != nil {
			return nil, err
		}

		items, err := r.GetItems(ctx, l.ID)
		if err != nil {
			return nil, err
		}
		l.Items = items

		logs = append(logs, l)
	}

	return logs, nil
}

// ListDistinctLogDates returns each calendar day in [from, to] that has at least one meal log.
func (r *mealLogRepository) ListDistinctLogDates(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]time.Time, error) {
	query := `
		SELECT DISTINCT log_date
		FROM meal_logs
		WHERE user_id = $1 AND log_date >= $2::date AND log_date <= $3::date
		ORDER BY log_date
	`
	rows, err := r.db.Query(ctx, query, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var d time.Time
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		dates = append(dates, d.UTC().Truncate(24*time.Hour))
	}
	return dates, rows.Err()
}

// Update applies partial updates (notes, mood, energy_level) to a meal log.
func (r *mealLogRepository) Update(ctx context.Context, id uuid.UUID, input meal.UpdateMealLogInput) error {
	setClauses := []string{"updated_at = CURRENT_TIMESTAMP"}
	args := []interface{}{id}
	argID := 2

	if input.LogDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("log_date = $%d", argID))
		args = append(args, *input.LogDate)
		argID++
	}
	if input.MealTime != nil {
		setClauses = append(setClauses, fmt.Sprintf("meal_time = $%d", argID))
		args = append(args, *input.MealTime)
		argID++
	}
	if input.Notes != nil {
		setClauses = append(setClauses, fmt.Sprintf("notes = $%d", argID))
		args = append(args, *input.Notes)
		argID++
	}
	if input.Mood != nil {
		setClauses = append(setClauses, fmt.Sprintf("mood = $%d", argID))
		args = append(args, *input.Mood)
		argID++
	}
	if input.EnergyLevel != nil {
		setClauses = append(setClauses, fmt.Sprintf("energy_level = $%d", argID))
		args = append(args, *input.EnergyLevel)
	}

	if len(setClauses) == 1 {
		return nil // nothing to update
	}

	query := fmt.Sprintf("UPDATE meal_logs SET %s WHERE id = $1", strings.Join(setClauses, ", "))
	_, err := r.db.Exec(ctx, query, args...)
	return err
}

// Delete removes a meal log (items cascade via FK).
func (r *mealLogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM meal_logs WHERE id = $1", id)
	return err
}

// AddItem inserts a new item into meal_log_items and recalculates the parent log totals.
func (r *mealLogRepository) AddItem(ctx context.Context, logID uuid.UUID, item *meal.MealLogItem) error {
	query := `
		INSERT INTO meal_log_items (
			id, meal_log_id, item_type, food_id, recipe_id,
			quantity, unit, serving_size,
			calories, protein_g, carbs_g, fat_g, "order"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(ctx, query,
		item.ID, logID, item.ItemType, item.FoodID, item.RecipeID,
		item.Quantity, item.Unit, item.ServingSize,
		item.Calories, item.ProteinG, item.CarbsG, item.FatG, item.Order,
	)
	if err != nil {
		return err
	}

	return r.RecalculateNutrition(ctx, logID)
}

// UpdateItem patches quantity/unit/serving_size/order of an item and recalculates parent totals.
func (r *mealLogRepository) UpdateItem(ctx context.Context, itemID uuid.UUID, input meal.UpdateItemInMealLogInput) error {
	setClauses := []string{}
	args := []interface{}{itemID}
	argID := 2

	if input.Quantity != nil {
		setClauses = append(setClauses, fmt.Sprintf("quantity = $%d", argID))
		args = append(args, *input.Quantity)
		argID++
	}
	if input.Unit != nil {
		setClauses = append(setClauses, fmt.Sprintf("unit = $%d", argID))
		args = append(args, *input.Unit)
		argID++
	}
	if input.ServingSize != nil {
		setClauses = append(setClauses, fmt.Sprintf("serving_size = $%d", argID))
		args = append(args, *input.ServingSize)
		argID++
	}
	if input.Order != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"order" = $%d`, argID))
		args = append(args, *input.Order)
	}

	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE meal_log_items SET %s WHERE id = $1", strings.Join(setClauses, ", "))
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	// Fetch the parent log ID to recalculate totals.
	var logID uuid.UUID
	if err := r.db.QueryRow(ctx, "SELECT meal_log_id FROM meal_log_items WHERE id = $1", itemID).Scan(&logID); err != nil {
		return err
	}

	return r.RecalculateNutrition(ctx, logID)
}

// RemoveItem deletes an item and recalculates the parent log's nutrition totals.
func (r *mealLogRepository) RemoveItem(ctx context.Context, itemID uuid.UUID) error {
	var logID uuid.UUID
	if err := r.db.QueryRow(ctx, "SELECT meal_log_id FROM meal_log_items WHERE id = $1", itemID).Scan(&logID); err != nil {
		return err
	}

	if _, err := r.db.Exec(ctx, "DELETE FROM meal_log_items WHERE id = $1", itemID); err != nil {
		return err
	}

	return r.RecalculateNutrition(ctx, logID)
}

// GetItems retrieves all items in a meal log, joined with food/recipe basic info.
func (r *mealLogRepository) GetItems(ctx context.Context, logID uuid.UUID) ([]meal.MealLogItem, error) {
	query := `
		SELECT
			mli.id, mli.meal_log_id, mli.item_type,
			mli.food_id, mli.recipe_id,
			mli.quantity, mli.unit, mli.serving_size,
			mli.calories, mli.protein_g, mli.carbs_g, mli.fat_g, mli."order",
			f.id, f.name, f.description, f.brand, f.image_url, f.barcode, f.serving_size, f.unit,
			f.calories, f.protein_g, f.carbs_g, f.fat_g, f.fiber_g, f.is_system, f.created_by_user_id,
			f.category, f.created_at, f.updated_at,
			r.id, r.user_id, r.name, r.description, r.prep_time_mins, r.cook_time_mins, r.servings,
			r.instructions, r.image_url, r.total_calories, r.total_protein_g, r.total_carbs_g, r.total_fat_g,
			r.per_serving_calories, r.per_serving_protein_g, r.per_serving_carbs_g, r.per_serving_fat_g,
			r.is_public, r.visibility, r.created_at, r.updated_at
		FROM meal_log_items mli
		LEFT JOIN foods    f ON mli.item_type = 'food'   AND mli.food_id   = f.id
		LEFT JOIN recipes  r ON mli.item_type = 'recipe' AND mli.recipe_id = r.id
		WHERE mli.meal_log_id = $1
		ORDER BY mli."order"
	`
	rows, err := r.db.Query(ctx, query, logID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []meal.MealLogItem
	for rows.Next() {
		var item meal.MealLogItem

		// Nullable food columns
		var (
			fID              *uuid.UUID
			fName            *string
			fDescription     *string
			fBrand           *string
			fImageURL        *string
			fBarcode         *string
			fServingSize     *float64
			fUnit            *string
			fCalories        *float64
			fProteinG        *float64
			fCarbsG          *float64
			fFatG            *float64
			fFiberG          *float64
			fIsSystem        *bool
			fCreatedByUserID *uuid.UUID
			fCategory        *string
			fCreatedAt       *time.Time
			fUpdatedAt       *time.Time
		)

		// Nullable recipe columns
		var (
			rID              *uuid.UUID
			rUserID          *uuid.UUID
			rName            *string
			rDescription     *string
			rPrepTimeMins    *int
			rCookTimeMins    *int
			rServings        *int
			rInstructions    *string
			rImageURL        *string
			rTotalCal        *float64
			rTotalPro        *float64
			rTotalCarbs      *float64
			rTotalFat        *float64
			rPerServingCal   *float64
			rPerServingPro   *float64
			rPerServingCarbs *float64
			rPerServingFat   *float64
			rIsPublic        *bool
			rVisibility      *string
			rCreatedAt       *time.Time
			rUpdatedAt       *time.Time
		)

		if err := rows.Scan(
			&item.ID, &item.MealLogID, &item.ItemType,
			&item.FoodID, &item.RecipeID,
			&item.Quantity, &item.Unit, &item.ServingSize,
			&item.Calories, &item.ProteinG, &item.CarbsG, &item.FatG, &item.Order,
			&fID, &fName, &fDescription, &fBrand, &fImageURL, &fBarcode, &fServingSize, &fUnit,
			&fCalories, &fProteinG, &fCarbsG, &fFatG, &fFiberG, &fIsSystem, &fCreatedByUserID,
			&fCategory, &fCreatedAt, &fUpdatedAt,
			&rID, &rUserID, &rName, &rDescription, &rPrepTimeMins, &rCookTimeMins, &rServings,
			&rInstructions, &rImageURL, &rTotalCal, &rTotalPro, &rTotalCarbs, &rTotalFat,
			&rPerServingCal, &rPerServingPro, &rPerServingCarbs, &rPerServingFat,
			&rIsPublic, &rVisibility, &rCreatedAt, &rUpdatedAt,
		); err != nil {
			return nil, err
		}

		if fID != nil {
			item.Food = &meal.Food{
				ID:              *fID,
				Name:            *fName,
				Description:     fDescription,
				Brand:           fBrand,
				ImageUrl:        fImageURL,
				Barcode:         fBarcode,
				ServingSize:     *fServingSize,
				Unit:            *fUnit,
				Calories:        *fCalories,
				ProteinG:        *fProteinG,
				CarbsG:          *fCarbsG,
				FatG:            *fFatG,
				FiberG:          fFiberG,
				IsSystem:        *fIsSystem,
				CreatedByUserID: fCreatedByUserID,
				Category:        fCategory,
				CreatedAt:       *fCreatedAt,
				UpdatedAt:       *fUpdatedAt,
			}
		}

		if rID != nil {
			item.Recipe = &meal.Recipe{
				ID:                 *rID,
				UserID:             *rUserID,
				Name:               *rName,
				Description:        rDescription,
				PrepTimeMins:       rPrepTimeMins,
				CookTimeMins:       rCookTimeMins,
				Servings:           *rServings,
				Instructions:       rInstructions,
				ImageURL:           rImageURL,
				TotalCalories:      derefFloat64(rTotalCal),
				TotalProteinG:      derefFloat64(rTotalPro),
				TotalCarbsG:        derefFloat64(rTotalCarbs),
				TotalFatG:          derefFloat64(rTotalFat),
				PerServingCalories: derefFloat64(rPerServingCal),
				PerServingProteinG: derefFloat64(rPerServingPro),
				PerServingCarbsG:   derefFloat64(rPerServingCarbs),
				PerServingFatG:     derefFloat64(rPerServingFat),
				IsPublic:           *rIsPublic,
				Visibility:         *rVisibility,
				CreatedAt:          *rCreatedAt,
				UpdatedAt:          *rUpdatedAt,
			}
		}

		items = append(items, item)
	}

	return items, nil
}

// RecalculateNutrition sums all items in a meal log and persists to meal_logs.
func (r *mealLogRepository) RecalculateNutrition(ctx context.Context, logID uuid.UUID) error {
	query := `
		UPDATE meal_logs SET
			total_calories  = (SELECT COALESCE(SUM(calories),  0) FROM meal_log_items WHERE meal_log_id = $1),
			total_protein_g = (SELECT COALESCE(SUM(protein_g), 0) FROM meal_log_items WHERE meal_log_id = $1),
			total_carbs_g   = (SELECT COALESCE(SUM(carbs_g),   0) FROM meal_log_items WHERE meal_log_id = $1),
			total_fat_g     = (SELECT COALESCE(SUM(fat_g),     0) FROM meal_log_items WHERE meal_log_id = $1),
			updated_at      = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, logID)
	return err
}

// GetDailySummary returns aggregated nutrition for a user on a given date.
func (r *mealLogRepository) GetDailySummary(ctx context.Context, userID uuid.UUID, date time.Time) (*meal.DailyNutritionSummary, error) {
	query := `
		SELECT
			COALESCE(SUM(total_calories),  0),
			COALESCE(SUM(total_protein_g), 0),
			COALESCE(SUM(total_carbs_g),   0),
			COALESCE(SUM(total_fat_g),     0),
			COUNT(*)
		FROM meal_logs
		WHERE user_id = $1 AND log_date = $2
	`
	var summary meal.DailyNutritionSummary
	summary.Date = date

	err := r.db.QueryRow(ctx, query, userID, date).Scan(
		&summary.TotalCalories,
		&summary.TotalProteinG,
		&summary.TotalCarbsG,
		&summary.TotalFatG,
		&summary.MealsLogged,
	)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

// GetStats returns average nutrition statistics for a given period.
func (r *mealLogRepository) GetStats(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, period string) (*meal.NutritionStats, error) {
	query := `
		SELECT
			COALESCE(AVG(daily_cal),  0),
			COALESCE(AVG(daily_pro),  0),
			COALESCE(AVG(daily_carb), 0),
			COALESCE(AVG(daily_fat),  0),
			COALESCE(SUM(daily_cal),  0),
			COALESCE(SUM(daily_pro),  0),
			COALESCE(SUM(daily_carb), 0),
			COALESCE(SUM(daily_fat),  0),
			COALESCE(SUM(meal_count), 0),
			COUNT(DISTINCT log_date)
		FROM (
			SELECT
				log_date,
				SUM(total_calories)  AS daily_cal,
				SUM(total_protein_g) AS daily_pro,
				SUM(total_carbs_g)   AS daily_carb,
				SUM(total_fat_g)     AS daily_fat,
				COUNT(*)             AS meal_count
			FROM meal_logs
			WHERE user_id = $1 AND log_date BETWEEN $2 AND $3
			GROUP BY log_date
		) daily
	`
	var stats meal.NutritionStats
	stats.Period = period

	err := r.db.QueryRow(ctx, query, userID, startDate, endDate).Scan(
		&stats.AverageCalories,
		&stats.AverageProteinG,
		&stats.AverageCarbsG,
		&stats.AverageFatG,
		&stats.TotalCalories,
		&stats.TotalProteinG,
		&stats.TotalCarbsG,
		&stats.TotalFatG,
		&stats.TotalMealsLogged,
		&stats.DaysTracked,
	)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// ListAllDistinctLogDateStrings returns sorted YYYY-MM-DD from meal_logs.
// Only dates where the log was created on the same VN calendar day are included
// so that retroactive logging does not inflate streaks.
func (r *mealLogRepository) ListAllDistinctLogDateStrings(ctx context.Context, userID uuid.UUID) ([]string, error) {
	query := `
		SELECT to_char(log_date, 'YYYY-MM-DD')
		FROM meal_logs
		WHERE user_id = $1
		  AND (created_at AT TIME ZONE 'Asia/Ho_Chi_Minh')::date = log_date
		GROUP BY log_date
		ORDER BY log_date
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// HasMealLogOnDate returns true if there is at least one meal log on that date.
func (r *mealLogRepository) HasMealLogOnDate(ctx context.Context, userID uuid.UUID, date time.Time) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM meal_logs WHERE user_id = $1 AND log_date = $2::date)`
	var ok bool
	err := r.db.QueryRow(ctx, query, userID, date).Scan(&ok)
	return ok, err
}

// derefFloat64 safely dereferences a *float64, returning 0 if nil.
func derefFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
