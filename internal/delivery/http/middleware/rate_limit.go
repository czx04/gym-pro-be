package middleware

import (
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    requestsPerMinute,
		window:   time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if a request is allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Get existing requests for this key
	requests, exists := rl.requests[key]
	if !exists {
		rl.requests[key] = []time.Time{now}
		return true
	}

	// Filter out requests outside the window
	var validRequests []time.Time
	for _, req := range requests {
		if req.After(windowStart) {
			validRequests = append(validRequests, req)
		}
	}

	// Check if limit is exceeded
	if len(validRequests) >= rl.limit {
		rl.requests[key] = validRequests
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[key] = validRequests
	return true
}

// cleanup removes old entries periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		windowStart := now.Add(-rl.window)

		for key, requests := range rl.requests {
			var validRequests []time.Time
			for _, req := range requests {
				if req.After(windowStart) {
					validRequests = append(validRequests, req)
				}
			}

			if len(validRequests) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validRequests
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware creates rate limiting middleware
func RateLimitMiddleware(cfg *config.RateLimitConfig) gin.HandlerFunc {
	limiter := NewRateLimiter(cfg.RequestsPerMinute)

	return func(c *gin.Context) {
		// Use IP address as key, or user ID if authenticated
		key := c.ClientIP()
		if userID, exists := c.Get(UserIDKey); exists {
			key = userID.(string)
		}

		if !limiter.Allow(key) {
			response.Error(c, errors.NewAppError(
				"RATE_LIMIT_EXCEEDED",
				"too many requests, please try again later",
				429,
				nil,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
