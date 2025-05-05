package main

import (
	"fmt"
	"log"
	"net"

	authpost "github.com/hoangNguyenDev3/backend/internal/app/authpost"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 1080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	service := authpost.NewAuthenticationAndPostService()
	grpcServer := grpc.NewServer(opts...)
	authpost.RegisterAuthenticationAndPostServer(grpcServer, service)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
