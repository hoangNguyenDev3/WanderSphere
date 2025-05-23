package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
)

// AddNewsfeedRouter adds newsfeed-related routes to input router
func AddNewsfeedRouter(r *gin.RouterGroup, svc *service.WebService) {
	newsfeedRouter := r.Group("newsfeed")

	// Protected routes that require authentication
	authRouter := newsfeedRouter.Group("")
	authRouter.Use(svc.AuthRequired())
	authRouter.GET("", svc.GetNewsfeed)
}
