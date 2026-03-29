package bootstrap

import (
	"context"
	"gym-pro-2026-ptit/internal/delivery/http/router"
	"gym-pro-2026-ptit/internal/delivery/ws"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/internal/port/socialnotify"

	"go.uber.org/fx"
)

func ProvideSocialBroadcaster(h *ws.Hub) socialnotify.Broadcaster {
	return h
}

func RegisterSocialWebSocket(r *router.Router, h *ws.Handler) {
	r.GetEngine().GET("/api/v1/social/ws", h.Handle)
}

func RegisterSocialWebSocketHook(lc fx.Lifecycle, r *router.Router, h *ws.Handler) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			RegisterSocialWebSocket(r, h)
			logger.Info("social websocket route registered", "path", "/api/v1/social/ws")
			return nil
		},
	})
}
