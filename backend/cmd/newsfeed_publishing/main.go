package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	newsfeed_publishing_svc "github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/newsfeed_publishing"
	pb_nfp "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed_publishing"
	"google.golang.org/grpc"
)

func main() {
	// Flags
	cfgPath := flag.String("conf", "configs/files/local_newsfeed_publishing.yml", "Path to config file for this service")
	flag.Parse()

	// Load configurations
	cfg, err := configs.GetNewsfeedPublishingConfigDirect(*cfgPath)
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
