package user

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"
	"time"

	"github.com/google/uuid"
)

// GetProfileUseCase handles getting user profile
type GetProfileUseCase struct {
	userRepo user.Repository
}

// NewGetProfileUseCase creates a new get profile use case
func NewGetProfileUseCase(userRepo user.Repository) *GetProfileUseCase {
	return &GetProfileUseCase{userRepo: userRepo}
}

// Execute executes the get profile use case
func (uc *GetProfileUseCase) Execute(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	// TODO: Get user by ID
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	u.PasswordHash = ""

	return u, nil
}

// UpdateProfileUseCase handles updating user profile
type UpdateProfileUseCase struct {
	userRepo  user.Repository
	validator *validator.Validator
}

// NewUpdateProfileUseCase creates a new update profile use case
func NewUpdateProfileUseCase(
	userRepo user.Repository,
	validator *validator.Validator,
) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		userRepo:  userRepo,
		validator: validator,
	}
}

// Execute executes the update profile use case
func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID uuid.UUID, input user.UpdateProfileInput) (*user.User, error) {
	// TODO: 1. Validate input
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	// TODO: 2. Check if user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// TODO: 3. Update user fields (only non-nil fields from input)
	// This is a simplified version - you should use UpdateProfile method
	// which does selective updates based on non-nil fields
	if err := uc.userRepo.UpdateProfile(ctx, userID, input); err != nil {
		return nil, err
	}

	// TODO: 4. Get updated user
	updated, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	updated.PasswordHash = ""
	updated.UpdatedAt = time.Now()

	return updated, nil
}
