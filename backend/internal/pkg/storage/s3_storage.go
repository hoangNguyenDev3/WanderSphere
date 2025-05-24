package storage

import (
	"fmt"
	"io"
	"time"

	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/utils"
	"go.uber.org/zap"
)

// BinaryStorage interface defines methods for binary storage operations
type BinaryStorage interface {
	UploadBinary(data []byte, filename string, contentType string) (*UploadResult, error)
	DownloadBinary(key string) ([]byte, error)
	DeleteBinary(key string) error
	GetBinaryInfo(key string) (*BinaryInfo, error)
	GenerateDownloadURL(key string, expiration time.Duration) (string, error)
	ListBinaries(prefix string, limit int) ([]string, error)
}

// UploadResult represents the result of a binary upload
type UploadResult struct {
	Key         string
	URL         string
	Size        int64
	ContentType string
	UploadedAt  time.Time
}

// BinaryInfo represents information about a stored binary
type BinaryInfo struct {
	Key          string
	Size         int64
	ContentType  string
	LastModified time.Time
	ETag         string
}

// S3BinaryStorage implements BinaryStorage using S3
type S3BinaryStorage struct {
	s3Service *utils.S3Service
	logger    *zap.Logger
}

// NewS3BinaryStorage creates a new S3 binary storage service
func NewS3BinaryStorage(config utils.S3Config, logger *zap.Logger) (BinaryStorage, error) {
	s3Service, err := utils.NewS3Service(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 service: %w", err)
	}

	return &S3BinaryStorage{
		s3Service: s3Service,
		logger:    logger,
	}, nil
}

// UploadBinary uploads binary data to S3 storage
func (s *S3BinaryStorage) UploadBinary(data []byte, filename string, contentType string) (*UploadResult, error) {
	s.logger.Info("Uploading binary to S3",
		zap.String("filename", filename),
		zap.String("contentType", contentType),
		zap.Int("size", len(data)))

	result, err := s.s3Service.UploadBinary(data, filename, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload binary: %w", err)
	}

	return &UploadResult{
		Key:         result.Key,
		URL:         result.Location,
		Size:        result.Size,
		ContentType: contentType,
		UploadedAt:  time.Now(),
	}, nil
}

// DownloadBinary downloads binary data from S3 storage
func (s *S3BinaryStorage) DownloadBinary(key string) ([]byte, error) {
	s.logger.Info("Downloading binary from S3", zap.String("key", key))

	data, err := s.s3Service.DownloadFile(key)
	if err != nil {
		return nil, fmt.Errorf("failed to download binary: %w", err)
	}

	return data, nil
}

// DeleteBinary deletes binary data from S3 storage
func (s *S3BinaryStorage) DeleteBinary(key string) error {
	s.logger.Info("Deleting binary from S3", zap.String("key", key))

	err := s.s3Service.DeleteFile(key)
	if err != nil {
		return fmt.Errorf("failed to delete binary: %w", err)
	}

	return nil
}

// GetBinaryInfo gets information about a binary file
func (s *S3BinaryStorage) GetBinaryInfo(key string) (*BinaryInfo, error) {
	s.logger.Info("Getting binary info from S3", zap.String("key", key))

	result, err := s.s3Service.GetFileInfo(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get binary info: %w", err)
	}

	return &BinaryInfo{
		Key:          key,
		Size:         *result.ContentLength,
		ContentType:  *result.ContentType,
		LastModified: *result.LastModified,
		ETag:         *result.ETag,
	}, nil
}

// GenerateDownloadURL generates a presigned URL for downloading a binary
func (s *S3BinaryStorage) GenerateDownloadURL(key string, expiration time.Duration) (string, error) {
	s.logger.Info("Generating download URL for binary",
		zap.String("key", key),
		zap.Duration("expiration", expiration))

	url, err := s.s3Service.GeneratePresignedURL(key, expiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate download URL: %w", err)
	}

	return url, nil
}

// ListBinaries lists binary files with optional prefix
func (s *S3BinaryStorage) ListBinaries(prefix string, limit int) ([]string, error) {
	s.logger.Info("Listing binaries from S3",
		zap.String("prefix", prefix),
		zap.Int("limit", limit))

	if limit <= 0 {
		limit = 100 // Default limit
	}

	keys, err := s.s3Service.ListFiles(prefix, int64(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to list binaries: %w", err)
	}

	return keys, nil
}

// UploadStream uploads a binary stream to S3 storage
func (s *S3BinaryStorage) UploadStream(reader io.Reader, filename string, contentType string, size int64) (*UploadResult, error) {
	s.logger.Info("Uploading binary stream to S3",
		zap.String("filename", filename),
		zap.String("contentType", contentType),
		zap.Int64("size", size))

	result, err := s.s3Service.UploadFile(utils.UploadFileRequest{
		Data:        reader,
		Key:         "", // Let the service generate a unique key
		ContentType: contentType,
		Size:        size,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload binary stream: %w", err)
	}

	return &UploadResult{
		Key:         result.Key,
		URL:         result.Location,
		Size:        result.Size,
		ContentType: contentType,
		UploadedAt:  time.Now(),
	}, nil
}
