package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
)

// AddNewsfeedRouter adds newsfeed-related routes to input router
func AddNewsfeedRouter(r *gin.RouterGroup, svc *service.WebService) {
	feedRouter := r.Group("newsfeed")

	// All newsfeed routes require authentication
	feedRouter.Use(svc.AuthRequired())
	feedRouter.GET("", svc.GetNewsfeed)
}
