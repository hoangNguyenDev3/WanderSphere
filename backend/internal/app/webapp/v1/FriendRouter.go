package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
)

// AddFriendRouter adds friend-related routes to input router
func AddFriendRouter(r *gin.RouterGroup, svc *service.WebService) {
	friendRouter := r.Group("friends")

	// Public routes
	friendRouter.GET(":user_id/followers", svc.GetUserFollowers)
	friendRouter.GET(":user_id/followings", svc.GetUserFollowings)
	friendRouter.GET(":user_id/posts", svc.GetUserPosts)

	// Protected routes that require authentication
	authRouter := friendRouter.Group("")
	authRouter.Use(svc.AuthRequired())
	authRouter.POST(":user_id", svc.FollowUser)
	authRouter.DELETE(":user_id", svc.UnfollowUser)
}
