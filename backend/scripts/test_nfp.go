package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb_nfp "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed_publishing"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:19004", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb_nfp.NewNewsfeedPublishingClient(conn)

	// Set up timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Send the request
	resp, err := client.PublishPost(ctx, &pb_nfp.PublishPostRequest{
		UserId: 123,
		PostId: 456,
	})

	if err != nil {
		log.Fatalf("Failed to publish post: %v", err)
	}

	fmt.Printf("Response status: %s\n", resp.Status.String())
}
