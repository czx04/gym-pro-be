package bootstrap

import (
	"context"
	"fmt"
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/delivery/http/handler"
	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	"gym-pro-2026-ptit/internal/delivery/http/router"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"net/http"

	"go.uber.org/fx"
)

// ProvideRouter creates a new router instance
func ProvideRouter(
	cfg *config.Config,
	authMiddleware middleware.AuthMiddleware,
	authHandler *handler.AuthHandler,
	workoutHandler *handler.WorkoutHandler,
	exerciseHandler *handler.ExerciseHandler,
	foodHandler *handler.FoodHandler,
	recipeHandler *handler.RecipeHandler,
	mealLogHandler *handler.MealLogHandler,
	userHandler *handler.UserHandler,
) *router.Router {
	return router.New(cfg, authMiddleware, authHandler, workoutHandler, exerciseHandler, foodHandler, recipeHandler, mealLogHandler, userHandler)
}

// RegisterRouterHooks registers lifecycle hooks for HTTP server
func RegisterRouterHooks(lc fx.Lifecycle, r *router.Router, cfg *config.Config) {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: r.GetEngine(),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting HTTP server", "host", cfg.Server.Host, "port", cfg.Server.Port, "mode", cfg.Server.GinMode)

			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("Failed to start HTTP server", "err", err)
				}
			}()

			logger.Info("HTTP server started successfully", "address", server.Addr)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP server")
			if err := server.Shutdown(ctx); err != nil {
				logger.Error("Failed to shutdown HTTP server gracefully", "err", err)
				return err
			}
			logger.Info("HTTP server stopped")
			return nil
		},
	})
}
