package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
)

func AddUserRouter(r *gin.RouterGroup, svc *service.WebService) {
	userRouter := r.Group("users")

	// Public routes
	userRouter.POST("signup", svc.CreateUser)
	userRouter.POST("login", svc.CheckUserAuthentication)
	userRouter.GET(":user_id", svc.GetUserDetailInfo)

	// Protected routes that require authentication
	authRouter := userRouter.Group("")
	authRouter.Use(svc.AuthRequired())
	authRouter.POST("edit", svc.EditUser)
}
