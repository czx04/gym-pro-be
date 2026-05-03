package database

import (
	"context"
	"fmt"
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
	tx   pgx.Tx
}

func New(cfg *config.DatabaseConfig) (*DB, error) {
	logger.Info("Connecting to PostgreSQL database")

	poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxConnections)
	poolConfig.MinConns = int32(cfg.MaxIdleConnections)
	poolConfig.MaxConnLifetime = time.Duration(cfg.MaxLifetimeMinutes) * time.Minute
	poolConfig.MaxConnIdleTime = 10 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL database", "max_connections", cfg.MaxConnections, "max_idle_connections", cfg.MaxIdleConnections)

	return &DB{pool: pool}, nil
}

func (db *DB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if db.tx != nil {
		return db.tx.Exec(ctx, sql, args...)
	}
	return db.pool.Exec(ctx, sql, args...)
}

func (db *DB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if db.tx != nil {
		return db.tx.Query(ctx, sql, args...)
	}
	return db.pool.Query(ctx, sql, args...)
}

func (db *DB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if db.tx != nil {
		return db.tx.QueryRow(ctx, sql, args...)
	}
	return db.pool.QueryRow(ctx, sql, args...)
}

func (db *DB) Begin(ctx context.Context) (*DB, error) {
	if db.tx != nil {
		return nil, fmt.Errorf("database: already in transaction")
	}
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &DB{pool: db.pool, tx: tx}, nil
}

func (db *DB) Commit(ctx context.Context) error {
	if db.tx == nil {
		return nil
	}
	return db.tx.Commit(ctx)
}

func (db *DB) Rollback(ctx context.Context) error {
	if db.tx == nil {
		return nil
	}
	return db.tx.Rollback(ctx)
}

func (db *DB) Close() {
	if db.tx != nil {
		return
	}
	if db.pool != nil {
		db.pool.Close()
	}
}

func (db *DB) Ping(ctx context.Context) error {
	if db.tx != nil {
		return nil
	}
	return db.pool.Ping(ctx)
}

func (db *DB) HealthCheck(ctx context.Context) error {
	return db.Ping(ctx)
}

func (db *DB) GetStats() *pgxpool.Stat {
	if db.pool == nil {
		return nil
	}
	return db.pool.Stat()
}
