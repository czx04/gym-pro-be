package postgres

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/workout"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/pkg/errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ExerciseRepository implementation
type exerciseRepository struct {
	db *database.DB
}

func NewExerciseRepository(db *database.DB) workout.ExerciseRepository {
	return &exerciseRepository{db: db}
}

func (r *exerciseRepository) WithTx(tx *database.DB) workout.ExerciseRepository {
	return &exerciseRepository{db: tx}
}

// TODO: Implement all ExerciseRepository methods
func (r *exerciseRepository) Create(ctx context.Context, exercise *workout.Exercise) error {
	query := `
		INSERT INTO exercises (
			id, name, description, category, muscle_groups, equipment_needed, difficulty_level, calories_per_minute, video_url, thumbnail_url, is_active, created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`
	_, err := r.db.Exec(ctx, query,
		exercise.ID, exercise.Name, exercise.Description, exercise.Category, exercise.MuscleGroups, exercise.EquipmentNeeded, exercise.DifficultyLevel, exercise.CaloriesPerMinute, exercise.VideoURL, exercise.ThumbnailURL, exercise.IsActive, exercise.CreatedBy, exercise.CreatedAt, exercise.UpdatedAt,
	)
	if err != nil {
		return errors.DatabaseError("create exercise", err)
	}
	return nil
}

func (r *exerciseRepository) GetByID(ctx context.Context, id uuid.UUID) (*workout.Exercise, error) {
	query := `
		SELECT * FROM exercises
		WHERE id = $1
	`

	var exercise workout.Exercise
	err := r.db.QueryRow(ctx, query, id).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
		&exercise.Category,
		&exercise.MuscleGroups,
		&exercise.EquipmentNeeded,
		&exercise.DifficultyLevel,
		&exercise.CaloriesPerMinute,
		&exercise.VideoURL,
		&exercise.ThumbnailURL,
		&exercise.IsActive,
		&exercise.CreatedBy,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("exercise")
		}
		return nil, errors.DatabaseError("get exercise by id", err)
	}
	return &exercise, nil
}

func (r *exerciseRepository) List(ctx context.Context, page, pageSize int) ([]workout.Exercise, int64, error) {
	query := `
		SELECT * FROM exercises
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, pageSize, (page-1)*pageSize)
	if err != nil {
		logger.Error("error listing exercises", "err", err)
		return nil, 0, errors.DatabaseError("list exercises", err)
	}
	defer rows.Close()

	var totalCount int64
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM exercises").Scan(&totalCount)
	if err != nil {
		logger.Error("error counting exercises", "err", err)
		return nil, 0, errors.DatabaseError("count exercises", err)
	}
	exercises, err := r.rowsToModels(rows)
	if err != nil {
		logger.Error("error rows to models", "err", err)
		return nil, 0, errors.DatabaseError("list exercises", err)
	}
	return exercises, totalCount, nil
}

func (r *exerciseRepository) Search(ctx context.Context, filter workout.SearchExercisesFilter) ([]workout.Exercise, int64, error) {

	query := `
		SELECT * FROM exercises
		WHERE 1=1
	`
	countQuery := `
		SELECT COUNT(*) FROM exercises WHERE 1=1
	`

	args := make([]interface{}, 0)
	if filter.Category != nil && *filter.Category != "" {
		query += ` AND category = $` + strconv.Itoa(len(args)+1)
		countQuery += ` AND category = $` + strconv.Itoa(len(args)+1)
		args = append(args, filter.Category)
	}
	if filter.MuscleGroup != nil && *filter.MuscleGroup != "" {
		query += ` AND muscle_groups @> $` + strconv.Itoa(len(args)+1)
		countQuery += ` AND muscle_groups @> $` + strconv.Itoa(len(args)+1)
		args = append(args, []string{*filter.MuscleGroup})
	}
	if filter.Equipment != nil && *filter.Equipment != "" {
		query += ` AND equipment_needed @> $` + strconv.Itoa(len(args)+1)
		countQuery += ` AND equipment_needed @> $` + strconv.Itoa(len(args)+1)
		args = append(args, []string{*filter.Equipment})
	}
	if filter.DifficultyLevel != nil && *filter.DifficultyLevel != "" {
		query += ` AND difficulty_level = $` + strconv.Itoa(len(args)+1)
		countQuery += ` AND difficulty_level = $` + strconv.Itoa(len(args)+1)
		args = append(args, filter.DifficultyLevel)
	}
	if filter.Query != nil && *filter.Query != "" {
		query += ` AND name ILIKE '%' || $` + strconv.Itoa(len(args)+1) + ` || '%'`
		countQuery += ` AND name ILIKE '%' || $` + strconv.Itoa(len(args)+1) + ` || '%'`
		args = append(args, filter.Query)
	}
	query += `
		ORDER BY created_at DESC
		LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)

	var totalCount int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		logger.Error("error counting exercises", "err", err)
		return nil, 0, errors.DatabaseError("count exercises", err)
	}

	args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		logger.Error("error querying exercises", "err", err)
		return nil, 0, errors.DatabaseError("search exercises", err)
	}
	defer rows.Close()
	exercises, err := r.rowsToModels(rows)
	if err != nil {
		logger.Error("error rows to models", "err", err)
		return nil, 0, errors.DatabaseError("search exercises", err)
	}
	return exercises, totalCount, nil
}

func (r *exerciseRepository) Update(ctx context.Context, exercise *workout.Exercise) error {
	// TODO: Update exercise
	return nil
}

func (r *exerciseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// TODO: Soft delete (set is_active = false)
	return nil
}

func (r *exerciseRepository) rowToModel(row pgx.Row) (*workout.Exercise, error) {
	var exercise workout.Exercise
	err := row.Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
		&exercise.Category,
		&exercise.MuscleGroups,
		&exercise.EquipmentNeeded,
		&exercise.DifficultyLevel,
		&exercise.CaloriesPerMinute,
		&exercise.VideoURL,
		&exercise.ThumbnailURL,
		&exercise.IsActive,
		&exercise.CreatedBy,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)
	if err != nil {
		logger.Error("error row to model", "err", err)
		return nil, errors.DatabaseError("row to model", err)
	}
	return &exercise, nil
}

func (r *exerciseRepository) rowsToModels(rows pgx.Rows) ([]workout.Exercise, error) {
	exercises := make([]workout.Exercise, 0)
	for rows.Next() {
		exercise, err := r.rowToModel(rows)
		if err != nil {
			logger.Error("error rows to models", "err", err)
			return nil, errors.DatabaseError("rows to models", err)
		}
		exercises = append(exercises, *exercise)
	}
	return exercises, nil
}
