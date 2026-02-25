package database

import (
	"context"
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/infrastructure/logger"

	"go.uber.org/fx"
)

// Module provides database dependency
var Module = fx.Module("database",
	fx.Provide(ProvideDatabase),
	fx.Invoke(registerHooks),
)

// ProvideDatabase creates a new database connection
func ProvideDatabase(cfg *config.Config, log logger.Logger) (*DB, error) {
	return New(&cfg.Database, log)
}

// registerHooks registers lifecycle hooks for database
func registerHooks(lc fx.Lifecycle, db *DB, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Database connection pool started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Closing database connection pool")
			db.Close()
			log.Info("Database connection pool closed")
			return nil
		},
	})
}
