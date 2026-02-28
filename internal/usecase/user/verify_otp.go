package user

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/infrastructure/auth"
	"gym-pro-2026-ptit/internal/infrastructure/otp"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"
	"time"

	"github.com/google/uuid"
)

type VerifyOTPInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	OTP      string `json:"otp" validate:"required,len=6"`
}

type TokenPair struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	User         *user.User `json:"user"`
}

type VerifyOTPUseCase struct {
	userRepo    user.Repository
	otpService  otp.Service
	passwordMgr *auth.PasswordManager
	jwtMgr      *auth.JWTManager
	validator   *validator.Validator
}

func NewVerifyOTPUseCase(
	userRepo user.Repository,
	otpService otp.Service,
	passwordMgr *auth.PasswordManager,
	jwtMgr *auth.JWTManager,
	validator *validator.Validator,
) *VerifyOTPUseCase {
	return &VerifyOTPUseCase{
		userRepo:    userRepo,
		otpService:  otpService,
		passwordMgr: passwordMgr,
		jwtMgr:      jwtMgr,
		validator:   validator,
	}
}

func (uc *VerifyOTPUseCase) Execute(ctx context.Context, input VerifyOTPInput) (*TokenPair, error) {
	// Validate input
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	// Verify OTP
	if err := uc.otpService.Verify(ctx, input.Email, input.OTP); err != nil {
		return nil, err
	}

	// Check if user was already created (edge case - duplicate verification)
	exists, err := uc.userRepo.Exists(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.Conflict("user already exists")
	}

	// Hash password
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

	// Save user to database
	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := uc.jwtMgr.GenerateTokenPair(newUser.ID, newUser.Email)
	if err != nil {
		return nil, errors.InternalServer("failed to generate tokens", err)
	}

	// Return tokens and user
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         newUser,
	}, nil
}
