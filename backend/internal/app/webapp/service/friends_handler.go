package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
	pb_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	"go.uber.org/zap"
)

// FollowUser godoc
// @Summary Follow user
// @Description Follow another user
// @Tags friends
// @Accept json
// @Produce json
// @Param user_id path int true "User ID to follow"
// @Success 200 {object} types.MessageResponse "User followed successfully"
// @Failure 400 {object} types.MessageResponse "Validation error or user not found"
// @Failure 401 {object} types.MessageResponse "Unauthorized"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /friends/{user_id} [post]
// @Security ApiKeyAuth
func (svc *WebService) FollowUser(ctx *gin.Context) {
	// Check authorization
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Check URL params
	followingId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	svc.Logger.Info("Web FollowUser request",
		zap.Int("user_id", userId),
		zap.Int("following_id", followingId))

	// Call FollowUser service
	resp, err := svc.AuthenticateAndPostClient.FollowUser(ctx, &pb_aap.FollowUserRequest{
		UserId:      int64(userId),
		FollowingId: int64(followingId),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.FollowUserResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.FollowUserResponse_ALREADY_FOLLOWED {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "already following this user"})
		return
	} else if resp.GetStatus() == pb_aap.FollowUserResponse_OK {
		ctx.JSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// UnfollowUser godoc
// @Summary Unfollow user
// @Description Unfollow another user
// @Tags friends
// @Accept json
// @Produce json
// @Param user_id path int true "User ID to unfollow"
// @Success 200 {object} types.MessageResponse "User unfollowed successfully"
// @Failure 400 {object} types.MessageResponse "Validation error or user not found"
// @Failure 401 {object} types.MessageResponse "Unauthorized"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /friends/{user_id} [delete]
// @Security ApiKeyAuth
func (svc *WebService) UnfollowUser(ctx *gin.Context) {
	// Check authorization
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Check URL params
	followingId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	// Call UnfollowUser service
	resp, err := svc.AuthenticateAndPostClient.UnfollowUser(ctx, &pb_aap.UnfollowUserRequest{
		UserId:      int64(userId),
		FollowingId: int64(followingId),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.UnfollowUserResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.UnfollowUserResponse_NOT_FOLLOWED {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "not following this user"})
		return
	} else if resp.GetStatus() == pb_aap.UnfollowUserResponse_OK {
		ctx.JSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// GetUserFollowers godoc
// @Summary Get user followers
// @Description Get the followers of a user
// @Tags friends
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} types.UserFollowerResponse "User's followers"
// @Failure 400 {object} types.MessageResponse "Validation error or user not found"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /friends/{user_id}/followers [get]
func (svc *WebService) GetUserFollowers(ctx *gin.Context) {
	// Check URL params
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	// Call GetUserFollower service
	resp, err := svc.AuthenticateAndPostClient.GetUserFollower(ctx, &pb_aap.GetUserFollowerRequest{
		UserId: int64(userId),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.GetUserFollowerResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.GetUserFollowerResponse_OK {
		ctx.JSON(http.StatusOK, types.UserFollowerResponse{
			FollowersIds: resp.GetFollowersIds(),
		})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// GetUserFollowings godoc
// @Summary Get user followings
// @Description Get the users followed by a user
// @Tags friends
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} types.UserFollowingResponse "User's followings"
// @Failure 400 {object} types.MessageResponse "Validation error or user not found"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /friends/{user_id}/followings [get]
func (svc *WebService) GetUserFollowings(ctx *gin.Context) {
	// Check URL params
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	// Call GetUserFollowing service
	resp, err := svc.AuthenticateAndPostClient.GetUserFollowing(ctx, &pb_aap.GetUserFollowingRequest{
		UserId: int64(userId),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.GetUserFollowingResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.GetUserFollowingResponse_OK {
		ctx.JSON(http.StatusOK, types.UserFollowingResponse{
			FollowingsIds: resp.GetFollowingsIds(),
		})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// GetUserPosts godoc
// @Summary Get user posts
// @Description Get all posts of a user
// @Tags friends
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} types.UserPostsResponse "User's posts"
// @Failure 400 {object} types.MessageResponse "Validation error or user not found"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /friends/{user_id}/posts [get]
func (svc *WebService) GetUserPosts(ctx *gin.Context) {
	// Check URL params
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	// Call GetUserPosts service
	resp, err := svc.AuthenticateAndPostClient.GetUserPosts(ctx, &pb_aap.GetUserPostsRequest{
		UserId: int64(userId),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.GetUserPostsResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.GetUserPostsResponse_OK {
		ctx.JSON(http.StatusOK, types.UserPostsResponse{
			PostsIds: resp.GetPostsIds(),
		})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}
