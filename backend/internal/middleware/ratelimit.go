package middleware

import (
	"sync"
	"time"
	"url-shortener-backend/internal/config"
	"url-shortener-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
	
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-rl.window)
	
	if _, exists := rl.requests[key]; !exists {
		rl.requests[key] = []time.Time{}
	}
	
	var validRequests []time.Time
	for _, req := range rl.requests[key] {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}
	
	rl.requests[key] = validRequests
	
	if len(validRequests) >= rl.limit {
		return false
	}
	
	rl.requests[key] = append(rl.requests[key], now)
	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window)
		
		for key, requests := range rl.requests {
			var validRequests []time.Time
			for _, req := range requests {
				if req.After(cutoff) {
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

func RateLimitMiddleware(cfg *config.Config) fiber.Handler {
	limiter := NewRateLimiter(cfg.RateLimitRequests, time.Duration(cfg.RateLimitWindow)*time.Second)
	
	return func(c *fiber.Ctx) error {
		key := c.IP()
		
		if !limiter.Allow(key) {
			return c.Status(fiber.StatusTooManyRequests).JSON(models.ErrorResponse{
				Error:   "rate_limit_exceeded",
				Message: "Too many requests, please try again later",
			})
		}
		
		return c.Next()
	}
}