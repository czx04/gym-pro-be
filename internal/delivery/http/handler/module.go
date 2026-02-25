package handler

import (
	"go.uber.org/fx"
)

// Module provides handler dependencies
var Module = fx.Module("handler",
	fx.Provide(
		NewAuthHandler,
		NewWorkoutHandler,
		// TODO: Add more handlers as you implement them
		// NewMealHandler,
		// NewSocialHandler,
	),
)
