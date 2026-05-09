package database

import (
	"gym-pro-2026-ptit/internal/config"
)

func ProvideDatabase(cfg *config.Config) (*DB, error) {
	return New(&cfg.Database)
}
