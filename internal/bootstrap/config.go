package bootstrap

import (
	"fmt"
	"gym-pro-2026-ptit/internal/config"
)

func LoadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return cfg, nil
}
