package middleware

import (
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware creates request logging middleware
func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		status := c.Writer.Status()

		// Log request details
		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// Add error if exists
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", c.Errors.String()))
		}

		// Log based on status code
		if status >= 500 {
			log.Error("Server error", fields...)
		} else if status >= 400 {
			log.Warn("Client error", fields...)
		} else {
			log.Info("Request processed", fields...)
		}
	}
}
