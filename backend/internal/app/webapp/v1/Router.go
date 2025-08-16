package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
)

func AddAllRouter(r *gin.RouterGroup, webService *service.WebService) {
	// Apply RefreshSession middleware first to set user_id in context if authenticated
	r.Use(webService.RefreshSession())
	// Apply rate limiting middleware second so it can use user_id for per-user limits
	r.Use(webService.RateLimit())

	// Add all the routers
	AddUserRouter(r, webService)
	AddFriendRouter(r, webService)
	AddPostRouter(r, webService)
	AddNewsfeedRouter(r, webService)
	AddBinaryRouter(r, webService)
}
