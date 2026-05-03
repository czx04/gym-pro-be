package bootstrap

import (
	"context"
	"gym-pro-2026-ptit/internal/delivery/http/websocket"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/infrastructure/auth"

	"go.uber.org/fx"
)

func ProvideWebSocketHub(jwtManager *auth.JWTManager, userRepo user.Repository) *websocket.Hub {
	return websocket.NewHub(jwtManager, userRepo)
}

func RegisterWebSocketHooks(lc fx.Lifecycle, hub *websocket.Hub) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			hub.Shutdown()
			return nil
		},
	})
}
