package usecase

import (
	useruc "gym-pro-2026-ptit/internal/usecase/user"
	workoutuc "gym-pro-2026-ptit/internal/usecase/workout"

	"go.uber.org/fx"
)

// Module provides use case dependencies
var Module = fx.Module("usecase",
	// User use cases
	fx.Provide(
		useruc.NewRegisterUseCase,
		useruc.NewLoginUseCase,
		useruc.NewGetProfileUseCase,
		useruc.NewUpdateProfileUseCase,
		// TODO: Add more user use cases
		// - RefreshTokenUseCase
		// - ChangePasswordUseCase
		// - etc.
	),

	// Workout use cases
	fx.Provide(
		workoutuc.NewCreateWorkoutPlanUseCase,
		workoutuc.NewAddExerciseToWorkoutUseCase,
		// TODO: Add more workout use cases
		// Exercise:
		// - ListExercisesUseCase
		// - SearchExercisesUseCase
		// - GetExerciseUseCase
		//
		// Workout Plan:
		// - GetWorkoutPlanUseCase
		// - UpdateWorkoutPlanUseCase
		// - DeleteWorkoutPlanUseCase
		// - ListWorkoutPlansUseCase
		//
		// Schedule:
		// - ScheduleWorkoutUseCase
		// - BulkScheduleWorkoutUseCase
		// - GetSchedulesUseCase
		// - UpdateScheduleUseCase
		// - DeleteScheduleUseCase
		//
		// Session:
		// - StartWorkoutSessionUseCase
		// - LogExerciseSetUseCase
		// - CompleteWorkoutSessionUseCase
		// - GetSessionHistoryUseCase
		// - GetWorkoutStatsUseCase
	),

	// TODO: Add Meal use cases
	// fx.Provide(
	//     mealuc.NewCreateFoodUseCase,
	//     mealuc.NewCreateRecipeUseCase,
	//     mealuc.NewCreateMealLogUseCase,
	//     mealuc.NewAddItemToMealLogUseCase,
	//     mealuc.NewGetDailySummaryUseCase,
	//     mealuc.NewGetNutritionStatsUseCase,
	//     // ... more meal use cases
	// ),

	// TODO: Add Social use cases
	// fx.Provide(
	//     socialuc.NewFollowUserUseCase,
	//     socialuc.NewUnfollowUserUseCase,
	//     socialuc.NewCreatePostUseCase,
	//     socialuc.NewLikePostUseCase,
	//     socialuc.NewCommentPostUseCase,
	//     socialuc.NewGetFeedUseCase,
	//     // ... more social use cases
	// ),
)
