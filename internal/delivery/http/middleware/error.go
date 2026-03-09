package middleware

import (
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware creates error handling middleware
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", "error", err, "path", c.Request.URL.Path, "method", c.Request.Method)
				response.Error(c, errors.InternalServer("an unexpected error occurred", nil))
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			logger.Error("Request error", "err", err, "path", c.Request.URL.Path, "method", c.Request.Method)

			if c.Writer.Written() {
				return
			}

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
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("Panic recovered", "panic", recovered, "path", c.Request.URL.Path, "method", c.Request.Method)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    errors.ErrCodeInternalServer,
				"message": "Internal server error",
			},
		})
	})
}
