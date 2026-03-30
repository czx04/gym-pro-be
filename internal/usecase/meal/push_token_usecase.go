package meal

import (
	"context"

	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"

	"github.com/google/uuid"
)

// PushTokenUseCases registers Expo push tokens for meal reminders.
type PushTokenUseCases struct {
	repo      meal.PushTokenRepository
	validator *validator.Validator
}

func NewPushTokenUseCases(repo meal.PushTokenRepository, validator *validator.Validator) *PushTokenUseCases {
	return &PushTokenUseCases{repo: repo, validator: validator}
}

func (uc *PushTokenUseCases) Register(ctx context.Context, userID uuid.UUID, input meal.RegisterPushTokenInput) error {
	if err := uc.validator.Validate(input); err != nil {
		return errors.Validation(err.Error())
	}
	if err := uc.repo.Upsert(ctx, userID, input.ExpoPushToken, input.Platform); err != nil {
		return errors.DatabaseError("failed to save push token", err)
	}
	return nil
}

func (uc *PushTokenUseCases) Delete(ctx context.Context, userID uuid.UUID, token string) error {
	if token == "" {
		return errors.BadRequest("expo_push_token is required")
	}
	if err := uc.repo.DeleteByUserAndToken(ctx, userID, token); err != nil {
		return errors.DatabaseError("failed to remove push token", err)
	}
	return nil
}
