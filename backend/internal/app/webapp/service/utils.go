package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
)

// sendErrorResponse standardizes error response format across the API
func sendErrorResponse(c *gin.Context, statusCode int, errorType, message string) {
	c.JSON(statusCode, types.ErrorResponse{
		Error:   errorType,
		Message: message,
		Code:    statusCode,
	})
}

// sendSuccessResponse standardizes success response format
func sendSuccessResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, types.MessageResponse{
		Message: message,
		Status:  "success",
	})
}

// sendDataResponse sends a response with data
func sendDataResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}
