package bootstrap

import (
	"gym-pro-2026-ptit/internal/delivery/http/handler"

	"go.uber.org/fx"
)

var HandlerProviders = fx.Options(
	fx.Provide(
		handler.NewAuthHandler,
		handler.NewWorkoutHandler,
		handler.NewExerciseHandler,
		handler.NewFoodHandler,
		handler.NewRecipeHandler,
		handler.NewMealLogHandler,
		handler.NewUserHandler,
		// TODO: Add more handlers
		// handler.NewMealHandler,
		// handler.NewSocialHandler,
	),
)
