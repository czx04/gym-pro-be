package meal

import (
	"time"

	"github.com/google/uuid"
)

// MealStreak is cached streak data for a user.
type MealStreak struct {
	UserID        uuid.UUID `json:"user_id"`
	CurrentStreak int       `json:"current_streak"`
	LongestStreak int       `json:"longest_streak"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// RegisterPushTokenInput registers or updates an Expo push token for the user.
type RegisterPushTokenInput struct {
	ExpoPushToken string `json:"expo_push_token" validate:"required"`
	Platform      string `json:"platform,omitempty"`
}
