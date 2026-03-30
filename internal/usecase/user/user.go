package user

import (
	"context"
	"time"

	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/helper"
	"gym-pro-2026-ptit/internal/infrastructure/auth"
	"gym-pro-2026-ptit/internal/infrastructure/email"
	"gym-pro-2026-ptit/internal/infrastructure/otp"
	mealuc "gym-pro-2026-ptit/internal/usecase/meal"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"

	"github.com/google/uuid"
)

// Input/Output types (shared across use cases)
type (
	RegisterRequestOTPInput struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	ResetPasswordRequestOTPInput struct {
		Email string `json:"email" validate:"required,email"`
	}
	VerifyOTPInput struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		OTP      string `json:"otp" validate:"required,len=6"`
	}
	VerifyOTPForgotPassword struct {
		Email string `json:"email" validate:"required,email"`
		OTP   string `json:"otp" validate:"required,len=6"`
	}
	ResetPasswordInput struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	TokenPair struct {
		AccessToken  string     `json:"access_token"`
		RefreshToken string     `json:"refresh_token"`
		User         *user.User `json:"user,omitempty"`
	}
	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	UpdateUserNutritionTargetInput struct {
		DailyCalorieTarget *int    `json:"daily_calorie_target,omitempty" validate:"omitempty,gte=500,lte=10000"`
		ProteinTargetG     *int    `json:"protein_target_g,omitempty" validate:"omitempty,gte=0,lte=500"`
		CarbsTargetG       *int    `json:"carbs_target_g,omitempty" validate:"omitempty,gte=0,lte=1000"`
		FatTargetG         *int    `json:"fat_target_g,omitempty" validate:"omitempty,gte=0,lte=300"`
		EffectiveDate      *string `json:"effective_date,omitempty"`
	}

	UpsertDailyStepsInput struct {
		Date   string `json:"date" validate:"required"` // YYYY-MM-DD
		Steps  int    `json:"steps" validate:"gte=0,lte=500000"`
		Source string `json:"source,omitempty"`
	}
)

// UserUseCases groups all user/auth use cases with a single dependency set.
type UserUseCases struct {
	userRepo      user.Repository
	mealDailyUC   *mealuc.MealDailyUseCases
	otpService    otp.Service
	emailService  email.Service
	passwordMgr   *auth.PasswordManager
	jwtMgr        *auth.JWTManager
	validator     *validator.Validator
}

// NewUserUseCases creates the user use cases container.
func NewUserUseCases(
	userRepo user.Repository,
	mealDailyUC *mealuc.MealDailyUseCases,
	otpService otp.Service,
	emailService email.Service,
	passwordMgr *auth.PasswordManager,
	jwtMgr *auth.JWTManager,
	validator *validator.Validator,
) *UserUseCases {
	return &UserUseCases{
		userRepo:     userRepo,
		mealDailyUC:  mealDailyUC,
		otpService:   otpService,
		emailService: emailService,
		passwordMgr:  passwordMgr,
		jwtMgr:       jwtMgr,
		validator:    validator,
	}
}

func (uc *UserUseCases) RegisterRequestOTP(ctx context.Context, input RegisterRequestOTPInput) error {
	if err := uc.validator.Validate(input); err != nil {
		return errors.Validation(err.Error())
	}
	exists, err := uc.userRepo.Exists(ctx, input.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.Conflict("email already registered")
	}
	otpCode, err := uc.otpService.Generate(ctx, input.Email)
	if err != nil {
		return err
	}
	if err := uc.emailService.SendOTP(input.Email, otpCode); err != nil {
		return errors.InternalServer("failed to send OTP email", err)
	}
	return nil
}

func (uc *UserUseCases) VerifyOTP(ctx context.Context, input VerifyOTPInput) (*TokenPair, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}
	if err := uc.otpService.Verify(ctx, input.Email, input.OTP); err != nil {
		return nil, err
	}
	exists, err := uc.userRepo.Exists(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.Conflict("user already exists")
	}
	passwordHash, err := uc.passwordMgr.HashPassword(input.Password)
	if err != nil {
		return nil, errors.InternalServer("failed to hash password", err)
	}
	newUser := &user.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: passwordHash,
		Name:         "User" + time.Now().Format("20060102150405"),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}
	accessToken, refreshToken, err := uc.jwtMgr.GenerateTokenPair(newUser.ID, newUser.Email)
	if err != nil {
		return nil, errors.InternalServer("failed to generate tokens", err)
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         newUser,
	}, nil
}

func (uc *UserUseCases) Login(ctx context.Context, input user.LoginInput) (*TokenPair, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}
	u, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.InvalidCredentials()
	}
	if !uc.passwordMgr.VerifyPassword(u.PasswordHash, input.Password) {
		return nil, errors.InvalidCredentials()
	}
	accessToken, refreshToken, err := uc.jwtMgr.GenerateTokenPair(u.ID, u.Email)
	if err != nil {
		return nil, errors.InternalServer("failed to generate tokens", err)
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         u,
	}, nil
}

func (uc *UserUseCases) RefreshToken(ctx context.Context, input RefreshTokenRequest) (*TokenPair, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}
	accessToken, err := uc.jwtMgr.RefreshAccessToken(input.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: input.RefreshToken,
	}, nil
}

func (uc *UserUseCases) GetProfile(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}

func (uc *UserUseCases) UpdateProfile(ctx context.Context, userID uuid.UUID, input user.UpdateProfileInput) (*user.User, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := uc.userRepo.UpdateProfile(ctx, userID, input); err != nil {
		return nil, err
	}
	updated, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	updated.PasswordHash = ""
	updated.UpdatedAt = time.Now()
	return updated, nil
}

func (uc *UserUseCases) ResetPasswordRequestOTP(ctx context.Context, input ResetPasswordRequestOTPInput) error {
	if err := uc.validator.Validate(input); err != nil {
		return errors.Validation(err.Error())
	}
	exists, err := uc.userRepo.Exists(ctx, input.Email)
	if !exists {
		return errors.NotFound("Email not found")
	}
	otpCode, err := uc.otpService.Generate(ctx, input.Email)
	if err != nil {
		return err
	}
	if err := uc.emailService.SendResetPasswordOTP(input.Email, otpCode); err != nil {
		return errors.InternalServer("failed to send OTP email", err)
	}
	return nil
}

func (uc *UserUseCases) VerifyOTPForgotPassword(ctx context.Context, input VerifyOTPForgotPassword) error {
	if err := uc.validator.Validate(input); err != nil {
		return errors.Validation(err.Error())
	}
	if err := uc.otpService.Verify(ctx, input.Email, input.OTP); err != nil {
		return err
	}
	return nil
}

func (uc *UserUseCases) ResetPassword(ctx context.Context, input ResetPasswordInput) (*user.User, error) {
	u, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.NotFound("Email not found")
	}
	passwordHash, err := uc.passwordMgr.HashPassword(input.Password)
	if err != nil {
		return nil, errors.InternalServer("Failed to hash password", err)
	}
	if err := uc.userRepo.UpdatePassword(ctx, u.ID, passwordHash); err != nil {
		return nil, errors.InternalServer("Failed to update password", err)
	}
	return u, nil
}

func (uc *UserUseCases) GetUserNutritionTarget(ctx context.Context, userID uuid.UUID) (*user.UserNutritionTarget, error) {
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.InternalServer("Failed to get user", err)
	}
	return &user.UserNutritionTarget{
		DailyCalorieTarget: u.DailyCalorieTarget,
		ProteinTargetG:     u.ProteinTargetG,
		CarbsTargetG:       u.CarbsTargetG,
		FatTargetG:         u.FatTargetG,
	}, nil
}

const maxWeightHistoryQueryRange = 732 * 24 * time.Hour // 2 years
const maxDailyStepsQueryRange = 732 * 24 * time.Hour   // 2 years

// ListMyWeightHistory returns chart points: latest measurement per calendar bucket in the given IANA timezone.
func (uc *UserUseCases) ListMyWeightHistory(ctx context.Context, userID uuid.UUID, from, to time.Time, tz string, granularity user.WeightHistoryGranularity) ([]user.WeightHistoryPoint, error) {
	if to.Before(from) {
		return nil, errors.BadRequest("to must be on or after from")
	}
	if to.Sub(from) > maxWeightHistoryQueryRange {
		return nil, errors.BadRequest("date range too large (max 2 years)")
	}
	switch granularity {
	case user.WeightHistoryGranularityDay, user.WeightHistoryGranularityWeek, user.WeightHistoryGranularityMonth:
	default:
		return nil, errors.BadRequest("granularity must be day, week, or month")
	}
	if tz == "" {
		tz = "UTC"
	}
	if _, err := time.LoadLocation(tz); err != nil {
		return nil, errors.BadRequest("invalid timezone")
	}
	return uc.userRepo.ListWeightHistoryByGranularity(ctx, userID, from, to, tz, granularity)
}

func (uc *UserUseCases) UpdateUserNutritionTarget(ctx context.Context, userID uuid.UUID, input UpdateUserNutritionTargetInput) (*user.UserNutritionTarget, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.InternalServer("Failed to get user", err)
	}

	helper.SetIfNotNil(&u.DailyCalorieTarget, &input.DailyCalorieTarget)
	helper.SetIfNotNil(&u.ProteinTargetG, &input.ProteinTargetG)
	helper.SetIfNotNil(&u.CarbsTargetG, &input.CarbsTargetG)
	helper.SetIfNotNil(&u.FatTargetG, &input.FatTargetG)
	if err := uc.userRepo.Update(ctx, u); err != nil {
		return nil, errors.InternalServer("Failed to update user", err)
	}

	effectiveDay := time.Now().UTC().Truncate(24 * time.Hour)
	if input.EffectiveDate != nil && *input.EffectiveDate != "" {
		d, err := time.Parse("2006-01-02", *input.EffectiveDate)
		if err != nil {
			return nil, errors.BadRequest("invalid effective_date, expected YYYY-MM-DD")
		}
		effectiveDay = d.UTC().Truncate(24 * time.Hour)
	}
	if err := uc.mealDailyUC.UpsertTargetsFromUserForDate(ctx, userID, effectiveDay); err != nil {
		return nil, errors.InternalServer("Failed to sync meal daily targets", err)
	}

	return &user.UserNutritionTarget{
		DailyCalorieTarget: u.DailyCalorieTarget,
		ProteinTargetG:     u.ProteinTargetG,
		CarbsTargetG:       u.CarbsTargetG,
		FatTargetG:         u.FatTargetG,
	}, nil
}

func (uc *UserUseCases) UpsertMyDailySteps(ctx context.Context, userID uuid.UUID, input UpsertDailyStepsInput) error {
	if err := uc.validator.Validate(input); err != nil {
		return errors.Validation(err.Error())
	}

	d, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return errors.BadRequest("invalid date, expected YYYY-MM-DD")
	}
	day := d.UTC().Truncate(24 * time.Hour)

	source := input.Source
	if source == "" {
		source = user.DailyStepsSourceAppleHealth
	}
	if source != user.DailyStepsSourceAppleHealth {
		return errors.BadRequest("invalid source")
	}

	return uc.userRepo.UpsertDailySteps(ctx, userID, day, source, input.Steps)
}

func (uc *UserUseCases) ListMyDailySteps(ctx context.Context, userID uuid.UUID, from, to time.Time, source string) ([]user.DailyStepsPoint, error) {
	if to.Before(from) {
		return nil, errors.BadRequest("to must be on or after from")
	}
	if to.Sub(from) > maxDailyStepsQueryRange {
		return nil, errors.BadRequest("date range too large (max 2 years)")
	}

	if source == "" {
		source = user.DailyStepsSourceAppleHealth
	}
	if source != user.DailyStepsSourceAppleHealth {
		return nil, errors.BadRequest("invalid source")
	}

	return uc.userRepo.ListDailySteps(ctx, userID, from.UTC().Truncate(24*time.Hour), to.UTC().Truncate(24*time.Hour), source)
}
