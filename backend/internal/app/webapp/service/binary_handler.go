package service

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UploadBinary handles binary file uploads
func (ws *WebService) UploadBinary(c *gin.Context) {
	// Get the uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		ws.Logger.Error("Failed to get uploaded file", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Read file data
	data := make([]byte, header.Size)
	_, err = file.Read(data)
	if err != nil {
		ws.Logger.Error("Failed to read uploaded file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Get content type from header
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to storage
	result, err := ws.BinaryStorage.UploadBinary(data, header.Filename, contentType)
	if err != nil {
		ws.Logger.Error("Failed to upload binary to storage",
			zap.String("filename", header.Filename),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	ws.Logger.Info("Binary uploaded successfully",
		zap.String("key", result.Key),
		zap.String("filename", header.Filename),
		zap.Int64("size", result.Size))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"key":          result.Key,
			"url":          result.URL,
			"size":         result.Size,
			"content_type": result.ContentType,
			"uploaded_at":  result.UploadedAt,
		},
	})
}

// DownloadBinary handles binary file downloads
func (ws *WebService) DownloadBinary(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File key is required"})
		return
	}

	// Decode the key in case it's URL encoded
	decodedKey, err := url.QueryUnescape(key)
	if err != nil {
		ws.Logger.Error("Failed to decode file key", zap.String("key", key), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file key"})
		return
	}

	// Download from storage
	data, err := ws.BinaryStorage.DownloadBinary(decodedKey)
	if err != nil {
		ws.Logger.Error("Failed to download binary from storage",
			zap.String("key", decodedKey),
			zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Get file info for content type
	info, err := ws.BinaryStorage.GetBinaryInfo(decodedKey)
	if err != nil {
		ws.Logger.Warn("Failed to get binary info, using default content type",
			zap.String("key", decodedKey),
			zap.Error(err))
	}

	// Set headers
	if info != nil {
		c.Header("Content-Type", info.ContentType)
		c.Header("Content-Length", strconv.FormatInt(info.Size, 10))
	} else {
		c.Header("Content-Type", "application/octet-stream")
	}

	// Return the file data
	c.Data(http.StatusOK, c.GetHeader("Content-Type"), data)
}

// GetBinaryInfo handles getting binary file information
func (ws *WebService) GetBinaryInfo(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File key is required"})
		return
	}

	// Decode the key in case it's URL encoded
	decodedKey, err := url.QueryUnescape(key)
	if err != nil {
		ws.Logger.Error("Failed to decode file key", zap.String("key", key), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file key"})
		return
	}

	// Get file info
	info, err := ws.BinaryStorage.GetBinaryInfo(decodedKey)
	if err != nil {
		ws.Logger.Error("Failed to get binary info",
			zap.String("key", decodedKey),
			zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"key":           info.Key,
			"size":          info.Size,
			"content_type":  info.ContentType,
			"last_modified": info.LastModified,
			"etag":          info.ETag,
		},
	})
}

// GenerateBinaryDownloadURL handles generating presigned download URLs
func (ws *WebService) GenerateBinaryDownloadURL(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File key is required"})
		return
	}

	// Get expiration from query parameter (default: 1 hour)
	expirationStr := c.DefaultQuery("expiration", "3600") // 1 hour in seconds
	expirationSeconds, err := strconv.Atoi(expirationStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiration value"})
		return
	}

	expiration := time.Duration(expirationSeconds) * time.Second

	// Decode the key in case it's URL encoded
	decodedKey, err := url.QueryUnescape(key)
	if err != nil {
		ws.Logger.Error("Failed to decode file key", zap.String("key", key), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file key"})
		return
	}

	// Generate presigned URL
	downloadURL, err := ws.BinaryStorage.GenerateDownloadURL(decodedKey, expiration)
	if err != nil {
		ws.Logger.Error("Failed to generate download URL",
			zap.String("key", decodedKey),
			zap.Duration("expiration", expiration),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate download URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"download_url": downloadURL,
			"expires_in":   expirationSeconds,
		},
	})
}

// ListBinaries handles listing binary files
func (ws *WebService) ListBinaries(c *gin.Context) {
	// Get query parameters
	prefix := c.DefaultQuery("prefix", "")
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	// List files
	keys, err := ws.BinaryStorage.ListBinaries(prefix, limit)
	if err != nil {
		ws.Logger.Error("Failed to list binaries",
			zap.String("prefix", prefix),
			zap.Int("limit", limit),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"files":  keys,
			"count":  len(keys),
			"prefix": prefix,
			"limit":  limit,
		},
	})
}

// DeleteBinary handles deleting binary files
func (ws *WebService) DeleteBinary(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File key is required"})
		return
	}

	// Decode the key in case it's URL encoded
	decodedKey, err := url.QueryUnescape(key)
	if err != nil {
		ws.Logger.Error("Failed to decode file key", zap.String("key", key), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file key"})
		return
	}

	// Delete from storage
	err = ws.BinaryStorage.DeleteBinary(decodedKey)
	if err != nil {
		ws.Logger.Error("Failed to delete binary",
			zap.String("key", decodedKey),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	ws.Logger.Info("Binary deleted successfully", zap.String("key", decodedKey))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("File %s deleted successfully", decodedKey),
	})
}
