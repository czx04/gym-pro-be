package logger

import (
	"gym-pro-2026-ptit/internal/config"

	"go.uber.org/fx"
)

// Module provides logger dependency
var Module = fx.Module("logger",
	fx.Provide(ProvideLogger),
)

// ProvideLogger creates a new logger instance
func ProvideLogger(cfg *config.Config) (Logger, error) {
	return New(&cfg.Logger)
}
