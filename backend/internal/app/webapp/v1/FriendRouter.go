package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
)

// AddFriendRouter adds friend-related routes to input router
func AddFriendRouter(r *gin.RouterGroup, svc *service.WebService) {
	friendRouter := r.Group("friends")
	friendRouter.GET(":user_id", svc.GetUserFollower)
	friendRouter.POST(":user_id", svc.FollowUser)
	friendRouter.DELETE(":user_id", svc.UnfollowUser)
	friendRouter.GET(":user_id/posts", svc.GetUserPost)

}
