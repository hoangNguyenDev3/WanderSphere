package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
)

// AddPostRouter adds post-related routes to input router
func AddPostRouter(r *gin.RouterGroup, svc *service.WebService) {
	postRouter := r.Group("posts")

	// Public routes
	postRouter.GET(":post_id", svc.GetPostDetail)

	// Protected routes that require authentication
	authRouter := postRouter.Group("")
	authRouter.Use(svc.AuthRequired())
	authRouter.POST("", svc.CreatePost)
	authRouter.PUT(":post_id", svc.EditPost)
	authRouter.DELETE(":post_id", svc.DeletePost)
	authRouter.POST(":post_id", svc.CommentPost)
	authRouter.POST(":post_id/likes", svc.LikePost)
	authRouter.GET("url", svc.GetS3PresignedUrl)
}
