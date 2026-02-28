package database

import (
	"context"
	"fmt"
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// DB wraps the pgxpool connection
type DB struct {
	*pgxpool.Pool
}

// New creates a new database connection pool
func New(cfg *config.DatabaseConfig, log logger.Logger) (*DB, error) {
	log.Info("Connecting to PostgreSQL database")

	// Create connection pool config
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set pool configuration
	poolConfig.MaxConns = int32(cfg.MaxConnections)
	poolConfig.MinConns = int32(cfg.MaxIdleConnections)
	poolConfig.MaxConnLifetime = time.Duration(cfg.MaxLifetimeMinutes) * time.Minute
	poolConfig.MaxConnIdleTime = 10 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	// Create connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Successfully connected to PostgreSQL database",
		zap.Int("max_connections", cfg.MaxConnections),
		zap.Int("max_idle_connections", cfg.MaxIdleConnections),
	)

	return &DB{Pool: pool}, nil
}

// Close closes the database connection pool
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// HealthCheck checks if the database is reachable
func (db *DB) HealthCheck(ctx context.Context) error {
	return db.Ping(ctx)
}

// GetStats returns database pool statistics
func (db *DB) GetStats() *pgxpool.Stat {
	return db.Stat()
}
