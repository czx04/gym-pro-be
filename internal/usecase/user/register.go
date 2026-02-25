package user

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/infrastructure/auth"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"
	"time"

	"github.com/google/uuid"
)

// RegisterUseCase handles user registration
type RegisterUseCase struct {
	userRepo    user.Repository
	passwordMgr *auth.PasswordManager
	jwtMgr      *auth.JWTManager
	validator   *validator.Validator
}

// NewRegisterUseCase creates a new register use case
func NewRegisterUseCase(
	userRepo user.Repository,
	passwordMgr *auth.PasswordManager,
	jwtMgr *auth.JWTManager,
	validator *validator.Validator,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:    userRepo,
		passwordMgr: passwordMgr,
		jwtMgr:      jwtMgr,
		validator:   validator,
	}
}

// TokenPair represents access and refresh tokens with user info
type TokenPair struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	User         *user.User `json:"user"`
}

// Execute executes the register use case
func (uc *RegisterUseCase) Execute(ctx context.Context, input user.CreateUserInput) (*TokenPair, error) {
	// TODO: 1. Validate input
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	// TODO: 2. Check if user with email already exists
	exists, err := uc.userRepo.Exists(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.Conflict("email already registered")
	}

	// TODO: 3. Hash password
	passwordHash, err := uc.passwordMgr.HashPassword(input.Password)
	if err != nil {
		return nil, errors.InternalServer("failed to hash password", err)
	}

	// TODO: 4. Create user entity
	newUser := &user.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: passwordHash,
		Name:         input.Name,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// TODO: 5. Save user to database
	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// TODO: 6. Generate JWT tokens
	accessToken, refreshToken, err := uc.jwtMgr.GenerateTokenPair(newUser.ID, newUser.Email)
	if err != nil {
		return nil, errors.InternalServer("failed to generate tokens", err)
	}

	// TODO: 7. Return tokens and user
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         newUser,
	}, nil
}
