package postgres

import (
	"context"
	"fmt"
	"gym-pro-2026-ptit/internal/domain/workout"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/pkg/errors"
	"strings"
	"time"

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

func (r *workoutSessionRepository) Create(ctx context.Context, session *workout.WorkoutSession) error {
	q := `INSERT INTO workout_sessions (
		id, workout_schedule_id, user_id, workout_plan_id, scheduled_date, status, started_at, completed_at, duration_mins, total_calories_burned, notes, mood, difficulty_rating, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	var startedAt, completedAt interface{}
	if session.StartedAt != nil {
		startedAt = *session.StartedAt
	}
	if session.CompletedAt != nil {
		completedAt = *session.CompletedAt
	}
	_, err := r.db.Exec(ctx, q,
		session.ID, session.WorkoutScheduleID, session.UserID, session.WorkoutPlanID, session.ScheduledDate, session.Status, startedAt, completedAt,
		session.DurationMins, session.TotalCaloriesBurned, session.Notes, session.Mood, session.DifficultyRating, session.CreatedAt, session.UpdatedAt,
	)
	if err != nil {
		return errors.DatabaseError("create workout session", err)
	}
	for i := range session.Exercises {
		ex := &session.Exercises[i]
		if err := r.insertSessionExercise(ctx, session.ID, ex); err != nil {
			return err
		}
	}
	return nil
}

func (r *workoutSessionRepository) insertSessionExercise(ctx context.Context, sessionID uuid.UUID, ex *workout.WorkoutSessionExercise) error {
	q := `INSERT INTO workout_session_exercises (id, workout_session_id, exercise_id, "order", target_sets, target_reps, duration_secs, notes, skipped)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(ctx, q, ex.ID, sessionID, ex.ExerciseID, ex.Order, ex.TargetSets, ex.TargetReps, ex.DurationSecs, ex.Notes, ex.Skipped)
	if err != nil {
		return errors.DatabaseError("insert session exercise", err)
	}
	sets := ex.TargetSets
	if sets == nil || *sets < 1 {
		return nil
	}
	for i := 0; i < *sets; i++ {
		setID := uuid.New()
		_, err = r.db.Exec(ctx, `INSERT INTO workout_session_sets (id, workout_session_exercise_id, set_index, completed) VALUES ($1, $2, $3, false)`,
			setID, ex.ID, i+1)
		if err != nil {
			return errors.DatabaseError("insert session set", err)
		}
	}
	return nil
}

func (r *workoutSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*workout.WorkoutSession, error) {
	q := `SELECT s.id, s.workout_schedule_id, s.user_id, s.workout_plan_id, s.scheduled_date::text, s.status, s.started_at, s.completed_at, s.duration_mins, s.total_calories_burned, s.notes, s.mood, s.difficulty_rating, s.created_at, s.updated_at, p.title
		FROM workout_sessions s LEFT JOIN workout_plans p ON p.id = s.workout_plan_id WHERE s.id = $1`
	var s workout.WorkoutSession
	var scheduledDate *string
	var title *string
	err := r.db.QueryRow(ctx, q, id).Scan(
		&s.ID, &s.WorkoutScheduleID, &s.UserID, &s.WorkoutPlanID, &scheduledDate, &s.Status, &s.StartedAt, &s.CompletedAt,
		&s.DurationMins, &s.TotalCaloriesBurned, &s.Notes, &s.Mood, &s.DifficultyRating, &s.CreatedAt, &s.UpdatedAt, &title,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("workout session")
		}
		return nil, errors.DatabaseError("get workout session", err)
	}
	if scheduledDate != nil {
		s.ScheduledDate = scheduledDate
	}
	if title != nil {
		s.Title = *title
	}
	s.Exercises, err = r.getSessionExercisesWithSets(ctx, id)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *workoutSessionRepository) getSessionExercisesWithSets(ctx context.Context, sessionID uuid.UUID) ([]workout.WorkoutSessionExercise, error) {
	q := `SELECT e.id, e.workout_session_id, e.exercise_id, e."order", e.target_sets, e.target_reps, e.duration_secs, e.notes, e.skipped,
		ex.id, ex.name, ex.description, ex.category, ex.muscle_groups, ex.equipment_needed, ex.difficulty_level, ex.calories_per_minute, ex.video_url, ex.thumbnail_url, ex.is_active
		FROM workout_session_exercises e JOIN exercises ex ON ex.id = e.exercise_id WHERE e.workout_session_id = $1 ORDER BY e."order"`
	rows, err := r.db.Query(ctx, q, sessionID)
	if err != nil {
		return nil, errors.DatabaseError("get session exercises", err)
	}
	defer rows.Close()
	var list []workout.WorkoutSessionExercise
	for rows.Next() {
		var ex workout.WorkoutSessionExercise
		ex.Exercise = &workout.Exercise{}
		err := rows.Scan(
			&ex.ID, &ex.WorkoutSessionID, &ex.ExerciseID, &ex.Order, &ex.TargetSets, &ex.TargetReps, &ex.DurationSecs, &ex.Notes, &ex.Skipped,
			&ex.Exercise.ID, &ex.Exercise.Name, &ex.Exercise.Description, &ex.Exercise.Category, &ex.Exercise.MuscleGroups, &ex.Exercise.EquipmentNeeded, &ex.Exercise.DifficultyLevel, &ex.Exercise.CaloriesPerMinute, &ex.Exercise.VideoURL, &ex.Exercise.ThumbnailURL, &ex.Exercise.IsActive,
		)
		if err != nil {
			return nil, errors.DatabaseError("scan session exercise", err)
		}
		ex.Sets, err = r.getSessionSets(ctx, ex.ID)
		if err != nil {
			return nil, err
		}
		list = append(list, ex)
	}
	return list, nil
}

func (r *workoutSessionRepository) getSessionSets(ctx context.Context, sessionExerciseID uuid.UUID) ([]workout.WorkoutSessionSet, error) {
	rows, err := r.db.Query(ctx, `SELECT id, workout_session_exercise_id, set_index, reps, weight_kg, completed, completed_at, created_at, updated_at FROM workout_session_sets WHERE workout_session_exercise_id = $1 ORDER BY set_index`, sessionExerciseID)
	if err != nil {
		return nil, errors.DatabaseError("get session sets", err)
	}
	defer rows.Close()
	var list []workout.WorkoutSessionSet
	for rows.Next() {
		var s workout.WorkoutSessionSet
		err := rows.Scan(&s.ID, &s.WorkoutSessionExerciseID, &s.SetIndex, &s.Reps, &s.WeightKg, &s.Completed, &s.CompletedAt, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, errors.DatabaseError("scan session set", err)
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *workoutSessionRepository) GetScheduledDates(ctx context.Context, userID uuid.UUID, month, year int) ([]string, error) {
	q := `SELECT DISTINCT scheduled_date::text FROM workout_sessions WHERE user_id = $1 AND scheduled_date IS NOT NULL AND EXTRACT(MONTH FROM scheduled_date) = $2 AND EXTRACT(YEAR FROM scheduled_date) = $3 ORDER BY 1`
	rows, err := r.db.Query(ctx, q, userID, month, year)
	if err != nil {
		return nil, errors.DatabaseError("get scheduled dates", err)
	}
	defer rows.Close()
	var dates []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, errors.DatabaseError("scan scheduled date", err)
		}
		dates = append(dates, d)
	}
	return dates, nil
}

func (r *workoutSessionRepository) GetByDate(ctx context.Context, userID uuid.UUID, date string) ([]workout.WorkoutSession, error) {
	q := `SELECT s.id, s.workout_schedule_id, s.user_id, s.workout_plan_id, s.scheduled_date::text, s.status, s.started_at, s.completed_at, s.duration_mins, s.total_calories_burned, s.notes, s.mood, s.difficulty_rating, s.created_at, s.updated_at, p.title
		FROM workout_sessions s LEFT JOIN workout_plans p ON p.id = s.workout_plan_id
		WHERE s.user_id = $1 AND (s.scheduled_date = $2::date OR (s.started_at IS NOT NULL AND s.started_at::date = $2::date)) ORDER BY s.scheduled_date NULLS LAST, s.started_at NULLS LAST`
	rows, err := r.db.Query(ctx, q, userID, date)
	if err != nil {
		return nil, errors.DatabaseError("get sessions by date", err)
	}
	defer rows.Close()
	var list []workout.WorkoutSession
	for rows.Next() {
		var s workout.WorkoutSession
		var scheduledDate *string
		var title *string
		err := rows.Scan(
			&s.ID, &s.WorkoutScheduleID, &s.UserID, &s.WorkoutPlanID, &scheduledDate, &s.Status, &s.StartedAt, &s.CompletedAt,
			&s.DurationMins, &s.TotalCaloriesBurned, &s.Notes, &s.Mood, &s.DifficultyRating, &s.CreatedAt, &s.UpdatedAt, &title,
		)
		if err != nil {
			return nil, errors.DatabaseError("scan session", err)
		}
		if scheduledDate != nil {
			s.ScheduledDate = scheduledDate
		}
		if title != nil {
			s.Title = *title
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *workoutSessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]workout.WorkoutSession, int64, error) {
	q := `SELECT id, workout_schedule_id, user_id, workout_plan_id, scheduled_date::text, status, started_at, completed_at, duration_mins, total_calories_burned, notes, mood, difficulty_rating, created_at, updated_at FROM workout_sessions WHERE user_id = $1 ORDER BY COALESCE(started_at, created_at) DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, q, userID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, errors.DatabaseError("get sessions by user", err)
	}
	defer rows.Close()
	var list []workout.WorkoutSession
	for rows.Next() {
		var s workout.WorkoutSession
		var scheduledDate *string
		err := rows.Scan(&s.ID, &s.WorkoutScheduleID, &s.UserID, &s.WorkoutPlanID, &scheduledDate, &s.Status, &s.StartedAt, &s.CompletedAt, &s.DurationMins, &s.TotalCaloriesBurned, &s.Notes, &s.Mood, &s.DifficultyRating, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, 0, errors.DatabaseError("scan session", err)
		}
		if scheduledDate != nil {
			s.ScheduledDate = scheduledDate
		}
		list = append(list, s)
	}
	var total int64
	_ = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM workout_sessions WHERE user_id = $1", userID).Scan(&total)
	return list, total, nil
}

func (r *workoutSessionRepository) Update(ctx context.Context, session *workout.WorkoutSession) error {
	q := `UPDATE workout_sessions SET status = $2, started_at = $3, completed_at = $4, duration_mins = $5, total_calories_burned = $6, notes = $7, mood = $8, difficulty_rating = $9, updated_at = $10 WHERE id = $1`
	_, err := r.db.Exec(ctx, q, session.ID, session.Status, session.StartedAt, session.CompletedAt, session.DurationMins, session.TotalCaloriesBurned, session.Notes, session.Mood, session.DifficultyRating, session.UpdatedAt)
	return errors.DatabaseError("update workout session", err)
}

func (r *workoutSessionRepository) Complete(ctx context.Context, id uuid.UUID, input workout.CompleteWorkoutSessionInput) error {
	var durationMins *int
	if input.DurationSecs != nil {
		m := *input.DurationSecs / 60
		durationMins = &m
	}
	q := `UPDATE workout_sessions SET status = 'completed', completed_at = COALESCE($2, NOW()), duration_mins = COALESCE($3, duration_mins), notes = COALESCE($4, notes), mood = $5, difficulty_rating = $6, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id, input.CompletedAt, durationMins, input.Notes, input.Mood, input.DifficultyRating)
	if err != nil {
		return errors.DatabaseError("complete workout session", err)
	}
	return nil
}

func (r *workoutSessionRepository) UpdateSet(ctx context.Context, setID uuid.UUID, input workout.UpdateSessionSetInput) error {
	// PATCH: reps/weight_kg/completed; weight_kg nullable (client có thể gửi null để xóa)
	q := `UPDATE workout_session_sets SET reps = COALESCE($2, reps), weight_kg = $3, completed = COALESCE($4, completed), completed_at = CASE WHEN $4 = true THEN COALESCE(completed_at, NOW()) WHEN $4 = false THEN NULL ELSE completed_at END, updated_at = NOW() WHERE id = $1`
	var reps interface{} = input.Reps
	var weightKg interface{} = input.WeightKg
	var completed interface{}
	if input.Completed != nil {
		completed = *input.Completed
	}
	_, err := r.db.Exec(ctx, q, setID, reps, weightKg, completed)
	if err != nil {
		return errors.DatabaseError("update session set", err)
	}
	return nil
}

func (r *workoutSessionRepository) AddExerciseLog(ctx context.Context, sessionID uuid.UUID, exercise *workout.WorkoutSessionExercise) error {
	return r.insertSessionExercise(ctx, sessionID, exercise)
}

func (r *workoutSessionRepository) GetExercises(ctx context.Context, sessionID uuid.UUID) ([]workout.WorkoutSessionExercise, error) {
	return r.getSessionExercisesWithSets(ctx, sessionID)
}

func (r *workoutSessionRepository) GetStats(ctx context.Context, userID uuid.UUID) (*workout.WorkoutStats, error) {
	// TODO: aggregate stats
	return nil, nil
}

func (r *workoutSessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM workout_sessions WHERE id = $1", id)
	if err != nil {
		return errors.DatabaseError("delete workout session", err)
	}
	return nil
}

func (r *workoutSessionRepository) GetExerciseStats(ctx context.Context, userID, exerciseID uuid.UUID) (*workout.ExerciseStats, error) {
	stats := &workout.ExerciseStats{}
	qMax := `SELECT MAX(s.weight_kg) FROM workout_session_sets s
		JOIN workout_session_exercises e ON e.id = s.workout_session_exercise_id
		JOIN workout_sessions sess ON sess.id = e.workout_session_id
		WHERE e.exercise_id = $1 AND sess.user_id = $2 AND s.weight_kg IS NOT NULL`
	if err := r.db.QueryRow(ctx, qMax, exerciseID, userID).Scan(&stats.MaxWeightKg); err != nil && err != pgx.ErrNoRows {
		return nil, errors.DatabaseError("get exercise max weight", err)
	}
	qLast := `SELECT s.weight_kg, COALESCE(s.completed_at, s.updated_at) AS logged_at
		FROM workout_session_sets s
		JOIN workout_session_exercises e ON e.id = s.workout_session_exercise_id
		JOIN workout_sessions sess ON sess.id = e.workout_session_id
		WHERE e.exercise_id = $1 AND sess.user_id = $2 AND s.weight_kg IS NOT NULL
		ORDER BY COALESCE(s.completed_at, s.updated_at) DESC, s.updated_at DESC
		LIMIT 1`
	var lastLoggedAt time.Time
	err := r.db.QueryRow(ctx, qLast, exerciseID, userID).Scan(&stats.LastWeightKg, &lastLoggedAt)
	if err != nil && err != pgx.ErrNoRows {
		return nil, errors.DatabaseError("get exercise last weight", err)
	}
	if err == nil {
		stats.LastLoggedAt = &lastLoggedAt
	}
	return stats, nil
}
