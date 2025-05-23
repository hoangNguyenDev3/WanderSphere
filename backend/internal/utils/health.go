package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// HealthStatus represents the health status of a service
type HealthStatus struct {
	Status       string            `json:"status"`
	Timestamp    time.Time         `json:"timestamp"`
	Version      string            `json:"version,omitempty"`
	ServiceName  string            `json:"service_name"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
	Uptime       string            `json:"uptime,omitempty"`
}

// HealthChecker provides health check functionality
type HealthChecker struct {
	ServiceName string
	Version     string
	StartTime   time.Time
	Logger      *zap.Logger
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(serviceName, version string, logger *zap.Logger) *HealthChecker {
	return &HealthChecker{
		ServiceName: serviceName,
		Version:     version,
		StartTime:   time.Now(),
		Logger:      logger,
	}
}

// HealthHandler creates an HTTP handler for health checks
func (hc *HealthChecker) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := hc.GetHealthStatus()

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode and send response
		if err := json.NewEncoder(w).Encode(status); err != nil {
			hc.Logger.Error("Failed to encode health response", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		hc.Logger.Debug("Health check requested", zap.String("status", status.Status))
	}
}

// GetHealthStatus returns the current health status
func (hc *HealthChecker) GetHealthStatus() HealthStatus {
	uptime := time.Since(hc.StartTime).String()

	return HealthStatus{
		Status:       "healthy",
		Timestamp:    time.Now(),
		Version:      hc.Version,
		ServiceName:  hc.ServiceName,
		Uptime:       uptime,
		Dependencies: make(map[string]string),
	}
}

// CheckDatabaseHealth checks if database connection is healthy
func (hc *HealthChecker) CheckDatabaseHealth(db *gorm.DB) string {
	if db == nil {
		return "unavailable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlDB, err := db.DB()
	if err != nil {
		hc.Logger.Warn("Database health check failed - DB() error", zap.Error(err))
		return "unhealthy"
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		hc.Logger.Warn("Database health check failed - ping error", zap.Error(err))
		return "unhealthy"
	}

	return "healthy"
}

// CheckRedisHealth checks if Redis connection is healthy
func (hc *HealthChecker) CheckRedisHealth(client *redis.Client) string {
	if client == nil {
		return "unavailable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		hc.Logger.Warn("Redis health check failed", zap.Error(err))
		return "unhealthy"
	}

	return "healthy"
}

// CheckKafkaAvailability checks Kafka availability based on service status
func (hc *HealthChecker) CheckKafkaAvailability(isAvailable bool) string {
	if isAvailable {
		return "healthy"
	}
	return "degraded"
}

// AddDependencyStatus adds a dependency status to the health check
func (hc *HealthChecker) AddDependencyStatus(dependencies map[string]string, name, status string) {
	dependencies[name] = status
}

// GetDetailedHealthStatus returns health status with dependency checks
func (hc *HealthChecker) GetDetailedHealthStatus(db *gorm.DB, redisClient *redis.Client) HealthStatus {
	status := hc.GetHealthStatus()

	// Check database health
	if db != nil {
		hc.AddDependencyStatus(status.Dependencies, "database", hc.CheckDatabaseHealth(db))
	}

	// Check Redis health
	if redisClient != nil {
		hc.AddDependencyStatus(status.Dependencies, "redis", hc.CheckRedisHealth(redisClient))
	}

	// Determine overall status based on dependencies
	for _, depStatus := range status.Dependencies {
		if depStatus == "unhealthy" {
			status.Status = "degraded"
			break
		}
	}

	return status
}

// DetailedHealthHandler creates an HTTP handler for detailed health checks
func (hc *HealthChecker) DetailedHealthHandler(db *gorm.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := hc.GetDetailedHealthStatus(db, redisClient)

		// Set response headers
		w.Header().Set("Content-Type", "application/json")

		// Set status code based on health
		statusCode := http.StatusOK
		if status.Status == "degraded" {
			statusCode = http.StatusServiceUnavailable
		}
		w.WriteHeader(statusCode)

		// Encode and send response
		if err := json.NewEncoder(w).Encode(status); err != nil {
			hc.Logger.Error("Failed to encode detailed health response", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		hc.Logger.Debug("Detailed health check requested",
			zap.String("status", status.Status),
			zap.Any("dependencies", status.Dependencies))
	}
}
