package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/newsfeed"
	pb_nf "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed"
	"google.golang.org/grpc"
)

func main() {
	// Flags
	cfgPath := flag.String("conf", "configs/files/local_newsfeed.yml", "Path to config file for this service")
	flag.Parse()

	// Load configurations
	cfg, err := configs.GetNewsfeedConfigDirect(*cfgPath)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	// Start new newsfeed service
	service, err := newsfeed.NewNewsfeedService(cfg)
	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		log.Fatalf("can not listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb_nf.RegisterNewsfeedServer(grpcServer, service)

	log.Printf("Starting Newsfeed service on port %d", cfg.Port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
