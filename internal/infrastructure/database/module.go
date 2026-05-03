package database

import (
	"context"
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/infrastructure/logger"

	"go.uber.org/fx"
)

func ProvideDatabase(cfg *config.Config) (*DB, error) {
	return New(&cfg.Database)
}

// registerHooks registers lifecycle hooks for database
func registerHooks(lc fx.Lifecycle, db *DB) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Database connection pool started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing database connection pool")
			db.Close()
			logger.Info("Database connection pool closed")
			return nil
		},
	})
}
