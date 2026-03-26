package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gym-pro-2026-ptit/internal/domain/admin"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type adminUserRepository struct {
	db *database.DB
}

func NewAdminUserRepository(db *database.DB) admin.UserRepository {
	return &adminUserRepository{db: db}
}

func (r *adminUserRepository) ListUsers(ctx context.Context, filter admin.ListUsersFilter) ([]admin.UserSummary, int64, error) {
	baseWhere := " WHERE 1=1"
	args := make([]interface{}, 0)

	if filter.Query != nil && *filter.Query != "" {
		baseWhere += ` AND (name ILIKE '%' || $` + strconv.Itoa(len(args)+1) + ` || '%' OR email ILIKE '%' || $` + strconv.Itoa(len(args)+1) + ` || '%')`
		args = append(args, *filter.Query)
	}
	if filter.Gender != nil && *filter.Gender != "" {
		baseWhere += ` AND gender = $` + strconv.Itoa(len(args)+1)
		args = append(args, *filter.Gender)
	}
	if filter.IsActive != nil {
		if *filter.IsActive {
			baseWhere += ` AND is_active = TRUE`
		} else {
			baseWhere += ` AND is_active = FALSE`
		}
	}

	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`+baseWhere, args...).Scan(&total); err != nil {
		return nil, 0, errors.DatabaseError("count users", err)
	}

	dataQuery := `SELECT id, email, name, avatar_url, gender, fitness_goal, activity_level, oauth_provider, is_active, created_at, updated_at
		FROM users` + baseWhere +
		` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, errors.DatabaseError("list users", err)
	}
	defer rows.Close()

	users := make([]admin.UserSummary, 0)
	for rows.Next() {
		var u admin.UserSummary
		if err := rows.Scan(
			&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.Gender,
			&u.FitnessGoal, &u.ActivityLevel, &u.OAuthProvider,
			&u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, 0, errors.DatabaseError("scan user", err)
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *adminUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*admin.UserDetail, error) {
	query := `
		SELECT id, email, name, bio, avatar_url, date_of_birth, gender,
			height_cm, weight_kg, fitness_goal, activity_level,
			daily_calorie_target, protein_target_g, carbs_target_g, fat_target_g,
			oauth_provider, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`
	var u admin.UserDetail
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.Name, &u.Bio, &u.AvatarURL,
		&u.DateOfBirth, &u.Gender, &u.HeightCm, &u.WeightKg,
		&u.FitnessGoal, &u.ActivityLevel,
		&u.DailyCalorieTarget, &u.ProteinTargetG, &u.CarbsTargetG, &u.FatTargetG,
		&u.OAuthProvider, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("user")
		}
		return nil, errors.DatabaseError("get user by id", err)
	}
	return &u, nil
}

func (r *adminUserRepository) UpdateUserStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	result, err := r.db.Exec(ctx, `UPDATE users SET is_active = $2, updated_at = NOW() WHERE id = $1`, id, isActive)
	if err != nil {
		return errors.DatabaseError("update user status", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound("user")
	}
	return nil
}

func (r *adminUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return errors.DatabaseError("delete user", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound("user")
	}
	return nil
}

func (r *adminUserRepository) GetOverviewStats(ctx context.Context) (*admin.OverviewStats, error) {
	stats := &admin.OverviewStats{}
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	queries := []struct {
		query  string
		args   []interface{}
		target *int64
	}{
		{`SELECT COUNT(*) FROM users`, nil, &stats.TotalUsers},
		{`SELECT COUNT(*) FROM users WHERE created_at >= $1`, []interface{}{thirtyDaysAgo}, &stats.NewUsersLast30d},
		{`SELECT COUNT(*) FROM exercises`, nil, &stats.TotalExercises},
		{`SELECT COUNT(*) FROM exercises WHERE is_active = TRUE`, nil, &stats.ActiveExercises},
		{`SELECT COUNT(*) FROM foods`, nil, &stats.TotalFoods},
		{`SELECT COUNT(*) FROM foods WHERE is_system = TRUE`, nil, &stats.SystemFoods},
		{`SELECT COUNT(*) FROM foods WHERE is_system = FALSE`, nil, &stats.UserFoods},
	}

	for _, q := range queries {
		if err := r.db.QueryRow(ctx, q.query, q.args...).Scan(q.target); err != nil {
			return nil, errors.DatabaseError("get overview stats", err)
		}
	}
	return stats, nil
}

type adminExerciseRepository struct {
	db *database.DB
}

func NewAdminExerciseRepository(db *database.DB) admin.ExerciseRepository {
	return &adminExerciseRepository{db: db}
}

func (r *adminExerciseRepository) ListExercises(ctx context.Context, filter admin.ListExercisesFilter) ([]admin.AdminExercise, int64, error) {
	baseWhere := " WHERE 1=1"
	args := make([]interface{}, 0)

	if filter.Query != nil && *filter.Query != "" {
		baseWhere += ` AND name ILIKE '%' || $` + strconv.Itoa(len(args)+1) + ` || '%'`
		args = append(args, *filter.Query)
	}
	if filter.Category != nil && *filter.Category != "" {
		baseWhere += ` AND category = $` + strconv.Itoa(len(args)+1)
		args = append(args, *filter.Category)
	}
	if filter.MuscleGroup != nil && *filter.MuscleGroup != "" {
		baseWhere += ` AND muscle_groups @> $` + strconv.Itoa(len(args)+1)
		args = append(args, []string{*filter.MuscleGroup})
	}
	if filter.DifficultyLevel != nil && *filter.DifficultyLevel != "" {
		baseWhere += ` AND difficulty_level = $` + strconv.Itoa(len(args)+1)
		args = append(args, *filter.DifficultyLevel)
	}
	if filter.IsActive != nil {
		if *filter.IsActive {
			baseWhere += ` AND is_active = TRUE`
		} else {
			baseWhere += ` AND is_active = FALSE`
		}
	}

	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM exercises`+baseWhere, args...).Scan(&total); err != nil {
		return nil, 0, errors.DatabaseError("count exercises", err)
	}

	dataQuery := `SELECT id, name, description, category, muscle_groups, equipment_needed,
		difficulty_level, calories_per_minute, video_url, thumbnail_url, is_active, created_by, created_at, updated_at
		FROM exercises` + baseWhere +
		` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, errors.DatabaseError("list exercises", err)
	}
	defer rows.Close()

	exercises := make([]admin.AdminExercise, 0)
	for rows.Next() {
		e, err := scanAdminExercise(rows)
		if err != nil {
			return nil, 0, err
		}
		exercises = append(exercises, *e)
	}
	return exercises, total, nil
}

func (r *adminExerciseRepository) GetExerciseByID(ctx context.Context, id uuid.UUID) (*admin.AdminExercise, error) {
	query := `SELECT id, name, description, category, muscle_groups, equipment_needed,
		difficulty_level, calories_per_minute, video_url, thumbnail_url, is_active, created_by, created_at, updated_at
		FROM exercises WHERE id = $1`

	e, err := scanAdminExercise(r.db.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("exercise")
		}
		return nil, err
	}
	return e, nil
}

func (r *adminExerciseRepository) CreateExercise(ctx context.Context, input admin.CreateExerciseInput, createdBy uuid.UUID) (*admin.AdminExercise, error) {
	id := uuid.New()
	now := time.Now()
	_, err := r.db.Exec(ctx, `
		INSERT INTO exercises (id, name, description, category, muscle_groups, equipment_needed,
			difficulty_level, calories_per_minute, video_url, thumbnail_url, is_active, created_by, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		id, input.Name, input.Description, input.Category, input.MuscleGroups, input.EquipmentNeeded,
		input.DifficultyLevel, input.CaloriesPerMinute, input.VideoURL, input.ThumbnailURL,
		input.IsActive, createdBy, now, now,
	)
	if err != nil {
		return nil, errors.DatabaseError("create exercise", err)
	}
	return r.GetExerciseByID(ctx, id)
}

func (r *adminExerciseRepository) UpdateExercise(ctx context.Context, id uuid.UUID, input admin.UpdateExerciseInput) (*admin.AdminExercise, error) {
	setClauses := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIdx := 1

	addField := func(col string, val interface{}) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, argIdx))
		args = append(args, val)
		argIdx++
	}

	if input.Name != nil {
		addField("name", *input.Name)
	}
	if input.Description != nil {
		addField("description", *input.Description)
	}
	if input.Category != nil {
		addField("category", *input.Category)
	}
	if input.MuscleGroups != nil {
		addField("muscle_groups", input.MuscleGroups)
	}
	if input.EquipmentNeeded != nil {
		addField("equipment_needed", input.EquipmentNeeded)
	}
	if input.DifficultyLevel != nil {
		addField("difficulty_level", *input.DifficultyLevel)
	}
	if input.CaloriesPerMinute != nil {
		addField("calories_per_minute", *input.CaloriesPerMinute)
	}
	if input.VideoURL != nil {
		addField("video_url", *input.VideoURL)
	}
	if input.ThumbnailURL != nil {
		addField("thumbnail_url", *input.ThumbnailURL)
	}
	if input.IsActive != nil {
		addField("is_active", *input.IsActive)
	}

	query := "UPDATE exercises SET "
	for i, c := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += c
	}
	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, errors.DatabaseError("update exercise", err)
	}
	if result.RowsAffected() == 0 {
		return nil, errors.NotFound("exercise")
	}
	return r.GetExerciseByID(ctx, id)
}

func (r *adminExerciseRepository) DeleteExercise(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM exercises WHERE id = $1`, id)
	if err != nil {
		return errors.DatabaseError("delete exercise", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound("exercise")
	}
	return nil
}

func scanAdminExercise(row interface {
	Scan(dest ...interface{}) error
}) (*admin.AdminExercise, error) {
	var e admin.AdminExercise
	err := row.Scan(
		&e.ID, &e.Name, &e.Description, &e.Category,
		&e.MuscleGroups, &e.EquipmentNeeded, &e.DifficultyLevel,
		&e.CaloriesPerMinute, &e.VideoURL, &e.ThumbnailURL,
		&e.IsActive, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, errors.DatabaseError("scan exercise", err)
	}
	return &e, nil
}

type adminFoodRepository struct {
	db *database.DB
}

func NewAdminFoodRepository(db *database.DB) admin.FoodRepository {
	return &adminFoodRepository{db: db}
}

func (r *adminFoodRepository) ListFoods(ctx context.Context, filter admin.ListFoodsFilter) ([]admin.AdminFood, int64, error) {
	baseWhere := " WHERE 1=1"
	args := make([]interface{}, 0)

	if filter.Query != nil && *filter.Query != "" {
		baseWhere += ` AND (name ILIKE '%' || $` + strconv.Itoa(len(args)+1) + ` || '%'` +
			` OR brand ILIKE '%' || $` + strconv.Itoa(len(args)+1) + ` || '%')`
		args = append(args, *filter.Query)
	}
	if filter.Category != nil && *filter.Category != "" {
		baseWhere += ` AND category = $` + strconv.Itoa(len(args)+1)
		args = append(args, *filter.Category)
	}
	if filter.IsSystem != nil {
		if *filter.IsSystem {
			baseWhere += ` AND is_system = TRUE`
		} else {
			baseWhere += ` AND is_system = FALSE`
		}
	}

	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM foods`+baseWhere, args...).Scan(&total); err != nil {
		return nil, 0, errors.DatabaseError("count foods", err)
	}

	dataQuery := `SELECT id, name, description, brand, image_url, barcode, serving_size, unit,
		calories, protein_g, carbs_g, fat_g, fiber_g, is_system, created_by_user_id, category, created_at, updated_at
		FROM foods` + baseWhere +
		` ORDER BY is_system DESC, created_at DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, errors.DatabaseError("list foods", err)
	}
	defer rows.Close()

	foods := make([]admin.AdminFood, 0)
	for rows.Next() {
		f, err := scanAdminFood(rows)
		if err != nil {
			return nil, 0, err
		}
		foods = append(foods, *f)
	}
	return foods, total, nil
}

func (r *adminFoodRepository) GetFoodByID(ctx context.Context, id uuid.UUID) (*admin.AdminFood, error) {
	query := `SELECT id, name, description, brand, image_url, barcode, serving_size, unit,
		calories, protein_g, carbs_g, fat_g, fiber_g, is_system, created_by_user_id, category, created_at, updated_at
		FROM foods WHERE id = $1`

	f, err := scanAdminFood(r.db.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("food")
		}
		return nil, err
	}
	return f, nil
}

func (r *adminFoodRepository) CreateSystemFood(ctx context.Context, input admin.CreateSystemFoodInput, adminID uuid.UUID) (*admin.AdminFood, error) {
	id := uuid.New()
	now := time.Now()
	_, err := r.db.Exec(ctx, `
		INSERT INTO foods (id, name, description, brand, image_url, barcode, serving_size, unit,
			calories, protein_g, carbs_g, fat_g, fiber_g, is_system, created_by_user_id, category, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`,
		id, input.Name, input.Description, input.Brand, input.ImageUrl, input.Barcode,
		input.ServingSize, input.Unit, input.Calories, input.ProteinG, input.CarbsG, input.FatG,
		input.FiberG, true, adminID, input.Category, now, now,
	)
	if err != nil {
		return nil, errors.DatabaseError("create system food", err)
	}
	return r.GetFoodByID(ctx, id)
}

func (r *adminFoodRepository) UpdateFood(ctx context.Context, id uuid.UUID, input admin.AdminUpdateFoodInput) (*admin.AdminFood, error) {
	setClauses := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIdx := 1

	addField := func(col string, val interface{}) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, argIdx))
		args = append(args, val)
		argIdx++
	}

	if input.Name != nil {
		addField("name", *input.Name)
	}
	if input.Description != nil {
		addField("description", *input.Description)
	}
	if input.Brand != nil {
		addField("brand", *input.Brand)
	}
	if input.ImageUrl != nil {
		addField("image_url", *input.ImageUrl)
	}
	if input.Barcode != nil {
		addField("barcode", *input.Barcode)
	}
	if input.ServingSize != nil {
		addField("serving_size", *input.ServingSize)
	}
	if input.Unit != nil {
		addField("unit", *input.Unit)
	}
	if input.Calories != nil {
		addField("calories", *input.Calories)
	}
	if input.ProteinG != nil {
		addField("protein_g", *input.ProteinG)
	}
	if input.CarbsG != nil {
		addField("carbs_g", *input.CarbsG)
	}
	if input.FatG != nil {
		addField("fat_g", *input.FatG)
	}
	if input.FiberG != nil {
		addField("fiber_g", *input.FiberG)
	}
	if input.Category != nil {
		addField("category", *input.Category)
	}
	if input.IsSystem != nil {
		addField("is_system", *input.IsSystem)
	}

	query := "UPDATE foods SET "
	for i, c := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += c
	}
	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, errors.DatabaseError("update food", err)
	}
	if result.RowsAffected() == 0 {
		return nil, errors.NotFound("food")
	}
	return r.GetFoodByID(ctx, id)
}

func (r *adminFoodRepository) DeleteFood(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM foods WHERE id = $1`, id)
	if err != nil {
		return errors.DatabaseError("delete food", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound("food")
	}
	return nil
}

func scanAdminFood(row interface {
	Scan(dest ...interface{}) error
}) (*admin.AdminFood, error) {
	var f admin.AdminFood
	err := row.Scan(
		&f.ID, &f.Name, &f.Description, &f.Brand, &f.ImageUrl, &f.Barcode,
		&f.ServingSize, &f.Unit, &f.Calories, &f.ProteinG, &f.CarbsG, &f.FatG,
		&f.FiberG, &f.IsSystem, &f.CreatedByUserID, &f.Category,
		&f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, errors.DatabaseError("scan food", err)
	}
	return &f, nil
}
