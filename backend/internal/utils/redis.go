package utils

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"go.uber.org/zap"
)

// RedisPool represents an enhanced Redis client pool
type RedisPool struct {
	Client *redis.Client
	Logger *zap.Logger
}

// NewRedisClient creates a new Redis client with enhanced connection pooling and error handling
func NewRedisClient(cfg *configs.RedisConfig) (*redis.Client, error) {
	// Create Redis client with enhanced options for better performance
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,

		// Enhanced connection pooling
		PoolSize:     20, // Increased from 10 to 20 for better concurrency
		MinIdleConns: 10, // Increased from 5 to 10 to reduce connection overhead
		MaxRetries:   5,  // Increased from 3 to 5 for better resilience

		// Connection timeouts optimized for microservices
		DialTimeout:  10 * time.Second, // Increased from 5s for better reliability
		ReadTimeout:  5 * time.Second,  // Increased from 3s for complex operations
		WriteTimeout: 5 * time.Second,  // Increased from 3s for better write performance
		PoolTimeout:  10 * time.Second, // Increased from 4s to reduce pool timeout errors

		// Connection health management
		IdleTimeout:        300 * time.Second, // Close idle connections after 5 minutes
		IdleCheckFrequency: 60 * time.Second,  // Check for idle connections every minute

		// Retry configuration
		MaxRetryBackoff: 512 * time.Millisecond,
		MinRetryBackoff: 8 * time.Millisecond,
	})

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}

// NewRedisPool creates a new Redis connection pool with enhanced features
func NewRedisPool(cfg *configs.RedisConfig, logger *zap.Logger) (*RedisPool, error) {
	client, err := NewRedisClient(cfg)
	if err != nil {
		return nil, err
	}

	pool := &RedisPool{
		Client: client,
		Logger: logger,
	}

	// Start connection monitoring
	go pool.monitorConnections()

	return pool, nil
}

// monitorConnections periodically logs Redis connection pool statistics
func (rp *RedisPool) monitorConnections() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		stats := rp.Client.PoolStats()
		rp.Logger.Info("Redis connection pool stats",
			zap.Uint32("hits", stats.Hits),
			zap.Uint32("misses", stats.Misses),
			zap.Uint32("timeouts", stats.Timeouts),
			zap.Uint32("total_conns", stats.TotalConns),
			zap.Uint32("idle_conns", stats.IdleConns),
			zap.Uint32("stale_conns", stats.StaleConns),
		)
	}
}

// HealthCheck performs a comprehensive health check on the Redis connection
func (rp *RedisPool) HealthCheck(ctx context.Context) error {
	// Test basic connectivity
	if err := rp.Client.Ping(ctx).Err(); err != nil {
		return err
	}

	// Test read/write operations
	testKey := "health_check_" + time.Now().Format("20060102150405")
	if err := rp.Client.Set(ctx, testKey, "test", time.Minute).Err(); err != nil {
		return err
	}

	if err := rp.Client.Del(ctx, testKey).Err(); err != nil {
		rp.Logger.Warn("Failed to delete health check key", zap.Error(err))
	}

	return nil
}

// Close gracefully closes the Redis connection pool
func (rp *RedisPool) Close() error {
	rp.Logger.Info("Closing Redis connection pool")
	return rp.Client.Close()
}

// Pipeline creates a new Redis pipeline for batch operations
func (rp *RedisPool) Pipeline() redis.Pipeliner {
	return rp.Client.Pipeline()
}

// TxPipeline creates a new Redis transaction pipeline
func (rp *RedisPool) TxPipeline() redis.Pipeliner {
	return rp.Client.TxPipeline()
}
