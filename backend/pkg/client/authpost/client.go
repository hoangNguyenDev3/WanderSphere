package authpost

import (
	"context"
	"log"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/hoangNguyenDev3/WanderSphere/pkg/types/proto/pb/authpost"
)

type RandomClient struct {
	clients []pb.AuthenticationAndPostClient
}

func NewClient(hosts []string) (pb.AuthenticationAndPostClient, error) {
	var opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	clients := make([]pb.AuthenticationAndPostClient, len(hosts))
	for _, host := range hosts {
		conn, err := grpc.Dial(host, opts...)
		if err != nil {
			log.Fatalf("failed to connect to %s: %v", host, err)
			return nil, err
		}
		client := pb.NewAuthenticationAndPostClient(conn)
		clients = append(clients, client)
	}
	return &RandomClient{clients: clients}, nil
}

func (c *RandomClient) CheckUserAuthentication(ctx context.Context, in *pb.UserInfo, opts ...grpc.CallOption) (*pb.UserResult, error) {
	return c.clients[rand.Intn(len(c.clients))].CheckUserAuthentication(ctx, in, opts...)
}

func (c *RandomClient) CreateUser(ctx context.Context, in *pb.UserDetailInfo, opts ...grpc.CallOption) (*pb.UserResult, error) {
	return c.clients[rand.Intn(len(c.clients))].CreateUser(ctx, in, opts...)
}

func (c *RandomClient) EditUser(ctx context.Context, in *pb.UserDetailInfo, opts ...grpc.CallOption) (*pb.UserResult, error) {
	return c.clients[rand.Intn(len(c.clients))].EditUser(ctx, in, opts...)
}

func (c *RandomClient) GetUserFollowers(ctx context.Context, in *pb.UserInfo, opts ...grpc.CallOption) (*pb.UserFollower, error) {
	return c.clients[rand.Intn(len(c.clients))].GetUserFollowers(ctx, in, opts...)
}

func (c *RandomClient) GetPostDetail(ctx context.Context, in *pb.GetPostRequest, opts ...grpc.CallOption) (*pb.Post, error) {
	return c.clients[rand.Intn(len(c.clients))].GetPostDetail(ctx, in, opts...)
}
