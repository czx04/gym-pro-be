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

		// Format latency for readability
		latencyMs := float64(latency.Microseconds()) / 1000.0

		// Build log fields for request logging [method, path, status, latency_ms, ip]
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Float64("latency_ms", latencyMs),
			zap.String("ip", c.ClientIP()),
		}

		// Add query
		if query != "" {
			fields = append(fields, zap.String("query", query))
		}

		// Add error if exists
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", c.Errors.String()))
		}

		// Create message
		msg := c.Request.Method + " " + path

		// Log based on status code
		if status >= 500 {
			log.Error(msg, fields...)
		} else if status >= 400 {
			log.Warn(msg, fields...)
		} else {
			log.Info(msg, fields...)
		}
	}
}
