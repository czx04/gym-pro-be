package bootstrap

import (
	"context"
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/infrastructure/ai"
	"gym-pro-2026-ptit/internal/infrastructure/auth"
	"gym-pro-2026-ptit/internal/infrastructure/cache"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/internal/infrastructure/email"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/internal/infrastructure/otp"

	"go.uber.org/fx"
)

func ProvideLogger(cfg *config.Config) (logger.Logger, error) {
	l, err := logger.New(&cfg.Logger)
	if err != nil {
		return nil, err
	}
	logger.SetGlobal(l)
	return l, nil
}

func ProvideDatabase(cfg *config.Config) (*database.DB, error) {
	return database.New(&cfg.Database)
}

func ProvideCache(cfg *config.Config) *cache.Cache {
	return cache.NewCache(&cfg.Cache)
}

func ProvideJWTManager(cfg *config.Config) *auth.JWTManager {
	return auth.NewJWTManager(&cfg.JWT)
}

func ProvidePasswordManager() *auth.PasswordManager {
	return auth.NewPasswordManager()
}

func ProvideOTPService(cache *cache.Cache) otp.Service {
	return otp.NewOTPService(cache)
}

func ProvideEmailService(cfg *config.Config) email.Service {
	return email.NewEmailService(&cfg.Email)
}

func ProvideAIService() ai.Service {
	return ai.NewGeminiService()
}

func RegisterInfrastructureHooks(lc fx.Lifecycle, db *database.DB, cache *cache.Cache) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Database connection pool started")
			logger.Info("Redis cache started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Redis cache")
			cache.Close()
			logger.Info("Redis cache closed")

			logger.Info("Closing database connection pool")
			db.Close()
			logger.Info("Database connection pool closed")
			return nil
		},
	})
}
