package repository

import (
	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/internal/domain/social"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/domain/workout"
	"gym-pro-2026-ptit/internal/repository/postgres"

	"go.uber.org/fx"
)

// Module provides repository dependencies
var Module = fx.Module("repository",
	// User repositories
	fx.Provide(
		fx.Annotate(
			postgres.NewUserRepository,
			fx.As(new(user.Repository)),
		),
	),

	// Workout repositories
	fx.Provide(
		fx.Annotate(
			postgres.NewExerciseRepository,
			fx.As(new(workout.ExerciseRepository)),
		),
		fx.Annotate(
			postgres.NewWorkoutPlanRepository,
			fx.As(new(workout.WorkoutPlanRepository)),
		),
		fx.Annotate(
			postgres.NewWorkoutScheduleRepository,
			fx.As(new(workout.WorkoutScheduleRepository)),
		),
		fx.Annotate(
			postgres.NewWorkoutSessionRepository,
			fx.As(new(workout.WorkoutSessionRepository)),
		),
	),

	// Meal repositories
	fx.Provide(
		fx.Annotate(
			postgres.NewFoodRepository,
			fx.As(new(meal.FoodRepository)),
		),
		fx.Annotate(
			postgres.NewRecipeRepository,
			fx.As(new(meal.RecipeRepository)),
		),
		fx.Annotate(
			postgres.NewMealLogRepository,
			fx.As(new(meal.MealLogRepository)),
		),
	),

	// Social repositories
	fx.Provide(
		fx.Annotate(
			postgres.NewFollowRepository,
			fx.As(new(social.FollowRepository)),
		),
		fx.Annotate(
			postgres.NewPostRepository,
			fx.As(new(social.PostRepository)),
		),
		fx.Annotate(
			postgres.NewLikeRepository,
			fx.As(new(social.LikeRepository)),
		),
		fx.Annotate(
			postgres.NewCommentRepository,
			fx.As(new(social.CommentRepository)),
		),
	),
)
