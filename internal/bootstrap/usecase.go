package bootstrap

import (
	exerciseuc "gym-pro-2026-ptit/internal/usecase/exercise"
	useruc "gym-pro-2026-ptit/internal/usecase/user"
	workoutuc "gym-pro-2026-ptit/internal/usecase/workout"

	"go.uber.org/fx"
)

var UseCaseProviders = fx.Options(
	// User use cases
	fx.Provide(
		useruc.NewRegisterRequestOTPUseCase,
		useruc.NewVerifyOTPUseCase,
		useruc.NewLoginUseCase,
		useruc.NewGetProfileUseCase,
		useruc.NewUpdateProfileUseCase,
		useruc.NewRefreshTokenUseCase,
	),

	// Workout use cases
	fx.Provide(
		workoutuc.NewCreateWorkoutPlanUseCase,
		workoutuc.NewAddExerciseToWorkoutUseCase,
	),

	// Exercise use cases
	fx.Provide(
		exerciseuc.NewListExercisesUseCase,
		exerciseuc.NewGetExerciseUseCase,
		exerciseuc.NewFilterExerciseUseCase,
	),

	// TODO: Add more use cases
	// Meal use cases: NewCreateFoodUseCase, NewCreateRecipeUseCase, NewCreateMealLogUseCase, etc.
	// Social use cases: NewFollowUserUseCase, NewUnfollowUserUseCase, NewCreatePostUseCase, etc.
)
