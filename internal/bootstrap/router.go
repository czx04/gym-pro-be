package bootstrap

import (
	"context"
	"fmt"
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/delivery/http/handler"
	"gym-pro-2026-ptit/internal/delivery/http/router"
	"gym-pro-2026-ptit/internal/infrastructure/auth"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ProvideRouter creates a new router instance
func ProvideRouter(
	cfg *config.Config,
	log logger.Logger,
	jwtManager *auth.JWTManager,
	authHandler *handler.AuthHandler,
	workoutHandler *handler.WorkoutHandler,
) *router.Router {
	return router.New(cfg, log, jwtManager, authHandler, workoutHandler)
}

// RegisterRouterHooks registers lifecycle hooks for HTTP server
func RegisterRouterHooks(lc fx.Lifecycle, r *router.Router, cfg *config.Config, log logger.Logger) {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: r.GetEngine(),
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
