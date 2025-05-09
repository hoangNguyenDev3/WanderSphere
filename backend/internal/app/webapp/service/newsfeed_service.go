package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"

	pb_nf "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed"
)

func (svc *WebService) GetNewsfeed(ctx *gin.Context) {
	// Check authorization
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Call GetNewsfeed service
	resp, err := svc.NewsfeedClient.GetNewsfeed(ctx, &pb_nf.GetNewsfeedRequest{UserId: int64(userId)})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "newsfeed empty"})
		return
	} else if resp.GetStatus() == pb_nf.GetNewsfeedResponse_OK {
		ctx.IndentedJSON(http.StatusOK, resp.GetPostsIds())
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}
