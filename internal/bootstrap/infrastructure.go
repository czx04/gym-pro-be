package bootstrap

import (
	"context"
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/infrastructure/auth"
	"gym-pro-2026-ptit/internal/infrastructure/cache"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/internal/infrastructure/email"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/internal/infrastructure/otp"

	"go.uber.org/fx"
)

func ProvideLogger(cfg *config.Config) (logger.Logger, error) {
	return logger.New(&cfg.Logger)
}

func ProvideDatabase(cfg *config.Config, log logger.Logger) (*database.DB, error) {
	return database.New(&cfg.Database, log)
}

func ProvideCache(cfg *config.Config, log logger.Logger) *cache.Cache {
	return cache.NewCache(&cfg.Cache, log)
}

func ProvideJWTManager(cfg *config.Config) *auth.JWTManager {
	return auth.NewJWTManager(&cfg.JWT)
}

func ProvidePasswordManager() *auth.PasswordManager {
	return auth.NewPasswordManager()
}

func ProvideOTPService(cache *cache.Cache, log logger.Logger) otp.Service {
	return otp.NewOTPService(cache, log)
}

func ProvideEmailService(cfg *config.Config, log logger.Logger) email.Service {
	return email.NewEmailService(&cfg.Email, log)
}

func RegisterInfrastructureHooks(lc fx.Lifecycle, db *database.DB, cache *cache.Cache, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Database connection pool started")
			log.Info("Redis cache started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Closing Redis cache")
			cache.Close()
			log.Info("Redis cache closed")

			log.Info("Closing database connection pool")
			db.Close()
			log.Info("Database connection pool closed")
			return nil
		},
	})
}
