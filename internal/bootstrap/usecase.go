package bootstrap

import (
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
	),

	// Workout use cases
	fx.Provide(
		workoutuc.NewCreateWorkoutPlanUseCase,
		workoutuc.NewAddExerciseToWorkoutUseCase,
	),

	// TODO: Add more use cases
	// Meal use cases: NewCreateFoodUseCase, NewCreateRecipeUseCase, NewCreateMealLogUseCase, etc.
	// Social use cases: NewFollowUserUseCase, NewUnfollowUserUseCase, NewCreatePostUseCase, etc.
)
