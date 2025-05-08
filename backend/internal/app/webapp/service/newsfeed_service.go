package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"

	pb_nf "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed"
)

func (svc *WebService) GetNewsfeed(ctx *gin.Context) {
	// Check authorization
	_, userId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Call GetNewsfeed service
	newsfeed, err := svc.NewsfeedClient.GetNewsfeed(ctx, &pb_nf.NewsfeedRequest{UserId: int64(userId)})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}

	// Return
	var posts []gin.H
	for _, postDetailInfo := range newsfeed.Posts {
		posts = append(posts, svc.newMapPost(postDetailInfo))
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"newsfeed": posts})
}
