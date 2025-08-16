package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// RateLimiterConfig holds rate limiting configuration
type RateLimiterConfig struct {
	RequestsPerMinute int
	BurstSize         int // Extra requests allowed above the per-minute rate
	Enabled           bool
}

// DefaultRateLimiterConfig returns sensible defaults
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerMinute: 100,
		BurstSize:         10,
		Enabled:           true,
	}
}

// RateLimiter implements Redis-backed per-user rate limiting
type RateLimiter struct {
	redisClient *redis.Client
	config      RateLimiterConfig
	logger      *zap.Logger
	keyPrefix   string
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redisClient *redis.Client, config RateLimiterConfig, logger *zap.Logger) *RateLimiter {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return &RateLimiter{
		redisClient: redisClient,
		config:      config,
		logger:      logger,
		keyPrefix:   "ratelimit:",
	}
}

// Middleware returns a Gin middleware handler function
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.config.Enabled {
			c.Next()
			return
		}

		// Determine the rate limit key
		// Use user_id if available (set by auth middleware), otherwise use client IP
		key := rl.buildKey(c)

		limit := int64(rl.config.RequestsPerMinute + rl.config.BurstSize)
		window := time.Minute

		// Use Redis TxPipeline for atomic INCR + EXPIRE
		ctx := context.Background()
		pipe := rl.redisClient.TxPipeline()
		incr := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, window)
		_, err := pipe.Exec(ctx)
		if err != nil {
			rl.logger.Error("Rate limiter Redis error, allowing request",
				zap.String("key", key),
				zap.Error(err))
			c.Next() // Fail open
			return
		}

		count := incr.Val()

		// Get TTL for reset header, with defensive fallback
		ttl, err := rl.redisClient.TTL(ctx, key).Result()
		if err != nil || ttl <= 0 {
			ttl = window // Fall back to full window if TTL is invalid
		}
		resetTime := time.Now().Add(ttl).Unix()

		remaining := limit - count
		if remaining < 0 {
			remaining = 0
		}

		// Set rate limit headers on every response
		c.Header("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

		if count > limit {
			rl.logger.Warn("Rate limit exceeded",
				zap.String("key", key),
				zap.Int64("count", count),
				zap.Int64("limit", limit))

			c.Header("Retry-After", strconv.FormatInt(int64(ttl.Seconds()), 10))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":               "rate limit exceeded",
				"message":             "Too many requests. Please try again later.",
				"retry_after_seconds": int(ttl.Seconds()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// buildKey constructs the rate limit key based on user identity
func (rl *RateLimiter) buildKey(c *gin.Context) string {
	// Try to get user_id from context (set by AuthRequired middleware or RefreshSession)
	userID, exists := c.Get("user_id")
	if exists {
		return fmt.Sprintf("%suser:%v:%s:%s", rl.keyPrefix, userID, c.Request.Method, c.FullPath())
	}
	// Fall back to client IP for unauthenticated requests
	return fmt.Sprintf("%sip:%s:%s:%s", rl.keyPrefix, c.ClientIP(), c.Request.Method, c.FullPath())
}
