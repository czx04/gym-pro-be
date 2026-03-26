package middleware

import (
	userdomain "gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequireRole returns a middleware that checks if the authenticated user has one of the required roles.
// Must be used after the AuthMiddleware.
func RequireRole(roles ...userdomain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUser(c)
		if err != nil {
			response.Error(c, errors.Unauthorized("user not authenticated"))
			c.Abort()
			return
		}

		if !user.IsActive {
			response.Error(c, errors.Forbidden("account is deactivated"))
			c.Abort()
			return
		}

		for _, role := range roles {
			if user.HasRole(role) {
				c.Next()
				return
			}
		}

		response.Error(c, errors.Forbidden("insufficient permissions"))
		c.Abort()
	}
}

// RequireAdmin is a convenience middleware that only allows admin users.
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(userdomain.RoleAdmin)
}

// RequireActiveUser ensures the authenticated user's account is active.
func RequireActiveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUser(c)
		if err != nil {
			response.Error(c, errors.Unauthorized("user not authenticated"))
			c.Abort()
			return
		}

		if !user.IsActive {
			response.Error(c, errors.Forbidden("account is deactivated"))
			c.Abort()
			return
		}

		c.Next()
	}
}
