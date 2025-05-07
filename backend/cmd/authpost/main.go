package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	authpost "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	"google.golang.org/grpc"
)

func main() {
	// Flags
	cfgPath := flag.String("conf", "configs/files/test.yml", "Path to config file for this service")

	// Load configurations
	cfg, err := configs.GetAuthenticateAndPostConfig(*cfgPath)
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
	authpost.RegisterAuthenticationAndPostServer(grpcServer, service)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("server stopped: %v", err)
	}

}
