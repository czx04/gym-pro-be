package middleware

import (
	"gym-pro-2026-ptit/internal/infrastructure/auth"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "user_id"
	UserEmailKey        = "user_email"
)

// AuthMiddleware creates JWT authentication middleware
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Error(c, errors.Unauthorized("missing authorization header"))
			c.Abort()
			return
		}

		// Check if it starts with Bearer
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			response.Error(c, errors.Unauthorized("invalid authorization header format"))
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			response.Error(c, errors.Unauthorized("missing token"))
			c.Abort()
			return
		}

		// Validate token
		claims, err := jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}

		// Set user info in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)

		c.Next()
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return uuid.Nil, errors.Unauthorized("user not authenticated")
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.Unauthorized("invalid user ID")
	}

	return id, nil
}

// GetUserEmail retrieves user email from context
func GetUserEmail(c *gin.Context) (string, error) {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return "", errors.Unauthorized("user not authenticated")
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", errors.Unauthorized("invalid user email")
	}

	return emailStr, nil
}

// MustGetUserID retrieves user ID from context or panics
func MustGetUserID(c *gin.Context) uuid.UUID {
	id, err := GetUserID(c)
	if err != nil {
		panic(err)
	}
	return id
}
