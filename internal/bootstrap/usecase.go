package bootstrap

import (
	exerciseuc "gym-pro-2026-ptit/internal/usecase/exercise"
	useruc "gym-pro-2026-ptit/internal/usecase/user"
	workoutuc "gym-pro-2026-ptit/internal/usecase/workout"

	"go.uber.org/fx"
)

var UseCaseProviders = fx.Options(
	fx.Provide(useruc.NewUserUseCases),
	fx.Provide(workoutuc.NewWorkoutUseCases),
	fx.Provide(exerciseuc.NewExerciseUseCases),
)
