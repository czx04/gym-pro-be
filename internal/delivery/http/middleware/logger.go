package middleware

import (
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware creates request logging middleware
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		latencyMs := float64(latency.Microseconds()) / 1000.0

		msg := c.Request.Method + " " + path
		kvs := []interface{}{"method", c.Request.Method, "path", path, "status", status, "latency_ms", latencyMs, "ip", c.ClientIP()}
		if query != "" {
			kvs = append(kvs, "query", query)
		}
		if len(c.Errors) > 0 {
			kvs = append(kvs, "error", c.Errors.String())
		}

		if status >= 500 {
			logger.Error(msg, kvs...)
		} else if status >= 400 {
			logger.Warn(msg, kvs...)
		} else {
			logger.Info(msg, kvs...)
		}
	}
}
