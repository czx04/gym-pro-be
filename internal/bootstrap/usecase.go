package bootstrap

import (
	adminuc "gym-pro-2026-ptit/internal/usecase/admin"
	exerciseuc "gym-pro-2026-ptit/internal/usecase/exercise"
	mealuc "gym-pro-2026-ptit/internal/usecase/meal"
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
	fx.Provide(mealuc.NewMealLogUseCases),
	fx.Provide(adminuc.NewAdminUseCases),
)
