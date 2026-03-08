package user

import (
	"context"
	"gym-pro-2026-ptit/internal/infrastructure/auth"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"
)

type RefreshTokenUseCase struct {
	jwtMgr    *auth.JWTManager
	validator *validator.Validator
}

func NewRefreshTokenUseCase(
	jwtMgr *auth.JWTManager,
	validator *validator.Validator,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		jwtMgr:    jwtMgr,
		validator: validator,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, input RefreshTokenRequest) (*TokenPair, error) {
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

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}
