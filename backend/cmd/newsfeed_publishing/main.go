package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	newsfeed_publishing_svc "github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/newsfeed_publishing"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/utils"
	pb_nfp "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed_publishing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Custom health handler for NFP service with Kafka status
func createNFPHealthHandler(healthChecker *utils.HealthChecker, service *newsfeed_publishing_svc.NewsfeedPublishingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := healthChecker.GetHealthStatus()

		// Add NFP-specific dependencies
		if service.GetRedis() != nil {
			healthChecker.AddDependencyStatus(status.Dependencies, "redis", healthChecker.CheckRedisHealth(service.GetRedis()))
		}
		healthChecker.AddDependencyStatus(status.Dependencies, "kafka", healthChecker.CheckKafkaAvailability(service.IsKafkaAvailable()))
		healthChecker.AddDependencyStatus(status.Dependencies, "authpost", "healthy") // Assume healthy since service started

		// Determine overall status based on dependencies
		for _, depStatus := range status.Dependencies {
			if depStatus == "unhealthy" {
				status.Status = "degraded"
				break
			}
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		statusCode := http.StatusOK
		if status.Status == "degraded" {
			statusCode = http.StatusServiceUnavailable
		}
		w.WriteHeader(statusCode)

		// Encode and send response
		if err := json.NewEncoder(w).Encode(status); err != nil {
			healthChecker.Logger.Error("Failed to encode NFP health response", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		healthChecker.Logger.Debug("NFP detailed health check requested",
			zap.String("status", status.Status),
			zap.Any("dependencies", status.Dependencies))
	}
}

func main() {
	// Flags - use environment-aware default config path
	defaultPath := os.Getenv("CONFIG_PATH")
	if defaultPath == "" {
		defaultPath = "/app/config.yaml"
	}
	cfgPath := flag.String("conf", defaultPath, "Path to config file for this service")
	flag.Parse()

	// Load configurations - use the function that extracts the specific section
	cfg, err := configs.GetNewsfeedPublishingConfig(*cfgPath)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	// Debug: print config
	cfgJSON, _ := json.MarshalIndent(cfg, "", "  ")
	log.Printf("Loaded config from %s: %s", *cfgPath, string(cfgJSON))

	// Start new newsfeed publishing service
	service, err := newsfeed_publishing_svc.NewNewsfeedPublishingService(cfg)
	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}

	// Create health checker
	healthChecker := utils.NewHealthChecker("newsfeed_publishing", "1.0.0", service.GetLogger())

	// Setup HTTP server for health checks (on port +100)
	healthPort := cfg.Port + 100
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", healthChecker.HealthHandler())
	healthMux.HandleFunc("/health/detailed", createNFPHealthHandler(healthChecker, service))

	healthServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", healthPort),
		Handler: healthMux,
	}

	// Start health check server in background
	go func() {
		log.Printf("Starting Newsfeed Publishing health check server on port %d", healthPort)
		if err := healthServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Health server error: %v", err)
		}
	}()

	// Setup graceful shutdown before starting background workers
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down Newsfeed Publishing service...")
		service.Shutdown() // Call service shutdown method
		healthServer.Close()
		log.Println("Newsfeed Publishing service stopped")
		os.Exit(0)
	}()

	// Run fanout worker
	go service.Run()

	// Start grpc server
	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("can not listen on %s: %v", addr, err)
	}

	grpcServer := grpc.NewServer()
	pb_nfp.RegisterNewsfeedPublishingServer(grpcServer, service)

	log.Printf("Starting Newsfeed Publishing service on port %d", cfg.Port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
