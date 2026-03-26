package admin

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	ListUsers(ctx context.Context, filter ListUsersFilter) ([]UserSummary, int64, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*UserDetail, error)
	UpdateUserStatus(ctx context.Context, id uuid.UUID, isActive bool) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetOverviewStats(ctx context.Context) (*OverviewStats, error)
}

type ExerciseRepository interface {
	ListExercises(ctx context.Context, filter ListExercisesFilter) ([]AdminExercise, int64, error)
	GetExerciseByID(ctx context.Context, id uuid.UUID) (*AdminExercise, error)
	CreateExercise(ctx context.Context, input CreateExerciseInput, createdBy uuid.UUID) (*AdminExercise, error)
	UpdateExercise(ctx context.Context, id uuid.UUID, input UpdateExerciseInput) (*AdminExercise, error)
	DeleteExercise(ctx context.Context, id uuid.UUID) error
}

type FoodRepository interface {
	ListFoods(ctx context.Context, filter ListFoodsFilter) ([]AdminFood, int64, error)
	GetFoodByID(ctx context.Context, id uuid.UUID) (*AdminFood, error)
	CreateSystemFood(ctx context.Context, input CreateSystemFoodInput, adminID uuid.UUID) (*AdminFood, error)
	UpdateFood(ctx context.Context, id uuid.UUID, input AdminUpdateFoodInput) (*AdminFood, error)
	DeleteFood(ctx context.Context, id uuid.UUID) error
}

type OverviewStats struct {
	TotalUsers      int64 `json:"total_users"`
	NewUsersLast30d int64 `json:"new_users_last_30d"`
	TotalExercises  int64 `json:"total_exercises"`
	ActiveExercises int64 `json:"active_exercises"`
	TotalFoods      int64 `json:"total_foods"`
	SystemFoods     int64 `json:"system_foods"`
	UserFoods       int64 `json:"user_foods"`
}
