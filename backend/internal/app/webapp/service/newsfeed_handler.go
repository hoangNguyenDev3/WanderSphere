package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/hoangNguyenDev3/WanderSphere/backend/docs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
	pb_newsfeed "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed"
	"go.uber.org/zap"
)

// GetNewsfeed godoc
// @Summary Get user's newsfeed
// @Description Get the current user's newsfeed with cursor-based pagination
// @Tags newsfeed
// @Accept json
// @Produce json
// @Param cursor query string false "Pagination cursor (empty for first page)"
// @Param limit query int false "Items per page (default 10, max 50)"
// @Success 200 {object} types.NewsfeedResponse "User's newsfeed"
// @Failure 400 {object} types.MessageResponse "Validation error"
// @Failure 401 {object} types.MessageResponse "Unauthorized"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /newsfeed [get]
// @Security ApiKeyAuth
func (svc *WebService) GetNewsfeed(ctx *gin.Context) {
	// Check authorization
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Parse cursor query param
	cursor := ctx.Query("cursor")

	// Parse limit query param with defaults
	limit := int32(10)
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if parsed, parseErr := strconv.Atoi(limitStr); parseErr == nil && parsed > 0 {
			limit = int32(parsed)
			if limit > 50 {
				limit = 50
			}
		}
	}

	// Parse optional page query param for list-based fallback pagination
	page := int32(1)
	if pageStr := ctx.Query("page"); pageStr != "" {
		if parsed, err := strconv.Atoi(pageStr); err == nil && parsed > 0 {
			page = int32(parsed)
		}
	}

	// Call GetNewsfeed service
	resp, err := svc.NewsfeedClient.GetNewsfeed(ctx, &pb_newsfeed.GetNewsfeedRequest{
		UserId:   int64(userId),
		Page:     page,
		PageSize: limit,
		Cursor:   cursor,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_newsfeed.GetNewsfeedResponse_NEWSFEED_EMPTY {
		ctx.JSON(http.StatusOK, types.NewsfeedResponse{
			PostsIds: []int64{}, // Return empty array
			HasMore:  false,
		})
		return
	} else if resp.GetStatus() == pb_newsfeed.GetNewsfeedResponse_OK {
		ctx.JSON(http.StatusOK, types.NewsfeedResponse{
			PostsIds:   resp.GetPostsIds(),
			NextCursor: resp.GetNextCursor(),
			HasMore:    resp.GetHasMore(),
		})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// RemovePostFromNewsfeed calls the newsfeed service to remove a post from all newsfeeds
// This should be called when a post is deleted
func (svc *WebService) RemovePostFromNewsfeed(postID int64) {
	// This is just a placeholder method showing how you'd call the RemovePostFromNewsfeed method
	// from the auth service when it needs to invalidate the cache after a post deletion
	// Since we haven't exposed this as a gRPC method yet, this is just for illustration
	svc.Logger.Info("Post deleted, would call newsfeed service to remove post from all feeds",
		zap.Int64("post_id", postID))
}
