package authpost

import (
	"context"

	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
	pb_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	"go.uber.org/zap"
)

func (a *AuthenticateAndPostService) GetUserFollower(ctx context.Context, info *pb_aap.GetUserFollowerRequest) (*pb_aap.GetUserFollowerResponse, error) {
	exist, _ := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.GetUserFollowerResponse{
			Status: pb_aap.GetUserFollowerResponse_USER_NOT_FOUND,
		}, nil
	}

	var user types.User
	result := a.db.Preload("Followers").First(&user, info.GetUserId())
	if result.Error != nil {
		return nil, result.Error
	}

	var followersIds []int64
	for _, follower := range user.Followers {
		followersIds = append(followersIds, int64(follower.ID))
	}
	return &pb_aap.GetUserFollowerResponse{
		Status:       pb_aap.GetUserFollowerResponse_OK,
		FollowersIds: followersIds,
	}, nil
}

func (a *AuthenticateAndPostService) GetUserFollowing(ctx context.Context, info *pb_aap.GetUserFollowingRequest) (*pb_aap.GetUserFollowingResponse, error) {
	exist, _ := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.GetUserFollowingResponse{
			Status: pb_aap.GetUserFollowingResponse_USER_NOT_FOUND,
		}, nil
	}

	var user types.User
	result := a.db.Preload("Followings").First(&user, info.GetUserId())
	if result.Error != nil {
		return nil, result.Error
	}

	var followingsIds []int64
	for _, following := range user.Followings {
		followingsIds = append(followingsIds, int64(following.ID))
	}
	return &pb_aap.GetUserFollowingResponse{
		Status:        pb_aap.GetUserFollowingResponse_OK,
		FollowingsIds: followingsIds,
	}, nil
}

func (a *AuthenticateAndPostService) FollowUser(ctx context.Context, info *pb_aap.FollowUserRequest) (*pb_aap.FollowUserResponse, error) {
	a.logger.Info("FollowUser request received",
		zap.Int64("user_id", info.GetUserId()),
		zap.Int64("following_id", info.GetFollowingId()))

	// Check if user is trying to follow themselves
	if info.GetUserId() == info.GetFollowingId() {
		a.logger.Info("Self-follow attempt detected, preventing",
			zap.Int64("user_id", info.GetUserId()))
		return &pb_aap.FollowUserResponse{Status: pb_aap.FollowUserResponse_ALREADY_FOLLOWED}, nil
	}

	// Check if the user exists
	exist, user := a.findUserById(info.GetUserId())
	if !exist {
		a.logger.Warn("User not found", zap.Int64("user_id", info.GetUserId()))
		return &pb_aap.FollowUserResponse{Status: pb_aap.FollowUserResponse_USER_NOT_FOUND}, nil
	}
	exist, friend := a.findUserById(info.GetFollowingId())
	if !exist {
		a.logger.Warn("Friend to follow not found", zap.Int64("following_id", info.GetFollowingId()))
		return &pb_aap.FollowUserResponse{Status: pb_aap.FollowUserResponse_USER_NOT_FOUND}, nil
	}

	// Load user with current followings to check for existing relationship
	result := a.db.Preload("Followings").First(&user, info.GetUserId())
	if result.Error != nil {
		a.logger.Error("Failed to load user with followings",
			zap.Int64("user_id", info.GetUserId()),
			zap.Error(result.Error))
		return nil, result.Error
	}

	// Check if already following
	for _, following := range user.Followings {
		if int64(following.ID) == info.GetFollowingId() {
			a.logger.Info("Already following user",
				zap.Int64("user_id", info.GetUserId()),
				zap.Int64("following_id", info.GetFollowingId()))
			return &pb_aap.FollowUserResponse{Status: pb_aap.FollowUserResponse_ALREADY_FOLLOWED}, nil
		}
	}

	// Add the following relationship
	err := a.db.Model(&user).Association("Followings").Append(&friend)
	if err != nil {
		a.logger.Error("Failed to create follow relationship",
			zap.Int64("user_id", info.GetUserId()),
			zap.Int64("following_id", info.GetFollowingId()),
			zap.Error(err))
		return nil, err
	}

	a.logger.Info("Successfully created follow relationship",
		zap.Int64("user_id", info.GetUserId()),
		zap.Int64("following_id", info.GetFollowingId()))
	return &pb_aap.FollowUserResponse{
		Status: pb_aap.FollowUserResponse_OK,
	}, nil
}

func (a *AuthenticateAndPostService) UnfollowUser(ctx context.Context, info *pb_aap.UnfollowUserRequest) (*pb_aap.UnfollowUserResponse, error) {
	a.logger.Info("UnfollowUser request received",
		zap.Int64("user_id", info.GetUserId()),
		zap.Int64("following_id", info.GetFollowingId()))

	// Check if users exist
	exist, user := a.findUserById(info.GetUserId())
	if !exist {
		a.logger.Warn("User not found", zap.Int64("user_id", info.GetUserId()))
		return &pb_aap.UnfollowUserResponse{Status: pb_aap.UnfollowUserResponse_USER_NOT_FOUND}, nil
	}
	exist, friend := a.findUserById(info.GetFollowingId())
	if !exist {
		a.logger.Warn("Friend to unfollow not found", zap.Int64("following_id", info.GetFollowingId()))
		return &pb_aap.UnfollowUserResponse{Status: pb_aap.UnfollowUserResponse_USER_NOT_FOUND}, nil
	}

	// Load user with current followings
	result := a.db.Preload("Followings").First(&user, info.GetUserId())
	if result.Error != nil {
		a.logger.Error("Failed to load user with followings",
			zap.Int64("user_id", info.GetUserId()),
			zap.Error(result.Error))
		return nil, result.Error
	}

	// Check if currently following
	currentlyFollowing := false
	for _, following := range user.Followings {
		if int64(following.ID) == info.GetFollowingId() {
			currentlyFollowing = true
			break
		}
	}
	if !currentlyFollowing {
		a.logger.Info("User is not following the target user",
			zap.Int64("user_id", info.GetUserId()),
			zap.Int64("following_id", info.GetFollowingId()))
		return &pb_aap.UnfollowUserResponse{Status: pb_aap.UnfollowUserResponse_NOT_FOLLOWED}, nil
	}

	// Remove the following relationship
	err := a.db.Model(&user).Association("Followings").Delete(&friend)
	if err != nil {
		a.logger.Error("Failed to remove follow relationship",
			zap.Int64("user_id", info.GetUserId()),
			zap.Int64("following_id", info.GetFollowingId()),
			zap.Error(err))
		return nil, err
	}

	a.logger.Info("Successfully removed follow relationship",
		zap.Int64("user_id", info.GetUserId()),
		zap.Int64("following_id", info.GetFollowingId()))
	return &pb_aap.UnfollowUserResponse{
		Status: pb_aap.UnfollowUserResponse_OK,
	}, nil
}

func (a *AuthenticateAndPostService) GetUserPosts(ctx context.Context, info *pb_aap.GetUserPostsRequest) (*pb_aap.GetUserPostsResponse, error) {
	exist, _ := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.GetUserPostsResponse{Status: pb_aap.GetUserPostsResponse_USER_NOT_FOUND}, nil
	}

	var user types.User
	a.db.Preload("Posts").First(&user, info.GetUserId())

	// Return
	var posts_ids []int64
	for _, post := range user.Posts {
		posts_ids = append(posts_ids, int64(post.ID))
	}

	return &pb_aap.GetUserPostsResponse{
		Status:   pb_aap.GetUserPostsResponse_OK,
		PostsIds: posts_ids,
	}, nil
}
