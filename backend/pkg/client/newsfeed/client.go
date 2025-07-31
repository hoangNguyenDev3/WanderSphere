package newsfeed

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/middleware"
	pb_nf "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed"
)

func NewClient(hosts []string) (pb_nf.NewsfeedClient, error) {
	logger, _ := zap.NewProduction()
	clients := make([]pb_nf.NewsfeedClient, 0, len(hosts))
	for _, host := range hosts {
		cb := middleware.NewCircuitBreaker(
			fmt.Sprintf("newsfeed-%s", host),
			middleware.DefaultCircuitBreakerConfig(),
			logger,
		)
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(middleware.UnaryClientInterceptor(cb)),
		}
		conn, err := grpc.Dial(host, opts...)
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
			return nil, err
		}

		client := pb_nf.NewNewsfeedClient(conn)
		clients = append(clients, client)
	}

	return &randomClient{clients: clients}, nil
}

type randomClient struct {
	clients []pb_nf.NewsfeedClient
}

func (a *randomClient) GetNewsfeed(ctx context.Context, in *pb_nf.GetNewsfeedRequest, opts ...grpc.CallOption) (*pb_nf.GetNewsfeedResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetNewsfeed(ctx, in, opts...)
}

func (a *randomClient) InvalidateCache(ctx context.Context, in *pb_nf.InvalidateCacheRequest, opts ...grpc.CallOption) (*pb_nf.InvalidateCacheResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].InvalidateCache(ctx, in, opts...)
}
