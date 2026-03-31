package workout

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/domain/workout"
	cacheinfra "gym-pro-2026-ptit/internal/infrastructure/cache"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/utils"
	"gym-pro-2026-ptit/pkg/validator"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type WorkoutUseCases struct {
	db              *database.DB
	cache           *cacheinfra.Cache
	userRepo        user.Repository
	workoutPlanRepo workout.WorkoutPlanRepository
	sessionRepo     workout.WorkoutSessionRepository
	validator       *validator.Validator
}

type ProfileWorkoutStats struct {
	TotalWorkouts    int64 `json:"total_workouts"`
	TotalWorkoutDays int64 `json:"total_workout_days"`
}

func (uc *WorkoutUseCases) GetProfileWorkoutStats(ctx context.Context, userID uuid.UUID) (*ProfileWorkoutStats, error) {
	totalWorkouts, totalDays, err := uc.sessionRepo.GetProfileWorkoutStats(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &ProfileWorkoutStats{
		TotalWorkouts:    totalWorkouts,
		TotalWorkoutDays: totalDays,
	}, nil
}

func NewWorkoutUseCases(
	db *database.DB,
	cache *cacheinfra.Cache,
	userRepo user.Repository,
	workoutPlanRepo workout.WorkoutPlanRepository,
	sessionRepo workout.WorkoutSessionRepository,
	validator *validator.Validator,
) *WorkoutUseCases {
	return &WorkoutUseCases{
		db:              db,
		cache:           cache,
		userRepo:        userRepo,
		workoutPlanRepo: workoutPlanRepo,
		sessionRepo:     sessionRepo,
		validator:       validator,
	}
}

const aiWeeklySummaryCacheTTL = 24 * time.Hour

type cachedAIWeeklySummary struct {
	Insights        []workout.WeeklyInsight `json:"insights"`
	Recommendations []string                `json:"recommendations"`
	AISummary       string                  `json:"ai_summary,omitempty"`
	AIModel         string                  `json:"ai_model,omitempty"`
}

func (uc *WorkoutUseCases) CreateWorkoutPlan(ctx context.Context, u *user.User, input workout.CreateWorkoutPlanInput) (*workout.WorkoutPlan, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}
	isTemplate := u.IsAdmin()

	tx, err := uc.db.Begin(ctx)
	if err != nil {
		return nil, errors.DatabaseError("begin transaction", err)
	}
	defer tx.Rollback(ctx)

	planRepo := uc.workoutPlanRepo.WithTx(tx)

	plan := &workout.WorkoutPlan{
		ID:              uuid.New(),
		UserID:          u.ID,
		Title:           input.Title,
		Description:     input.Description,
		DifficultyLevel: input.DifficultyLevel,
		IsTemplate:      isTemplate,
		IsPublic:        input.IsPublic,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := planRepo.Create(ctx, plan); err != nil {
		return nil, err
	}

	exercises := make([]*workout.WorkoutPlanExercise, len(input.Exercises))
	for i, exercise := range input.Exercises {
		exercises[i] = &workout.WorkoutPlanExercise{
			WorkoutPlanID: plan.ID,
			ExerciseID:    exercise.ExerciseID,
			Order:         exercise.Order,
			Sets:          exercise.Sets,
			Reps:          exercise.Reps,
			DurationSecs:  exercise.DurationSecs,
			RestSecs:      exercise.RestSecs,
			Notes:         exercise.Notes,
		}
	}

	if len(exercises) > 0 {
		if err := planRepo.AddExercise(ctx, plan.ID, exercises); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.DatabaseError("commit transaction", err)
	}
	return plan, nil
}

func (uc *WorkoutUseCases) ListWorkoutPlans(ctx context.Context, user user.User, page, pageSize int) ([]workout.WorkoutPlan, int64, error) {
	return uc.workoutPlanRepo.GetByUserID(ctx, user.ID, page, pageSize)
}

func (uc *WorkoutUseCases) GetWorkoutPlan(ctx context.Context, userID uuid.UUID, planID string) (*workout.WorkoutPlan, error) {
	uuidPlanID, err := uuid.Parse(planID)
	if err != nil {
		logger.Error("error parsing plan ID", "err", err, "planID", planID)
		return nil, errors.BadRequest("invalid plan ID")
	}
	plan, err := uc.workoutPlanRepo.GetByID(ctx, uuidPlanID)
	if err != nil {
		logger.Error("error getting workout plan by ID", "err", err, "planID", planID)
		return nil, errors.DatabaseError("get workout plan by ID", err)
	}
	if plan.UserID != userID {
		logger.Error("user is not allowed to get this workout plan", "userID", userID, "planID", planID)
		return nil, errors.Forbidden("you are not allowed to get this workout plan")
	}

	exercises, err := uc.workoutPlanRepo.GetExercises(ctx, plan.ID)
	if err != nil {
		logger.Error("error getting exercises by plan id", "err", err, "planID", planID)
		return nil, errors.DatabaseError("get exercises by plan id", err)
	}
	plan.Exercises = exercises
	return plan, nil
}
func (uc *WorkoutUseCases) DeleteWorkoutPlan(ctx context.Context, userID uuid.UUID, planID string) error {
	uuidPlanID, err := uuid.Parse(planID)
	if err != nil {
		logger.Error("error parsing plan ID", "err", err, "planID", planID)
		return errors.BadRequest("invalid plan ID")
	}
	plan, err := uc.workoutPlanRepo.GetByID(ctx, uuidPlanID)
	if err != nil {
		logger.Error("error getting workout plan by ID", "err", err, "planID", planID)
		return errors.DatabaseError("get workout plan by ID", err)
	}
	if plan.UserID != userID {
		logger.Error("user is not allowed to delete this workout plan", "userID", userID, "planID", planID)
		return errors.Forbidden("you are not allowed to delete this workout plan")
	}
	return uc.workoutPlanRepo.Delete(ctx, plan.ID)
}

func (uc *WorkoutUseCases) UpdateWorkoutPlan(ctx context.Context, userID uuid.UUID, input workout.UpdateWorkoutPlanInput) (*workout.WorkoutPlan, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	plan, err := uc.workoutPlanRepo.GetByID(ctx, input.ID)
	if err != nil {
		logger.Error("error getting workout plan by ID", "err", err, "planID", input.ID)
		return nil, errors.DatabaseError("get workout plan by ID", err)
	}
	if plan.UserID != userID {
		logger.Error("user is not allowed to update this workout plan", "userID", userID, "planID", input.ID)
		return nil, errors.Forbidden("you are not allowed to update this workout plan")
	}

	uc.buildWorkoutPlanFromUpdateInput(plan, input)

	db, err := uc.db.Begin(ctx)
	if err != nil {
		return nil, errors.DatabaseError("begin transaction", err)
	}
	defer db.Rollback(ctx)

	planRepo := uc.workoutPlanRepo.WithTx(db)

	if err := planRepo.Update(ctx, plan); err != nil {
		return nil, errors.DatabaseError("update workout plan", err)
	}

	if input.IsUpdateExercises && len(input.Exercises) > 0 {
		if err := planRepo.RemoveExercise(ctx, plan.ID); err != nil {
			return nil, errors.DatabaseError("remove exercise from workout plan", err)
		}
		exercises := make([]*workout.WorkoutPlanExercise, len(input.Exercises))
		for i, exercise := range input.Exercises {
			exercises[i] = &workout.WorkoutPlanExercise{
				WorkoutPlanID: plan.ID,
				ExerciseID:    exercise.ID,
				Order:         exercise.Order,
				Sets:          exercise.Sets,
				Reps:          exercise.Reps,
				DurationSecs:  exercise.DurationSecs,
				RestSecs:      exercise.RestSecs,
				Notes:         exercise.Notes,
			}
		}
		if err := planRepo.AddExercise(ctx, plan.ID, exercises); err != nil {
			return nil, errors.DatabaseError("update exercise in workout plan", err)
		}
	}

	if err := db.Commit(ctx); err != nil {
		return nil, errors.DatabaseError("commit transaction", err)
	}
	return plan, nil
}

func (uc *WorkoutUseCases) buildWorkoutPlanFromUpdateInput(currentPlan *workout.WorkoutPlan, input workout.UpdateWorkoutPlanInput) {
	if input.Title != nil {
		currentPlan.Title = *input.Title
	}
	if input.Description != nil {
		currentPlan.Description = input.Description
	}
	if input.DifficultyLevel != nil {
		currentPlan.DifficultyLevel = *input.DifficultyLevel
	}
	if input.IsTemplate != nil {
		currentPlan.IsTemplate = *input.IsTemplate
	}
	if input.IsPublic != nil {
		currentPlan.IsPublic = *input.IsPublic
	}
	currentPlan.UpdatedAt = time.Now()
}

// ——— Workout Session (Calendar / Tracking) ———

func (uc *WorkoutUseCases) GetScheduledDates(ctx context.Context, userID uuid.UUID, month, year int) ([]string, error) {
	return uc.sessionRepo.GetScheduledDates(ctx, userID, month, year)
}

func (uc *WorkoutUseCases) GetSessionsByDate(ctx context.Context, userID uuid.UUID, date string) ([]workout.WorkoutSession, error) {
	return uc.sessionRepo.GetByDate(ctx, userID, date)
}

func (uc *WorkoutUseCases) DeleteWorkoutSession(ctx context.Context, userID uuid.UUID, sessionID string) error {
	id, err := uuid.Parse(sessionID)
	if err != nil {
		return errors.BadRequest("invalid session ID")
	}

	s, err := uc.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if s.UserID != userID {
		return errors.Forbidden("not your session")
	}

	now := time.Now().In(time.Local)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	target := time.Date(s.CreatedAt.Year(), s.CreatedAt.Month(), s.CreatedAt.Day(), 0, 0, 0, 0, time.Local)
	if s.ScheduledDate != nil {
		if scheduledDate, parseErr := time.ParseInLocation("2006-01-02", *s.ScheduledDate, time.Local); parseErr == nil {
			target = time.Date(scheduledDate.Year(), scheduledDate.Month(), scheduledDate.Day(), 0, 0, 0, 0, time.Local)
		}
	} else if s.StartedAt != nil {
		target = time.Date(s.StartedAt.Year(), s.StartedAt.Month(), s.StartedAt.Day(), 0, 0, 0, 0, time.Local)
	}

	if target.Before(today) {
		return errors.Forbidden("can only delete sessions for today or future dates")
	}

	return uc.sessionRepo.Delete(ctx, id)
}

func (uc *WorkoutUseCases) GetSessionByID(ctx context.Context, userID uuid.UUID, sessionID string) (*workout.WorkoutSession, error) {
	id, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, errors.BadRequest("invalid session ID")
	}
	s, err := uc.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.UserID != userID {
		return nil, errors.Forbidden("not your session")
	}
	return s, nil
}

func (uc *WorkoutUseCases) CreateWorkoutSession(ctx context.Context, userID uuid.UUID, input workout.CreateWorkoutSessionInput) (*workout.WorkoutSession, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}
	plan, err := uc.workoutPlanRepo.GetByID(ctx, input.WorkoutPlanID)
	if err != nil {
		return nil, err
	}
	if plan.UserID != userID && !plan.IsPublic {
		return nil, errors.Forbidden("plan not found or not yours")
	}
	exercises, err := uc.workoutPlanRepo.GetExercises(ctx, plan.ID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	session := &workout.WorkoutSession{
		ID:            uuid.New(),
		UserID:        userID,
		WorkoutPlanID: plan.ID,
		ScheduledDate: &input.ScheduledDate,
		Status:        workout.SessionStatusScheduled,
		CreatedAt:     now,
		UpdatedAt:     now,
		Title:         plan.Title,
	}
	if input.StartNow {
		session.Status = workout.SessionStatusInProgress
		session.StartedAt = &now
	}
	session.Exercises = make([]workout.WorkoutSessionExercise, 0, len(exercises))
	for _, pe := range exercises {
		ex := workout.WorkoutSessionExercise{
			ID:               uuid.New(),
			WorkoutSessionID: session.ID,
			ExerciseID:       pe.ExerciseID,
			Order:            pe.Order,
			TargetSets:       pe.Sets,
			TargetReps:       pe.Reps,
			DurationSecs:     pe.DurationSecs,
			Notes:            pe.Notes,
			Skipped:          false,
			Exercise:         pe.Exercise,
		}
		session.Exercises = append(session.Exercises, ex)
	}
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}
	session.Title = plan.Title
	return session, nil
}

func (uc *WorkoutUseCases) UpdateWorkoutSession(ctx context.Context, userID uuid.UUID, sessionID string, input workout.UpdateWorkoutSessionInput) (*workout.WorkoutSession, error) {
	id, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, errors.BadRequest("invalid session ID")
	}
	s, err := uc.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.UserID != userID {
		return nil, errors.Forbidden("not your session")
	}
	if input.Status != nil {
		s.Status = *input.Status
	}
	if input.StartedAt != nil {
		s.StartedAt = input.StartedAt
	}
	s.UpdatedAt = time.Now()
	if err := uc.sessionRepo.Update(ctx, s); err != nil {
		return nil, err
	}
	return uc.sessionRepo.GetByID(ctx, id)
}

func (uc *WorkoutUseCases) UpdateSessionSet(ctx context.Context, userID uuid.UUID, sessionID, setID string, input workout.UpdateSessionSetInput) error {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return errors.BadRequest("invalid session ID")
	}
	setUUID, err := uuid.Parse(setID)
	if err != nil {
		return errors.BadRequest("invalid set ID")
	}
	s, err := uc.sessionRepo.GetByID(ctx, sid)
	if err != nil {
		return err
	}
	if s.UserID != userID {
		return errors.Forbidden("not your session")
	}
	return uc.sessionRepo.UpdateSet(ctx, setUUID, input)
}

func (uc *WorkoutUseCases) FinishWorkoutSession(ctx context.Context, userID uuid.UUID, sessionID string, input workout.CompleteWorkoutSessionInput) (*workout.WorkoutSession, error) {
	id, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, errors.BadRequest("invalid session ID")
	}
	s, err := uc.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.UserID != userID {
		return nil, errors.Forbidden("not your session")
	}
	if err := uc.sessionRepo.Complete(ctx, id, input); err != nil {
		return nil, err
	}
	return uc.sessionRepo.GetByID(ctx, id)
}

func (uc *WorkoutUseCases) GetWeeklySummary(ctx context.Context, userID uuid.UUID, input workout.GetWeeklySummaryRequest) (*workout.WeeklyWorkoutSummary, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, errors.BadRequest("invalid start_date format, expected YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, errors.BadRequest("invalid end_date format, expected YYYY-MM-DD")
	}
	if startDate.After(endDate) {
		return nil, errors.BadRequest("start date must be before end date")
	}

	currentStart := startDate
	currentEndExclusive := endDate.AddDate(0, 0, 1)
	days := int(currentEndExclusive.Sub(currentStart).Hours() / 24)
	if days <= 0 {
		return nil, errors.BadRequest("date range must contain at least 1 day")
	}

	previousEndExclusive := currentStart
	previousStart := previousEndExclusive.AddDate(0, 0, -days)

	current, err := uc.sessionRepo.GetWeeklyAggregate(ctx, userID, currentStart, currentEndExclusive)
	if err != nil {
		return nil, err
	}
	previous, err := uc.sessionRepo.GetWeeklyAggregate(ctx, userID, previousStart, previousEndExclusive)
	if err != nil {
		return nil, err
	}

	roundWeeklyMetrics(current)
	roundWeeklyMetrics(previous)

	summary := &workout.WeeklyWorkoutSummary{
		StartDate:         currentStart.Format("2006-01-02"),
		EndDate:           endDate.Format("2006-01-02"),
		PreviousStartDate: previousStart.Format("2006-01-02"),
		PreviousEndDate:   previousEndExclusive.AddDate(0, 0, -1).Format("2006-01-02"),
		Current:           *current,
		Previous:          *previous,
		StrengthTrend:     buildTrendDelta(current.TotalVolumeKg, previous.TotalVolumeKg, 0.05),
		RestTrend:         buildTrendDelta(current.AvgRestSecs, previous.AvgRestSecs, 0.05),
		MoodTrend:         buildTrendDelta(current.AvgMoodScore, previous.AvgMoodScore, 0.1),
	}

	if current.RestSamples == 0 || previous.RestSamples == 0 {
		summary.RestTrend.Trend = "insufficient_data"
	}

	currentWeight, err := uc.userRepo.GetLatestWeightInRange(ctx, userID, currentStart, currentEndExclusive)
	if err != nil {
		return nil, err
	}
	previousWeight, err := uc.userRepo.GetLatestWeightInRange(ctx, userID, previousStart, previousEndExclusive)
	if err != nil {
		return nil, err
	}

	if currentWeight != nil && previousWeight != nil {
		summary.BodyWeightTrend = buildTrendDelta(currentWeight.WeightKg, previousWeight.WeightKg, 0.3)
	} else {
		summary.BodyWeightTrend = workout.TrendDelta{Trend: "insufficient_data"}
	}

	summary.Insights, summary.Recommendations = buildWeeklyInsightsAndRecommendations(summary)
	summary.RecommendationSource = "rule_based"
	aiCacheKey := uc.buildAIWeeklySummaryCacheKey(userID, summary.StartDate, summary.EndDate)

	if cachedAI, ok := uc.getCachedAIWeeklySummary(ctx, aiCacheKey); ok {
		summary.Insights = cachedAI.Insights
		summary.Recommendations = cachedAI.Recommendations
		summary.RecommendationSource = "gemini"
		summary.AISummary = cachedAI.AISummary
		summary.AIModel = cachedAI.AIModel
		return summary, nil
	}

	var fitnessGoal *string
	if u, err := uc.userRepo.GetByID(ctx, userID); err == nil {
		fitnessGoal = u.FitnessGoal
	} else {
		logger.Warn("failed to load user profile for ai recommendation", "user_id", userID, "err", err)
	}

	aiInsights, aiRecommendations, aiSummary, aiModel, err := tryGeminiRecommendationsWithRetry(ctx, summary, fitnessGoal)
	if err != nil {
		logger.Warn("gemini recommendation fallback to rule-based", "user_id", userID, "err", err)
		return summary, nil
	}

	if len(aiInsights) > 0 {
		summary.Insights = aiInsights
	}
	if len(aiRecommendations) > 0 {
		summary.Recommendations = aiRecommendations
	}
	summary.RecommendationSource = "gemini"
	summary.AISummary = aiSummary
	summary.AIModel = aiModel

	uc.setCachedAIWeeklySummary(ctx, aiCacheKey, cachedAIWeeklySummary{
		Insights:        summary.Insights,
		Recommendations: summary.Recommendations,
		AISummary:       summary.AISummary,
		AIModel:         summary.AIModel,
	})

	return summary, nil
}

func (uc *WorkoutUseCases) buildAIWeeklySummaryCacheKey(userID uuid.UUID, startDate, endDate string) string {
	return fmt.Sprintf("weekly_ai_summary:%s:%s:%s", userID.String(), startDate, endDate)
}

func (uc *WorkoutUseCases) getCachedAIWeeklySummary(ctx context.Context, cacheKey string) (*cachedAIWeeklySummary, bool) {
	if uc.cache == nil {
		return nil, false
	}

	raw, err := uc.cache.Get(ctx, cacheKey)
	if err != nil {
		if err != redis.Nil {
			logger.Warn("failed to read ai summary cache", "cache_key", cacheKey, "err", err)
		}
		return nil, false
	}

	var cached cachedAIWeeklySummary
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		logger.Warn("failed to unmarshal ai summary cache", "cache_key", cacheKey, "err", err)
		return nil, false
	}

	if len(cached.Recommendations) == 0 {
		return nil, false
	}

	return &cached, true
}

func (uc *WorkoutUseCases) setCachedAIWeeklySummary(ctx context.Context, cacheKey string, payload cachedAIWeeklySummary) {
	if uc.cache == nil {
		return
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		logger.Warn("failed to marshal ai summary cache payload", "cache_key", cacheKey, "err", err)
		return
	}

	if err := uc.cache.Set(ctx, cacheKey, raw, aiWeeklySummaryCacheTTL); err != nil {
		logger.Warn("failed to set ai summary cache", "cache_key", cacheKey, "err", err)
	}
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type geminiWeeklyOutput struct {
	Summary         string                  `json:"summary"`
	Insights        []workout.WeeklyInsight `json:"insights"`
	Recommendations []string                `json:"recommendations"`
}

func tryGeminiRecommendations(ctx context.Context, summary *workout.WeeklyWorkoutSummary, fitnessGoal *string) ([]workout.WeeklyInsight, []string, string, string, error) {
	apiKey := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	if apiKey == "" {
		return nil, nil, "", "", fmt.Errorf("missing GEMINI_API_KEY")
	}

	model := strings.TrimSpace(os.Getenv("GEMINI_MODEL"))
	if model == "" {
		model = "gemini-1.5-flash"
	}

	timeoutSeconds := 8
	if rawTimeout := strings.TrimSpace(os.Getenv("GEMINI_TIMEOUT_SECONDS")); rawTimeout != "" {
		if t, err := strconv.Atoi(rawTimeout); err == nil && t > 0 {
			timeoutSeconds = t
		}
	}

	inputPayload := map[string]interface{}{
		"period": map[string]string{
			"start_date":          summary.StartDate,
			"end_date":            summary.EndDate,
			"previous_start_date": summary.PreviousStartDate,
			"previous_end_date":   summary.PreviousEndDate,
		},
		"current":                  summary.Current,
		"previous":                 summary.Previous,
		"strength_trend":           summary.StrengthTrend,
		"rest_trend":               summary.RestTrend,
		"mood_trend":               summary.MoodTrend,
		"body_weight_trend":        summary.BodyWeightTrend,
		"fallback_insights":        summary.Insights,
		"fallback_recommendations": summary.Recommendations,
	}
	if fitnessGoal != nil {
		inputPayload["fitness_goal"] = *fitnessGoal
	}

	metricsJSON, err := json.Marshal(inputPayload)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("marshal gemini metrics payload: %w", err)
	}

	prompt := "You are a professional but safe fitness coach. Analyze weekly workout metrics and return concise Vietnamese recommendations. " +
		"Do not provide medical diagnosis. Focus on practical next-week actions. " +
		"Return ONLY valid JSON with this schema: {\"summary\": string, \"insights\": [{\"code\": string, \"severity\": \"positive\"|\"warning\"|\"neutral\", \"message\": string, \"evidence\": string, \"metric_key\": string}], \"recommendations\": [string]}. " +
		"Input metrics JSON: " + string(metricsJSON)

	reqBody := geminiRequest{
		Contents: []geminiContent{{
			Parts: []geminiPart{{Text: prompt}},
		}},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("marshal gemini request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("create gemini request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("call gemini: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("read gemini response: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, nil, "", "", fmt.Errorf("gemini returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var gResp geminiResponse
	if err := json.Unmarshal(respBytes, &gResp); err != nil {
		return nil, nil, "", "", fmt.Errorf("unmarshal gemini response wrapper: %w", err)
	}
	if len(gResp.Candidates) == 0 || len(gResp.Candidates[0].Content.Parts) == 0 {
		return nil, nil, "", "", fmt.Errorf("empty gemini candidates")
	}

	rawText := strings.TrimSpace(gResp.Candidates[0].Content.Parts[0].Text)
	jsonText := extractJSONBody(rawText)

	var output geminiWeeklyOutput
	if err := json.Unmarshal([]byte(jsonText), &output); err != nil {
		return nil, nil, "", "", fmt.Errorf("unmarshal gemini content json: %w", err)
	}

	if len(output.Recommendations) == 0 {
		return nil, nil, "", "", fmt.Errorf("gemini output has empty recommendations")
	}

	cleanedInsights := make([]workout.WeeklyInsight, 0, len(output.Insights))
	for _, item := range output.Insights {
		if item.Code == "" || item.Message == "" {
			continue
		}
		if item.Severity == "" {
			item.Severity = "neutral"
		}
		cleanedInsights = append(cleanedInsights, item)
	}

	cleanedRecommendations := make([]string, 0, len(output.Recommendations))
	for _, rec := range output.Recommendations {
		rec = strings.TrimSpace(rec)
		if rec != "" {
			cleanedRecommendations = append(cleanedRecommendations, rec)
		}
	}
	if len(cleanedRecommendations) == 0 {
		return nil, nil, "", "", fmt.Errorf("gemini output has no valid recommendations")
	}

	return cleanedInsights, cleanedRecommendations, strings.TrimSpace(output.Summary), model, nil
}

func tryGeminiRecommendationsWithRetry(ctx context.Context, summary *workout.WeeklyWorkoutSummary, fitnessGoal *string) ([]workout.WeeklyInsight, []string, string, string, error) {
	attempts := 3
	if rawAttempts := strings.TrimSpace(os.Getenv("GEMINI_RETRY_ATTEMPTS")); rawAttempts != "" {
		if v, err := strconv.Atoi(rawAttempts); err == nil && v > 0 {
			attempts = v
		}
	}

	backoff := 300 * time.Millisecond
	if rawBackoff := strings.TrimSpace(os.Getenv("GEMINI_RETRY_BACKOFF_MS")); rawBackoff != "" {
		if v, err := strconv.Atoi(rawBackoff); err == nil && v > 0 {
			backoff = time.Duration(v) * time.Millisecond
		}
	}

	var lastErr error
	for i := 1; i <= attempts; i++ {
		insights, recommendations, summaryText, model, err := tryGeminiRecommendations(ctx, summary, fitnessGoal)
		if err == nil {
			return insights, recommendations, summaryText, model, nil
		}

		lastErr = err
		if i == attempts {
			break
		}

		logger.Warn("gemini request failed, retrying", "attempt", i, "max_attempts", attempts, "err", err)

		wait := backoff * time.Duration(i)
		select {
		case <-ctx.Done():
			return nil, nil, "", "", ctx.Err()
		case <-time.After(wait):
		}
	}

	return nil, nil, "", "", fmt.Errorf("gemini failed after %d attempts: %w", attempts, lastErr)
}

func extractJSONBody(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "```") {
		raw = strings.TrimPrefix(raw, "```json")
		raw = strings.TrimPrefix(raw, "```")
		raw = strings.TrimSuffix(raw, "```")
		raw = strings.TrimSpace(raw)
	}

	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start >= 0 && end > start {
		return raw[start : end+1]
	}
	return raw
}

func roundWeeklyMetrics(m *workout.WeeklyWorkoutMetrics) {
	if m == nil {
		return
	}
	m.TotalVolumeKg = utils.RoundToTwo(m.TotalVolumeKg)
	m.AvgWeightKg = utils.RoundToTwo(m.AvgWeightKg)
	m.AvgRestSecs = utils.RoundToTwo(m.AvgRestSecs)
	m.AvgMoodScore = utils.RoundToTwo(m.AvgMoodScore)
	m.AvgDifficulty = utils.RoundToTwo(m.AvgDifficulty)
	m.CompletionRate = utils.RoundToTwo(m.CompletionRate * 100)
}

func buildTrendDelta(current, previous, stableThreshold float64) workout.TrendDelta {
	delta := current - previous
	trend := "stable"
	threshold := stableThreshold
	if threshold > 0 && math.Abs(previous) > 0 {
		threshold = math.Abs(previous) * stableThreshold
	}
	if math.Abs(delta) <= threshold {
		trend = "stable"
	} else if delta > 0 {
		trend = "up"
	} else {
		trend = "down"
	}

	return workout.TrendDelta{
		Current:  utils.RoundToTwo(current),
		Previous: utils.RoundToTwo(previous),
		Delta:    utils.RoundToTwo(delta),
		Trend:    trend,
	}
}

func buildWeeklyInsightsAndRecommendations(summary *workout.WeeklyWorkoutSummary) ([]workout.WeeklyInsight, []string) {
	insights := make([]workout.WeeklyInsight, 0, 8)
	recommendations := make([]string, 0, 6)
	seenRecommendations := map[string]struct{}{}

	for _, rule := range defaultWeeklyRecommendationRules() {
		if !matchesRule(summary, rule) {
			continue
		}

		insight := workout.WeeklyInsight{
			Code:      rule.Code,
			Severity:  rule.Severity,
			MetricKey: rule.MetricKey,
			Message:   rule.Message,
			Evidence:  buildEvidenceFromRule(summary, rule),
		}
		insights = append(insights, insight)

		if rule.Recommendation != "" {
			if _, exists := seenRecommendations[rule.Recommendation]; !exists {
				recommendations = append(recommendations, rule.Recommendation)
				seenRecommendations[rule.Recommendation] = struct{}{}
			}
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Duy tri lich tap hien tai va tiep tuc theo doi trend tai, mood va thoi gian nghi moi tuan.")
	}

	return insights, recommendations
}

type weeklyRecommendationRule struct {
	Code           string
	Severity       string
	MetricKey      string
	Message        string
	EvidenceLabel  string
	EvidenceValue  string
	EvidenceUnit   string
	Recommendation string
	Conditions     []ruleCondition
	ConditionMode  string // all|any
}

type ruleCondition struct {
	Metric   string
	Operator string // eq|neq|gt|gte|lt|lte|abs_gte|abs_lte
	NumValue float64
	StrValue string
}

func defaultWeeklyRecommendationRules() []weeklyRecommendationRule {
	return []weeklyRecommendationRule{
		{
			Code:          "strength_progress",
			Severity:      "positive",
			MetricKey:     "total_volume_kg",
			Message:       "Tong khoi luong ta cua ban dang tang tuan qua tuan.",
			EvidenceLabel: "Delta volume",
			EvidenceValue: "strength_delta",
			EvidenceUnit:  "kg",
			Conditions: []ruleCondition{{
				Metric:   "strength_trend",
				Operator: "eq",
				StrValue: "up",
			}},
		},
		{
			Code:           "strength_drop",
			Severity:       "warning",
			MetricKey:      "total_volume_kg",
			Message:        "Tong khoi luong ta giam so voi tuan truoc.",
			EvidenceLabel:  "Delta volume",
			EvidenceValue:  "strength_delta",
			EvidenceUnit:   "kg",
			Recommendation: "Giam tai nhe 1 tuan hoac tang thoi gian phuc hoi neu ban cam thay met moi.",
			Conditions: []ruleCondition{{
				Metric:   "strength_trend",
				Operator: "eq",
				StrValue: "down",
			}},
		},
		{
			Code:           "mood_drop",
			Severity:       "warning",
			MetricKey:      "avg_mood_score",
			Message:        "Mood tap luyen co xu huong giam.",
			EvidenceLabel:  "Mood delta",
			EvidenceValue:  "mood_delta",
			Recommendation: "Can nhac giam khoi luong 5-10% va uu tien ngu de cai thien recovery.",
			Conditions: []ruleCondition{{
				Metric:   "mood_trend",
				Operator: "eq",
				StrValue: "down",
			}},
		},
		{
			Code:          "rest_increase",
			Severity:      "neutral",
			MetricKey:     "avg_rest_secs",
			Message:       "Thoi gian nghi trung binh giua cac set dang tang.",
			EvidenceLabel: "Rest delta",
			EvidenceValue: "rest_delta",
			EvidenceUnit:  "sec",
			Conditions: []ruleCondition{
				{Metric: "rest_trend", Operator: "eq", StrValue: "up"},
				{Metric: "rest_trend", Operator: "neq", StrValue: "insufficient_data"},
			},
			ConditionMode: "all",
		},
		{
			Code:          "body_weight_change",
			Severity:      "neutral",
			MetricKey:     "body_weight",
			Message:       "Can nang co thay doi so voi tuan truoc.",
			EvidenceLabel: "Weight delta",
			EvidenceValue: "body_weight_delta",
			EvidenceUnit:  "kg",
			Conditions: []ruleCondition{
				{Metric: "body_weight_trend", Operator: "eq", StrValue: "up"},
				{Metric: "body_weight_trend", Operator: "eq", StrValue: "down"},
			},
			ConditionMode: "any",
		},
		{
			Code:           "body_weight_insufficient_data",
			Severity:       "neutral",
			MetricKey:      "body_weight",
			Message:        "Chua du du lieu can nang de so sanh voi tuan truoc.",
			Recommendation: "Cap nhat can nang deu dan trong profile de theo doi xu huong chinh xac hon.",
			Conditions: []ruleCondition{{
				Metric:   "body_weight_trend",
				Operator: "eq",
				StrValue: "insufficient_data",
			}},
		},
		{
			Code:           "no_completed_workouts",
			Severity:       "warning",
			MetricKey:      "total_workouts",
			Message:        "Tuan nay ban chua hoan thanh buoi tap nao.",
			Recommendation: "Hay dat muc tieu toi thieu 2-3 buoi tap cho tuan toi va theo doi muc do hoan thanh.",
			Conditions: []ruleCondition{{
				Metric:   "current_total_workouts",
				Operator: "eq",
				NumValue: 0,
			}},
		},
	}
}

func matchesRule(summary *workout.WeeklyWorkoutSummary, rule weeklyRecommendationRule) bool {
	if len(rule.Conditions) == 0 {
		return true
	}

	modeAny := rule.ConditionMode == "any"
	if modeAny {
		for _, cond := range rule.Conditions {
			if evaluateCondition(summary, cond) {
				return true
			}
		}
		return false
	}

	for _, cond := range rule.Conditions {
		if !evaluateCondition(summary, cond) {
			return false
		}
	}
	return true
}

func evaluateCondition(summary *workout.WeeklyWorkoutSummary, cond ruleCondition) bool {
	if num, ok := getNumericMetric(summary, cond.Metric); ok {
		switch cond.Operator {
		case "eq":
			return math.Abs(num-cond.NumValue) < 1e-9
		case "neq":
			return math.Abs(num-cond.NumValue) >= 1e-9
		case "gt":
			return num > cond.NumValue
		case "gte":
			return num >= cond.NumValue
		case "lt":
			return num < cond.NumValue
		case "lte":
			return num <= cond.NumValue
		case "abs_gte":
			return math.Abs(num) >= cond.NumValue
		case "abs_lte":
			return math.Abs(num) <= cond.NumValue
		default:
			return false
		}
	}

	if str, ok := getStringMetric(summary, cond.Metric); ok {
		switch cond.Operator {
		case "eq":
			return str == cond.StrValue
		case "neq":
			return str != cond.StrValue
		default:
			return false
		}
	}

	return false
}

func buildEvidenceFromRule(summary *workout.WeeklyWorkoutSummary, rule weeklyRecommendationRule) string {
	if rule.EvidenceLabel == "" || rule.EvidenceValue == "" {
		return ""
	}
	v, ok := getNumericMetric(summary, rule.EvidenceValue)
	if !ok {
		return ""
	}
	if rule.EvidenceUnit != "" {
		return rule.EvidenceLabel + ": " + formatFloat(v) + " " + rule.EvidenceUnit
	}
	return rule.EvidenceLabel + ": " + formatFloat(v)
}

func getNumericMetric(summary *workout.WeeklyWorkoutSummary, key string) (float64, bool) {
	switch key {
	case "strength_delta":
		return summary.StrengthTrend.Delta, true
	case "rest_delta":
		return summary.RestTrend.Delta, true
	case "mood_delta":
		return summary.MoodTrend.Delta, true
	case "body_weight_delta":
		return summary.BodyWeightTrend.Delta, true
	case "current_total_workouts":
		return float64(summary.Current.TotalWorkouts), true
	default:
		return 0, false
	}
}

func getStringMetric(summary *workout.WeeklyWorkoutSummary, key string) (string, bool) {
	switch key {
	case "strength_trend":
		return summary.StrengthTrend.Trend, true
	case "rest_trend":
		return summary.RestTrend.Trend, true
	case "mood_trend":
		return summary.MoodTrend.Trend, true
	case "body_weight_trend":
		return summary.BodyWeightTrend.Trend, true
	default:
		return "", false
	}
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(utils.RoundToTwo(v), 'f', 2, 64)
}
