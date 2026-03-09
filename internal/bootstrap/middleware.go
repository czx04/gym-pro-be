package bootstrap

import (
	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/infrastructure/auth"
)

// ProvideAuthMiddleware provides the JWT auth middleware for injection into router.
func ProvideAuthMiddleware(jwtManager *auth.JWTManager, userRepo user.Repository) middleware.AuthMiddleware {
	return middleware.NewAuthMiddleware(jwtManager, userRepo)
}
