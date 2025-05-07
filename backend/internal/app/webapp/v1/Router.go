package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
)

func AddAllRouter(r *gin.RouterGroup, webService *service.WebService) {
	AddUserRouter(r, webService)
	AddFriendRouter(r, webService)
}
