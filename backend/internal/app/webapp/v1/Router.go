package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
)

func AddAllRouter(r *gin.RouterGroup, webService *service.WebService) {
	// Apply RefreshSession middleware to all routes to handle session extension
	r.Use(webService.RefreshSession())

	// Add all the routers
	AddUserRouter(r, webService)
	AddFriendRouter(r, webService)
	AddPostRouter(r, webService)
	AddNewsfeedRouter(r, webService)
}
