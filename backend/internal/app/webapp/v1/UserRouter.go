package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
)

func AddUserRouter(r *gin.RouterGroup, svc *service.WebService) {
	userRouter := r.Group("users")
	userRouter.POST("signup", svc.CreateUser)
	userRouter.POST("login", svc.CheckUserAuthentication)
	userRouter.POST("edit", svc.EditUser)
}
