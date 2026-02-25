package auth

import (
	"gym-pro-2026-ptit/internal/config"

	"go.uber.org/fx"
)

// Module provides auth dependencies
var Module = fx.Module("auth",
	fx.Provide(
		ProvideJWTManager,
		ProvidePasswordManager,
	),
)

// ProvideJWTManager creates a new JWT manager
func ProvideJWTManager(cfg *config.Config) *JWTManager {
	return NewJWTManager(&cfg.JWT)
}

// ProvidePasswordManager creates a new password manager
func ProvidePasswordManager() *PasswordManager {
	return NewPasswordManager()
}
