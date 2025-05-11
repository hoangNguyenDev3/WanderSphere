package newsfeed_publishing

import (
	"context"
	"log"
	"math/rand"

	pb_nfp "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed_publishing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client defines the interface for the Newsfeed Publishing client
type Client interface {
	PublishPost(ctx context.Context, in *pb_nfp.PublishPostRequest) (*pb_nfp.PublishPostResponse, error)
}

// NewClient creates a new client for the Newsfeed Publishing service
func NewClient(hosts []string) (Client, error) {
	var opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	clients := make([]pb_nfp.NewsfeedPublishingClient, 0, len(hosts))
	for _, host := range hosts {
		conn, err := grpc.Dial(host, opts...)
		if err != nil {
			log.Printf("Failed to dial host %s: %v", host, err)
			continue
		}

		client := pb_nfp.NewNewsfeedPublishingClient(conn)
		clients = append(clients, client)
	}

	if len(clients) == 0 {
		return nil, log.Output(1, "no available newsfeed_publishing service hosts")
	}

	return &randomClient{clients: clients}, nil
}

type randomClient struct {
	clients []pb_nfp.NewsfeedPublishingClient
}

// PublishPost forwards to a random client
func (rc *randomClient) PublishPost(ctx context.Context, in *pb_nfp.PublishPostRequest) (*pb_nfp.PublishPostResponse, error) {
	return rc.clients[rand.Intn(len(rc.clients))].PublishPost(ctx, in)
}
