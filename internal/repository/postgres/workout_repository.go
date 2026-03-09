package postgres

import (
	"context"
	"fmt"
	"gym-pro-2026-ptit/internal/domain/workout"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/pkg/errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// WorkoutPlanRepository implementation
type workoutPlanRepository struct {
	db *database.DB
}

func NewWorkoutPlanRepository(db *database.DB) workout.WorkoutPlanRepository {
	return &workoutPlanRepository{db: db}
}

func (r *workoutPlanRepository) WithTx(tx *database.DB) workout.WorkoutPlanRepository {
	return &workoutPlanRepository{db: tx}
}

// TODO: Implement all WorkoutPlanRepository methods
func (r *workoutPlanRepository) Create(ctx context.Context, plan *workout.WorkoutPlan) error {
	query := `
		INSERT INTO workout_plans (
			id, user_id, title, description, difficulty_level, estimated_duration_mins, estimated_calories, is_template, is_public, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`
	_, err := r.db.Exec(ctx, query,
		plan.ID, plan.UserID, plan.Title, plan.Description, plan.DifficultyLevel, plan.EstimatedDurationMins, plan.EstimatedCalories, plan.IsTemplate, plan.IsPublic, plan.CreatedAt, plan.UpdatedAt,
	)
	if err != nil {
		logger.Error("error creating workout plan", "err", err, "plan", plan)
		return errors.DatabaseError("create workout plan", err)
	}
	return nil
}

func (r *workoutPlanRepository) GetByID(ctx context.Context, id uuid.UUID) (*workout.WorkoutPlan, error) {
	query := `
		SELECT * FROM workout_plans
		WHERE id = $1
	`
	var plan workout.WorkoutPlan
	err := r.db.QueryRow(ctx, query, id).Scan(
		&plan.ID,
		&plan.UserID,
		&plan.Title,
		&plan.Description,
		&plan.DifficultyLevel,
		&plan.EstimatedDurationMins,
		&plan.EstimatedCalories,
		&plan.IsTemplate,
		&plan.IsPublic,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("workout plan")
		}
		return nil, errors.DatabaseError("get workout plan by id", err)
	}
	return &plan, nil
}

func (r *workoutPlanRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]workout.WorkoutPlan, int64, error) {
	query := `
		SELECT * FROM workout_plans
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, pageSize, (page-1)*pageSize)
	if err != nil {
		logger.Error("error getting workout plans by user id", "err", err)
		return nil, 0, errors.DatabaseError("get workout plans by user id", err)
	}
	defer rows.Close()

	var totalCount int64
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM workout_plans WHERE user_id = $1", userID).Scan(&totalCount)
	if err != nil {
		logger.Error("error getting total workout plans by user id", "err", err)
		return nil, 0, errors.DatabaseError("get total workout plans by user id", err)
	}

	plans, err := r.rowsToModels(rows)
	if err != nil {
		logger.Error("error rows to models", "err", err)
		return nil, 0, errors.DatabaseError("get workout plans by user id", err)
	}
	return plans, totalCount, nil
}

func (r *workoutPlanRepository) rowToModel(row pgx.Row) (*workout.WorkoutPlan, error) {
	var plan workout.WorkoutPlan
	err := row.Scan(
		&plan.ID,
		&plan.UserID,
		&plan.Title,
		&plan.Description,
		&plan.DifficultyLevel,
		&plan.EstimatedDurationMins,
		&plan.EstimatedCalories,
		&plan.IsTemplate,
		&plan.IsPublic,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err != nil {
		return nil, errors.DatabaseError("row to model", err)
	}
	return &plan, nil
}

func (r *workoutPlanRepository) rowsToModels(rows pgx.Rows) ([]workout.WorkoutPlan, error) {
	plans := make([]workout.WorkoutPlan, 0)
	for rows.Next() {
		plan, err := r.rowToModel(rows)
		if err != nil {
			return nil, errors.DatabaseError("rows to models", err)
		}
		plans = append(plans, *plan)
	}
	return plans, nil
}

func (r *workoutPlanRepository) Update(ctx context.Context, plan *workout.WorkoutPlan) error {
	query := `
		UPDATE workout_plans SET
		title = $2,
		description = $3,
		difficulty_level = $4,
		estimated_duration_mins = $5,
		estimated_calories = $6,
		is_template = $7,
		is_public = $8,
		updated_at = $9
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, plan.ID, plan.Title, plan.Description, plan.DifficultyLevel, plan.EstimatedDurationMins, plan.EstimatedCalories, plan.IsTemplate, plan.IsPublic, plan.UpdatedAt)
	if err != nil {
		logger.Error("error updating workout plan", "err", err)
		return errors.DatabaseError("update workout plan", err)
	}
	return nil
}

func (r *workoutPlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM workout_plans WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		logger.Error("error deleting workout plan", "err", err)
		return errors.DatabaseError("delete workout plan", err)
	}
	return nil
}

func (r *workoutPlanRepository) AddExercise(ctx context.Context, planID uuid.UUID, exercises []*workout.WorkoutPlanExercise) error {
	valueStrings := make([]string, 0, len(exercises))
	valueArgs := make([]interface{}, 0, len(exercises)*8)
	for i, exercise := range exercises {
		n := i * 8
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8))
		valueArgs = append(valueArgs, exercise.WorkoutPlanID, exercise.ExerciseID, exercise.Order,
			exercise.Sets, exercise.Reps, exercise.DurationSecs, exercise.RestSecs, exercise.Notes)
	}

	query := fmt.Sprintf("INSERT INTO workout_plan_exercises (workout_plan_id, exercise_id, \"order\", sets, reps, duration_secs, rest_secs, notes) VALUES %s",
		strings.Join(valueStrings, ","))
	_, err := r.db.Exec(ctx, query, valueArgs...)
	if err != nil {
		logger.Error("error adding exercise to workout plan", "err", err)
		return errors.DatabaseError("add exercise to workout plan", err)
	}
	return nil
}

func (r *workoutPlanRepository) UpdateExercise(ctx context.Context, planExerciseID uuid.UUID, input workout.UpdateExerciseInWorkoutInput) error {
	// TODO: Update exercise configuration in plan
	return nil
}

func (r *workoutPlanRepository) RemoveExercise(ctx context.Context, planID uuid.UUID) error {
	query := `
		DELETE FROM workout_plan_exercises WHERE workout_plan_id = $1
	`

	_, err := r.db.Exec(ctx, query, planID)
	if err != nil {
		logger.Error("error removing exercise from workout plan", "err", err)
		return errors.DatabaseError("remove exercise from workout plan", err)
	}
	return nil
}

func (r *workoutPlanRepository) GetExercises(ctx context.Context, planID uuid.UUID) ([]workout.WorkoutPlanExercise, error) {
	query := `
		SELECT workout_plan_exercises.*, exercises.* FROM workout_plan_exercises
		JOIN exercises ON workout_plan_exercises.exercise_id = exercises.id
		WHERE workout_plan_id = $1
	`
	rows, err := r.db.Query(ctx, query, planID)
	if err != nil {
		logger.Error("error getting exercises by plan id", "err", err)
		return nil, errors.DatabaseError("get exercises by plan id", err)
	}
	defer rows.Close()
	exercises := make([]workout.WorkoutPlanExercise, 0)
	for rows.Next() {
		exercise := workout.WorkoutPlanExercise{
			Exercise: &workout.Exercise{},
		}
		err := rows.Scan(
			&exercise.ID,
			&exercise.WorkoutPlanID,
			&exercise.ExerciseID,
			&exercise.Order,
			&exercise.Sets,
			&exercise.Reps,
			&exercise.DurationSecs,
			&exercise.RestSecs,
			&exercise.Notes,
			&exercise.Exercise.ID,
			&exercise.Exercise.Name,
			&exercise.Exercise.Description,
			&exercise.Exercise.Category,
			&exercise.Exercise.MuscleGroups,
			&exercise.Exercise.EquipmentNeeded,
			&exercise.Exercise.DifficultyLevel,
			&exercise.Exercise.CaloriesPerMinute,
			&exercise.Exercise.VideoURL,
			&exercise.Exercise.ThumbnailURL,
			&exercise.Exercise.IsActive,
			&exercise.Exercise.CreatedBy,
			&exercise.Exercise.CreatedAt,
			&exercise.Exercise.UpdatedAt,
		)
		if err != nil {
			logger.Error("error scanning exercise", "err", err)
			return nil, errors.DatabaseError("get exercises by plan id", err)
		}
		exercises = append(exercises, exercise)
	}
	return exercises, nil
}

// WorkoutScheduleRepository implementation
type workoutScheduleRepository struct {
	db *database.DB
}

func NewWorkoutScheduleRepository(db *database.DB) workout.WorkoutScheduleRepository {
	return &workoutScheduleRepository{db: db}
}

func (r *workoutScheduleRepository) WithTx(tx *database.DB) workout.WorkoutScheduleRepository {
	return &workoutScheduleRepository{db: tx}
}

// TODO: Implement all WorkoutScheduleRepository methods
func (r *workoutScheduleRepository) Create(ctx context.Context, schedule *workout.WorkoutSchedule) error {
	// TODO: Insert into workout_schedules
	return nil
}

func (r *workoutScheduleRepository) GetByID(ctx context.Context, id uuid.UUID) (*workout.WorkoutSchedule, error) {
	// TODO: Query schedule with workout plan details
	return nil, nil
}

func (r *workoutScheduleRepository) GetByUserID(ctx context.Context, userID uuid.UUID, filter workout.GetScheduleFilter) ([]workout.WorkoutSchedule, error) {
	// TODO: Query schedules with filters (date range, completed status)
	return nil, nil
}

func (r *workoutScheduleRepository) GetByDateRange(ctx context.Context, userID uuid.UUID, filter workout.GetScheduleFilter) ([]workout.WorkoutSchedule, error) {
	// TODO: Query schedules in date range
	return nil, nil
}

func (r *workoutScheduleRepository) Update(ctx context.Context, schedule *workout.WorkoutSchedule) error {
	// TODO: Update schedule
	return nil
}

func (r *workoutScheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// TODO: Delete schedule
	return nil
}

func (r *workoutScheduleRepository) MarkCompleted(ctx context.Context, id uuid.UUID) error {
	// TODO: Set is_completed = true, completed_at = NOW()
	return nil
}

func (r *workoutScheduleRepository) BulkCreate(ctx context.Context, schedules []workout.WorkoutSchedule) error {
	// TODO: Batch insert multiple schedules
	return nil
}

// WorkoutSessionRepository implementation
type workoutSessionRepository struct {
	db *database.DB
}

func NewWorkoutSessionRepository(db *database.DB) workout.WorkoutSessionRepository {
	return &workoutSessionRepository{db: db}
}

func (r *workoutSessionRepository) WithTx(tx *database.DB) workout.WorkoutSessionRepository {
	return &workoutSessionRepository{db: tx}
}

// TODO: Implement all WorkoutSessionRepository methods
func (r *workoutSessionRepository) Create(ctx context.Context, session *workout.WorkoutSession) error {
	// TODO: Insert into workout_sessions
	return nil
}

func (r *workoutSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*workout.WorkoutSession, error) {
	// TODO: Query session with exercises and workout plan details
	return nil, nil
}

func (r *workoutSessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]workout.WorkoutSession, int64, error) {
	// TODO: Query user's sessions with pagination
	return nil, 0, nil
}

func (r *workoutSessionRepository) Update(ctx context.Context, session *workout.WorkoutSession) error {
	// TODO: Update session
	return nil
}

func (r *workoutSessionRepository) Complete(ctx context.Context, id uuid.UUID, input workout.CompleteWorkoutSessionInput) error {
	// TODO: Set completed_at, calculate duration_mins, update notes/mood/difficulty
	return nil
}

func (r *workoutSessionRepository) AddExerciseLog(ctx context.Context, sessionID uuid.UUID, exercise *workout.WorkoutSessionExercise) error {
	// TODO: Insert or update workout_session_exercises
	// Handle actual_sets_completed as JSONB
	return nil
}

func (r *workoutSessionRepository) GetExercises(ctx context.Context, sessionID uuid.UUID) ([]workout.WorkoutSessionExercise, error) {
	// TODO: Query exercises in session with exercise details
	return nil, nil
}

func (r *workoutSessionRepository) GetStats(ctx context.Context, userID uuid.UUID) (*workout.WorkoutStats, error) {
	// TODO: Calculate statistics:
	// - Total workouts
	// - Total duration
	// - Total calories
	// - Average duration
	// - Completion rate (scheduled vs completed)
	return nil, nil
}

func (r *workoutSessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// TODO: Delete session
	return nil
}
