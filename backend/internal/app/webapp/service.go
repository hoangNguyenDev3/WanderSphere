package webapp

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/hoangNguyenDev3/WanderSphere/internal/app/webapp/v1"
)

func Run() {
	// Create server using Gin
	router := gin.Default()

	v1Router := router.Group("/v1")
	v1.AddUserRouter(v1Router)

	router.Run(":8080")
}
