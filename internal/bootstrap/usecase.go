package bootstrap

import (
	exerciseuc "gym-pro-2026-ptit/internal/usecase/exercise"
	mealuc "gym-pro-2026-ptit/internal/usecase/meal"
	socialuc "gym-pro-2026-ptit/internal/usecase/social"
	useruc "gym-pro-2026-ptit/internal/usecase/user"
	workoutuc "gym-pro-2026-ptit/internal/usecase/workout"

	"go.uber.org/fx"
)

var UseCaseProviders = fx.Options(
	fx.Provide(useruc.NewUserUseCases),
	fx.Provide(workoutuc.NewWorkoutUseCases),
	fx.Provide(exerciseuc.NewExerciseUseCases),
	fx.Provide(mealuc.NewFoodUseCases),
	fx.Provide(mealuc.NewRecipeUseCases),
	fx.Provide(mealuc.NewMealStreakUseCases),
	fx.Provide(mealuc.NewPushTokenUseCases),
	fx.Provide(mealuc.NewMealLogUseCases),
	fx.Provide(mealuc.NewMealDailyUseCases),
	fx.Provide(socialuc.NewSocialUseCases),
)
