package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
)

// AddPostRouter adds post-related routes to input router
func AddPostRouter(r *gin.RouterGroup, svc *service.WebService) {
	postRouter := r.Group("posts")

	// Public routes
	postRouter.GET(":post_id", svc.GetPostDetailInfo)

	// Protected routes that require authentication
	authRouter := postRouter.Group("")
	authRouter.Use(svc.AuthRequired())
	authRouter.POST("", svc.CreatePost)
	authRouter.PUT(":post_id", svc.EditPost)
	authRouter.DELETE(":post_id", svc.DeletePost)
	authRouter.POST(":post_id/comments", svc.CommentPost)
	authRouter.POST(":post_id/likes", svc.LikePost)
}
