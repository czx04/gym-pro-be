package main

import (
	"context"
	"fmt"
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/delivery/http/handler"
	"gym-pro-2026-ptit/internal/delivery/http/router"
	"gym-pro-2026-ptit/internal/infrastructure/auth"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/internal/repository"
	"gym-pro-2026-ptit/internal/usecase"
	"gym-pro-2026-ptit/pkg/validator"
	"os"

	"go.uber.org/fx"
)

// @title Gym Pro API
// @version 1.0
// @description Backend API for Gym Pro - A fitness tracking mobile application
// @contact.name API Support
// @contact.email support@gympro.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	app := fx.New(
		// Configuration
		fx.Provide(loadConfig),

		// Utilities
		fx.Provide(validator.New),

		// Infrastructure
		logger.Module,
		database.Module,
		auth.Module,

		// Data Layer
		repository.Module,

		// Business Logic Layer
		usecase.Module,

		// HTTP Layer
		handler.Module,
		router.Module,

		// Lifecycle logging
		fx.Invoke(logLifecycle),
	)

	app.Run()
}

// loadConfig loads application configuration
func loadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return cfg, nil
}

// logLifecycle logs application lifecycle events
func logLifecycle(lc fx.Lifecycle, log logger.Logger) {
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
				// Ignore sync errors on shutdown
				fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
			}
			return nil
		},
	})
}
