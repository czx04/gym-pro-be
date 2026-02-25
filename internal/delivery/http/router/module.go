package router

import (
	"context"
	"fmt"
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides router dependency
var Module = fx.Module("router",
	fx.Provide(New),
	fx.Invoke(registerHooks),
)

// registerHooks registers lifecycle hooks for HTTP server
func registerHooks(lc fx.Lifecycle, router *Router, cfg *config.Config, log logger.Logger) {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: router.GetEngine(),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting HTTP server",
				zap.String("host", cfg.Server.Host),
				zap.String("port", cfg.Server.Port),
				zap.String("mode", cfg.Server.GinMode),
			)

			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal("Failed to start HTTP server", zap.Error(err))
				}
			}()

			log.Info("HTTP server started successfully",
				zap.String("address", server.Addr),
			)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down HTTP server")
			if err := server.Shutdown(ctx); err != nil {
				log.Error("Failed to shutdown HTTP server gracefully", zap.Error(err))
				return err
			}
			log.Info("HTTP server stopped")
			return nil
		},
	})
}
