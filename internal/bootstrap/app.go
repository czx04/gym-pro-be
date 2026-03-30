package bootstrap

import (
	"context"
	"fmt"
	"os"

	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/pkg/validator"

	"go.uber.org/fx"
)

// NewApp creates and configures the fx application
func NewApp() *fx.App {
	return fx.New(
		// Configuration
		fx.Provide(LoadConfig),

		// Utilities
		fx.Provide(validator.New),

		// Infrastructure Layer
		fx.Provide(
			ProvideLogger,
			ProvideDatabase,
			ProvideCache,
			ProvideJWTManager,
			ProvidePasswordManager,
			ProvideOTPService,
			ProvideEmailService,
			ProvideAIService,
		),

		// Data Layer (Repositories)
		RepositoryProviders,

		// Business Logic Layer (Use Cases)
		UseCaseProviders,

		// HTTP Layer (Handlers & Router)
		HandlerProviders,
		fx.Provide(ProvideAuthMiddleware),
		fx.Provide(ProvideWebSocketHub),
		fx.Provide(ProvideRouter),

		// Lifecycle hooks
		fx.Invoke(
			InitGlobalLogger,
			RegisterAutoMigrateHook,
			RegisterInfrastructureHooks,
			RegisterWebSocketHooks,
			RegisterMealReminderCron,
			RegisterRouterHooks,
			RegisterAppLifecycle,
		),
	)
}

// RegisterAutoMigrateHook runs pending migrations on startup (uses same db as repositories).
func RegisterAutoMigrateHook(lc fx.Lifecycle, db *database.DB) {
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations"
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return RunAutoMigrate(ctx, db, migrationsPath)
		},
	})
}

func RegisterAppLifecycle(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("===========================================")
			logger.Info("🚀 Gym Pro API Server Starting...")
			logger.Info("===========================================")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("===========================================")
			logger.Info("👋 Gym Pro API Server Stopping...")
			logger.Info("===========================================")
			if err := logger.Sync(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
			}
			return nil
		},
	})
}

func InitGlobalLogger(l logger.Logger) {
	logger.SetGlobal(l)
}
