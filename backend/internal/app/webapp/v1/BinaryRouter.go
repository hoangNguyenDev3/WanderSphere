package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
)

func AddBinaryRouter(r *gin.RouterGroup, webService *service.WebService) {
	binaryGroup := r.Group("/binaries")
	{
		// Upload binary file
		binaryGroup.POST("/upload", webService.UploadBinary)

		// Download binary file
		binaryGroup.GET("/:key", webService.DownloadBinary)

		// Get binary file info
		binaryGroup.GET("/:key/info", webService.GetBinaryInfo)

		// Generate download URL
		binaryGroup.GET("/:key/download-url", webService.GenerateBinaryDownloadURL)

		// List binary files
		binaryGroup.GET("/", webService.ListBinaries)

		// Delete binary file
		binaryGroup.DELETE("/:key", webService.DeleteBinary)
	}
}
