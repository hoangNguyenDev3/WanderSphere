package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/authpost"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/utils"
	pb "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	"google.golang.org/grpc"
)

func main() {
	// Flags - use environment-aware default config path
	defaultPath := os.Getenv("CONFIG_PATH")
	if defaultPath == "" {
		defaultPath = "/app/config.yaml"
	}
	cfgPath := flag.String("conf", defaultPath, "Path to config file for this service")
	flag.Parse()

	// Load configurations - use the function that extracts the specific section
	cfg, err := configs.GetAuthenticateAndPostConfig(*cfgPath)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	// Start new authenticate and post service
	service, err := authpost.NewAuthenticateAndPostService(cfg)
	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}

	// Create health checker
	healthChecker := utils.NewHealthChecker("authpost", "1.0.0", service.GetLogger())

	// Setup HTTP server for health checks (on port +100)
	healthPort := cfg.Port + 100
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", healthChecker.HealthHandler())
	healthMux.HandleFunc("/health/detailed", healthChecker.DetailedHealthHandler(service.GetDB(), nil))

	healthServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", healthPort),
		Handler: healthMux,
	}

	// Start health check server in background
	go func() {
		log.Printf("Starting AuthPost health check server on port %d", healthPort)
		if err := healthServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Health server error: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		log.Fatalf("can not listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthenticateAndPostServer(grpcServer, service)

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down AuthPost service...")
		grpcServer.GracefulStop()
		healthServer.Close()
		log.Println("AuthPost service stopped")
		os.Exit(0)
	}()

	log.Printf("Starting AuthPost service on port %d", cfg.Port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
