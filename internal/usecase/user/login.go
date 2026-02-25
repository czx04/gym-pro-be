package user

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/infrastructure/auth"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"
)

// LoginUseCase handles user login
type LoginUseCase struct {
	userRepo    user.Repository
	passwordMgr *auth.PasswordManager
	jwtMgr      *auth.JWTManager
	validator   *validator.Validator
}

// NewLoginUseCase creates a new login use case
func NewLoginUseCase(
	userRepo user.Repository,
	passwordMgr *auth.PasswordManager,
	jwtMgr *auth.JWTManager,
	validator *validator.Validator,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:    userRepo,
		passwordMgr: passwordMgr,
		jwtMgr:      jwtMgr,
		validator:   validator,
	}
}

// Execute executes the login use case
func (uc *LoginUseCase) Execute(ctx context.Context, input user.LoginInput) (*TokenPair, error) {
	// TODO: 1. Validate input
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	// TODO: 2. Find user by email
	u, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.InvalidCredentials()
	}

	// TODO: 3. Verify password
	if !uc.passwordMgr.VerifyPassword(u.PasswordHash, input.Password) {
		return nil, errors.InvalidCredentials()
	}

	// TODO: 4. Generate JWT tokens
	accessToken, refreshToken, err := uc.jwtMgr.GenerateTokenPair(u.ID, u.Email)
	if err != nil {
		return nil, errors.InternalServer("failed to generate tokens", err)
	}

	// TODO: 5. Return tokens and user
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         u,
	}, nil
}
