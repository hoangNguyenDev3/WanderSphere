package v1

import "github.com/gin-gonic/gin"

func AddUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("/users")

	userRouter.GET("", func(ctx *gin.Context) {})
	userRouter.POST("", func(ctx *gin.Context) {})
	userRouter.PUT("", func(ctx *gin.Context) {})

}
