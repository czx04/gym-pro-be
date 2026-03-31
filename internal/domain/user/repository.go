package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository defines the interface for user data access
type Repository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// GetByOAuth retrieves a user by OAuth provider and ID
	GetByOAuth(ctx context.Context, provider, oauthID string) (*User, error)

	// Update updates a user
	Update(ctx context.Context, user *User) error

	// Delete deletes a user
	Delete(ctx context.Context, id uuid.UUID) error

	// UpdateProfile updates user profile information
	UpdateProfile(ctx context.Context, id uuid.UUID, input UpdateProfileInput) error

	// UpdatePassword updates user password
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error

	// UpdateEmail updates user's email (must remain unique).
	UpdateEmail(ctx context.Context, id uuid.UUID, email string) error

	// Exists checks if a user with given email exists
	Exists(ctx context.Context, email string) (bool, error)

	// InsertWeightHistory inserts a weight history record for a user
	InsertWeightHistory(ctx context.Context, item *WeightHistory) error

	// GetLatestWeightInRange gets the latest weight entry in [start, end]
	GetLatestWeightInRange(ctx context.Context, userID uuid.UUID, start, end time.Time) (*WeightHistory, error)

	// GetLatestWeightBefore gets the latest weight entry before a given timestamp
	GetLatestWeightBefore(ctx context.Context, userID uuid.UUID, before time.Time) (*WeightHistory, error)

	// ListWeightHistoryByGranularity returns one point per bucket (latest measured_at in that bucket).
	ListWeightHistoryByGranularity(ctx context.Context, userID uuid.UUID, from, to time.Time, tz string, granularity WeightHistoryGranularity) ([]WeightHistoryPoint, error)
}
