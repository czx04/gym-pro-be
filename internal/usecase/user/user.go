package user

import (
	"context"
	"mime/multipart"
	"strings"
	"time"

	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/infrastructure/auth"
	"gym-pro-2026-ptit/internal/infrastructure/email"
	"gym-pro-2026-ptit/internal/infrastructure/otp"
	mealuc "gym-pro-2026-ptit/internal/usecase/meal"
	"gym-pro-2026-ptit/pkg/cloudinary"
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
		ProteinTargetG     *int    `json:"protein_target_g,omitempty" validate:"omitempty,gt=0,lte=500"`
		CarbsTargetG       *int    `json:"carbs_target_g,omitempty" validate:"omitempty,gt=0,lte=1000"`
		FatTargetG         *int    `json:"fat_target_g,omitempty" validate:"omitempty,gt=0,lte=300"`
		EffectiveDate      *string `json:"effective_date,omitempty"`
	}
	UploadAvatarImageInput struct {
		File *multipart.FileHeader `form:"file" validate:"required"`
	}
	UploadAvatarImageOutput struct {
		AvatarURL string `json:"avatar_url"`
	}
	RequestChangeEmailOTPInput struct {
		NewEmail string `json:"new_email" validate:"required,email"`
	}
	VerifyChangeEmailOTPInput struct {
		NewEmail string `json:"new_email" validate:"required,email"`
		OTP      string `json:"otp" validate:"required,len=6"`
	}
)

// UserUseCases groups all user/auth use cases with a single dependency set.
type UserUseCases struct {
	userRepo     user.Repository
	mealDailyUC  *mealuc.MealDailyUseCases
	otpService   otp.Service
	emailService email.Service
	passwordMgr  *auth.PasswordManager
	jwtMgr       *auth.JWTManager
	validator    *validator.Validator
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

func (uc *UserUseCases) RequestChangeEmailOTP(ctx context.Context, userID uuid.UUID, input RequestChangeEmailOTPInput) error {
	if err := uc.validator.Validate(input); err != nil {
		return errors.Validation(err.Error())
	}

	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if u != nil && strings.EqualFold(strings.TrimSpace(u.Email), strings.TrimSpace(input.NewEmail)) {
		return errors.BadRequest("new_email must be different from current email")
	}

	exists, err := uc.userRepo.Exists(ctx, input.NewEmail)
	if err != nil {
		return err
	}
	if exists {
		return errors.Conflict("email already registered")
	}

	otpCode, err := uc.otpService.Generate(ctx, input.NewEmail)
	if err != nil {
		return err
	}
	if err := uc.emailService.SendOTP(input.NewEmail, otpCode); err != nil {
		return errors.InternalServer("failed to send OTP email", err)
	}
	return nil
}

func (uc *UserUseCases) VerifyChangeEmailOTP(ctx context.Context, userID uuid.UUID, input VerifyChangeEmailOTPInput) error {
	if err := uc.validator.Validate(input); err != nil {
		return errors.Validation(err.Error())
	}

	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if u != nil && strings.EqualFold(strings.TrimSpace(u.Email), strings.TrimSpace(input.NewEmail)) {
		return errors.BadRequest("new_email must be different from current email")
	}

	if err := uc.otpService.Verify(ctx, input.NewEmail, input.OTP); err != nil {
		return err
	}

	// Re-check existence to avoid race (and still rely on unique constraint in DB).
	exists, err := uc.userRepo.Exists(ctx, input.NewEmail)
	if err != nil {
		return err
	}
	if exists {
		return errors.Conflict("email already registered")
	}

	if err := uc.userRepo.UpdateEmail(ctx, userID, input.NewEmail); err != nil {
		return err
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
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if input.ProteinTargetG != nil || input.CarbsTargetG != nil || input.FatTargetG != nil || input.DailyCalorieTarget != nil {
		var p, c, f *int
		if input.ProteinTargetG != nil {
			p = input.ProteinTargetG
		} else {
			p = u.ProteinTargetG
		}
		if input.CarbsTargetG != nil {
			c = input.CarbsTargetG
		} else {
			c = u.CarbsTargetG
		}
		if input.FatTargetG != nil {
			f = input.FatTargetG
		} else {
			f = u.FatTargetG
		}

		if p == nil || c == nil || f == nil {
			return nil, errors.BadRequest("missing macro targets to compute calories")
		}
		if *p <= 0 || *c <= 0 || *f <= 0 {
			return nil, errors.BadRequest("macro targets must be greater than 0")
		}
		cal := (*p)*4 + (*c)*4 + (*f)*9
		input.DailyCalorieTarget = &cal
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

	if input.ProteinTargetG != nil {
		u.ProteinTargetG = input.ProteinTargetG
	}
	if input.CarbsTargetG != nil {
		u.CarbsTargetG = input.CarbsTargetG
	}
	if input.FatTargetG != nil {
		u.FatTargetG = input.FatTargetG
	}

	if u.ProteinTargetG == nil || u.CarbsTargetG == nil || u.FatTargetG == nil {
		return nil, errors.BadRequest("missing macro targets to compute calories")
	}
	if *u.ProteinTargetG <= 0 || *u.CarbsTargetG <= 0 || *u.FatTargetG <= 0 {
		return nil, errors.BadRequest("macro targets must be greater than 0")
	}

	cal := (*u.ProteinTargetG)*4 + (*u.CarbsTargetG)*4 + (*u.FatTargetG)*9
	u.DailyCalorieTarget = &cal

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

func (uc *UserUseCases) UploadAvatarImage(ctx context.Context, userID uuid.UUID, input UploadAvatarImageInput) (*UploadAvatarImageOutput, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	updated, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.InternalServer("failed to get user", err)
	}
	oldAvatarURL := updated.AvatarURL

	f, err := input.File.Open()
	if err != nil {
		return nil, errors.BadRequest("invalid file")
	}
	defer func() { _ = f.Close() }()

	imageURL, err := cloudinary.UploadAvatarImage(ctx, f, userID)
	if err != nil {
		return nil, errors.InternalServer("failed to upload avatar image", err)
	}
	updated.AvatarURL = &imageURL
	updated.UpdatedAt = time.Now()
	if err := uc.userRepo.Update(ctx, updated); err != nil {
		return nil, errors.InternalServer("failed to update user", err)
	}

	// Best-effort cleanup old avatar to avoid orphaned uploads.
	if oldAvatarURL != nil && *oldAvatarURL != "" && *oldAvatarURL != imageURL {
		_ = cloudinary.DeleteImage(ctx, *oldAvatarURL)
	}
	return &UploadAvatarImageOutput{AvatarURL: *updated.AvatarURL}, nil
}

func (uc *UserUseCases) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	// Fetch user first to get avatar URL for best-effort cleanup.
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	var avatarURL *string
	if u != nil {
		avatarURL = u.AvatarURL
	}

	if err := uc.userRepo.Delete(ctx, userID); err != nil {
		// Preserve NotFound (and other domain errors) when applicable.
		return err
	}

	// Best-effort cleanup avatar to avoid orphaned uploads.
	if avatarURL != nil && *avatarURL != "" {
		_ = cloudinary.DeleteImage(ctx, *avatarURL)
	}
	return nil
}
