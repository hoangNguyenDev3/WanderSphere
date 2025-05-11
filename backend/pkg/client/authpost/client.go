package authpost

import (
	"context"
	"log"

	"math/rand"

	pb_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient creates a new client for the AuthenticateAndPost service
func NewClient(hosts []string) (pb_aap.AuthenticateAndPostClient, error) {
	var opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	clients := make([]pb_aap.AuthenticateAndPostClient, 0, len(hosts))
	for _, host := range hosts {
		conn, err := grpc.Dial(host, opts...)
		if err != nil {
			log.Printf("Failed to dial host %s: %v", host, err)
			continue
		}

		client := pb_aap.NewAuthenticateAndPostClient(conn)
		clients = append(clients, client)
	}

	if len(clients) == 0 {
		return nil, log.Output(1, "no available authpost service hosts")
	}

	return &randomClient{clients: clients}, nil
}

type randomClient struct {
	clients []pb_aap.AuthenticateAndPostClient
}

// Group: Users
func (a *randomClient) CheckUserAuthentication(ctx context.Context, in *pb_aap.CheckUserAuthenticationRequest, opts ...grpc.CallOption) (*pb_aap.CheckUserAuthenticationResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].CheckUserAuthentication(ctx, in, opts...)
}

func (a *randomClient) CreateUser(ctx context.Context, in *pb_aap.CreateUserRequest, opts ...grpc.CallOption) (*pb_aap.CreateUserResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].CreateUser(ctx, in, opts...)
}

func (a *randomClient) EditUser(ctx context.Context, in *pb_aap.EditUserRequest, opts ...grpc.CallOption) (*pb_aap.EditUserResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].EditUser(ctx, in, opts...)
}

func (a *randomClient) GetUserDetailInfo(ctx context.Context, in *pb_aap.GetUserDetailInfoRequest, opts ...grpc.CallOption) (*pb_aap.GetUserDetailInfoResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetUserDetailInfo(ctx, in, opts...)
}

// Group: Friends

func (a *randomClient) GetUserFollower(ctx context.Context, in *pb_aap.GetUserFollowerRequest, opts ...grpc.CallOption) (*pb_aap.GetUserFollowerResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetUserFollower(ctx, in, opts...)
}

func (a *randomClient) GetUserFollowing(ctx context.Context, in *pb_aap.GetUserFollowingRequest, opts ...grpc.CallOption) (*pb_aap.GetUserFollowingResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetUserFollowing(ctx, in, opts...)
}

func (a *randomClient) FollowUser(ctx context.Context, in *pb_aap.FollowUserRequest, opts ...grpc.CallOption) (*pb_aap.FollowUserResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].FollowUser(ctx, in, opts...)
}

func (a *randomClient) UnfollowUser(ctx context.Context, in *pb_aap.UnfollowUserRequest, opts ...grpc.CallOption) (*pb_aap.UnfollowUserResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].UnfollowUser(ctx, in, opts...)
}

func (a *randomClient) GetUserPosts(ctx context.Context, in *pb_aap.GetUserPostsRequest, opts ...grpc.CallOption) (*pb_aap.GetUserPostsResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetUserPosts(ctx, in, opts...)
}

// Group: Posts

func (a *randomClient) CreatePost(ctx context.Context, in *pb_aap.CreatePostRequest, opts ...grpc.CallOption) (*pb_aap.CreatePostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].CreatePost(ctx, in, opts...)
}

func (a *randomClient) GetPostDetailInfo(ctx context.Context, in *pb_aap.GetPostDetailInfoRequest, opts ...grpc.CallOption) (*pb_aap.GetPostDetailInfoResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetPostDetailInfo(ctx, in, opts...)
}

func (a *randomClient) EditPost(ctx context.Context, in *pb_aap.EditPostRequest, opts ...grpc.CallOption) (*pb_aap.EditPostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].EditPost(ctx, in, opts...)
}

func (a *randomClient) DeletePost(ctx context.Context, in *pb_aap.DeletePostRequest, opts ...grpc.CallOption) (*pb_aap.DeletePostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].DeletePost(ctx, in, opts...)
}

func (a *randomClient) CommentPost(ctx context.Context, in *pb_aap.CommentPostRequest, opts ...grpc.CallOption) (*pb_aap.CommentPostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].CommentPost(ctx, in, opts...)
}

func (a *randomClient) LikePost(ctx context.Context, in *pb_aap.LikePostRequest, opts ...grpc.CallOption) (*pb_aap.LikePostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].LikePost(ctx, in, opts...)
}
