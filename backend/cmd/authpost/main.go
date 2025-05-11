package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/authpost"
	pb "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	"google.golang.org/grpc"
)

func main() {
	// Flags
	cfgPath := flag.String("conf", "configs/files/local.yml", "Path to config file for this service")
	flag.Parse()

	// Load configurations
	cfg, err := configs.GetAuthenticateAndPostConfigDirect(*cfgPath)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	// Start new authenticate and post service
	service, err := authpost.NewAuthenticateAndPostService(cfg)
	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		log.Fatalf("can not listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthenticateAndPostServer(grpcServer, service)

	log.Printf("Starting AuthPost service on port %d", cfg.Port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
