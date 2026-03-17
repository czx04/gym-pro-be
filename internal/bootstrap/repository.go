package bootstrap

import (
	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/internal/domain/social"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/domain/workout"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/internal/repository/postgres"

	"go.uber.org/fx"
)

// User repositories
func ProvideUserRepository(db *database.DB) user.Repository {
	return postgres.NewUserRepository(db)
}

// Workout repositories
func ProvideExerciseRepository(db *database.DB) workout.ExerciseRepository {
	return postgres.NewExerciseRepository(db)
}

func ProvideWorkoutPlanRepository(db *database.DB) workout.WorkoutPlanRepository {
	return postgres.NewWorkoutPlanRepository(db)
}

func ProvideWorkoutScheduleRepository(db *database.DB) workout.WorkoutScheduleRepository {
	return postgres.NewWorkoutScheduleRepository(db)
}

func ProvideWorkoutSessionRepository(db *database.DB) workout.WorkoutSessionRepository {
	return postgres.NewWorkoutSessionRepository(db)
}

// Meal repositories
func ProvideFoodRepository(db *database.DB) meal.FoodRepository {
	return postgres.NewFoodRepository(db)
}

func ProvideRecipeRepository(db *database.DB) meal.RecipeRepository {
	return postgres.NewRecipeRepository(db)
}

func ProvideMealLogRepository(db *database.DB) meal.MealLogRepository {
	return postgres.NewMealLogRepository(db)
}

// Social repositories
func ProvideFollowRepository(db *database.DB) social.FollowRepository {
	return postgres.NewFollowRepository(db)
}

func ProvidePostRepository(db *database.DB) social.PostRepository {
	return postgres.NewPostRepository(db)
}

func ProvideLikeRepository(db *database.DB) social.LikeRepository {
	return postgres.NewLikeRepository(db)
}

func ProvideCommentRepository(db *database.DB) social.CommentRepository {
	return postgres.NewCommentRepository(db)
}

func ProvideMediaAssetRepository(db *database.DB) social.MediaAssetRepository {
	return postgres.NewMediaAssetRepository(db)
}

// RepositoryProviders returns all repository providers
var RepositoryProviders = fx.Options(
	fx.Provide(
		ProvideUserRepository,
		ProvideExerciseRepository,
		ProvideWorkoutPlanRepository,
		ProvideWorkoutScheduleRepository,
		ProvideWorkoutSessionRepository,
		ProvideFoodRepository,
		ProvideRecipeRepository,
		ProvideMealLogRepository,
		ProvideFollowRepository,
		ProvidePostRepository,
		ProvideLikeRepository,
		ProvideCommentRepository,
		ProvideMediaAssetRepository,
	),
)
