package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
)

func AddBinaryRouter(r *gin.RouterGroup, webService *service.WebService) {
	binaryGroup := r.Group("/binaries")

	// Protected routes (authentication required)
	authGroup := binaryGroup.Group("")
	authGroup.Use(webService.AuthRequired())
	{
		// Upload binary file
		authGroup.POST("/upload", webService.UploadBinary)

		// Generate download URL
		authGroup.GET("/:key/download-url", webService.GenerateBinaryDownloadURL)

		// List binary files
		authGroup.GET("/", webService.ListBinaries)

		// Delete binary file
		authGroup.DELETE("/:key", webService.DeleteBinary)
	}

	// Public routes (no authentication required)
	{
		// Get binary file info
		binaryGroup.GET("/:key/info", webService.GetBinaryInfo)

		// Download binary file (single parameter)
		binaryGroup.GET("/:key", webService.DownloadBinary)
	}
}
