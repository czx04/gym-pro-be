package middleware

import (
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandlerMiddleware creates error handling middleware
func ErrorHandlerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log panic
				log.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// Return internal server error
				response.Error(c, errors.InternalServer("an unexpected error occurred", nil))
			}
		}()

		c.Next()

		// Handle errors from handlers
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Log error
			log.Error("Request error",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)

			// Check if response was already sent
			if c.Writer.Written() {
				return
			}

			// Send error response
			var appErr *errors.AppError
			if e, ok := err.(*errors.AppError); ok {
				appErr = e
			} else {
				appErr = errors.InternalServer("an unexpected error occurred", err)
			}

			c.JSON(appErr.Status, gin.H{
				"success": false,
				"error": gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				},
			})
		}
	}
}

// RecoveryMiddleware creates recovery middleware
func RecoveryMiddleware(log logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Error("Panic recovered",
			zap.Any("panic", recovered),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    errors.ErrCodeInternalServer,
				"message": "Internal server error",
			},
		})
	})
}
