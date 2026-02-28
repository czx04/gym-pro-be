package bootstrap

import (
	"context"
	"fmt"
	"os"

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
		),

		// Data Layer (Repositories)
		RepositoryProviders,

		// Business Logic Layer (Use Cases)
		UseCaseProviders,

		// HTTP Layer (Handlers & Router)
		HandlerProviders,
		fx.Provide(ProvideRouter),

		// Lifecycle hooks
		fx.Invoke(
			RegisterInfrastructureHooks,
			RegisterRouterHooks,
			registerAppLifecycle,
		),
	)
}

// registerAppLifecycle registers application lifecycle hooks
func registerAppLifecycle(lc fx.Lifecycle, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("===========================================")
			log.Info("🚀 Gym Pro API Server Starting...")
			log.Info("===========================================")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("===========================================")
			log.Info("👋 Gym Pro API Server Stopping...")
			log.Info("===========================================")
			if err := log.Sync(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
			}
			return nil
		},
	})
}
