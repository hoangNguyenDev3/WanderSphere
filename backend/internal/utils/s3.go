package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// S3Config holds the configuration for S3 storage
type S3Config struct {
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	Region          string `yaml:"region"`
	Bucket          string `yaml:"bucket"`
	Endpoint        string `yaml:"endpoint"` // Optional: for S3-compatible services like MinIO
	DisableSSL      bool   `yaml:"disable_ssl"`
	ForcePathStyle  bool   `yaml:"force_path_style"`
}

// S3Service provides methods for interacting with S3 storage
type S3Service struct {
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	bucket     string
	logger     *zap.Logger
}

// NewS3Service creates a new S3 service instance
func NewS3Service(config S3Config, logger *zap.Logger) (*S3Service, error) {
	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(config.Region),
		Credentials:      credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
		Endpoint:         aws.String(config.Endpoint),
		DisableSSL:       aws.Bool(config.DisableSSL),
		S3ForcePathStyle: aws.Bool(config.ForcePathStyle),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Create S3 client
	client := s3.New(sess)

	// Create uploader and downloader
	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)

	return &S3Service{
		client:     client,
		uploader:   uploader,
		downloader: downloader,
		bucket:     config.Bucket,
		logger:     logger,
	}, nil
}

// UploadFileRequest represents a file upload request
type UploadFileRequest struct {
	Data        io.Reader
	Key         string
	ContentType string
	Size        int64
	Metadata    map[string]*string
}

// UploadFileResponse represents a file upload response
type UploadFileResponse struct {
	Key      string
	Location string
	ETag     string
	Size     int64
}

// UploadFile uploads a file to S3 and returns the file information
func (s *S3Service) UploadFile(req UploadFileRequest) (*UploadFileResponse, error) {
	// Generate unique key if not provided
	if req.Key == "" {
		req.Key = s.generateUniqueKey()
	}

	// Auto-detect content type if not provided
	if req.ContentType == "" {
		req.ContentType = s.detectContentType(req.Key)
	}

	// Prepare upload input
	uploadInput := &s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(req.Key),
		Body:        req.Data,
		ContentType: aws.String(req.ContentType),
		Metadata:    req.Metadata,
	}

	// Upload the file
	result, err := s.uploader.Upload(uploadInput)
	if err != nil {
		s.logger.Error("Failed to upload file to S3",
			zap.String("key", req.Key),
			zap.Error(err))
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	s.logger.Info("File uploaded successfully",
		zap.String("key", req.Key),
		zap.String("location", result.Location))

	return &UploadFileResponse{
		Key:      req.Key,
		Location: result.Location,
		ETag:     aws.StringValue(result.ETag),
		Size:     req.Size,
	}, nil
}

// DownloadFile downloads a file from S3
func (s *S3Service) DownloadFile(key string) ([]byte, error) {
	// Create a buffer to write the file to
	buf := aws.NewWriteAtBuffer([]byte{})

	// Download the file
	_, err := s.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		s.logger.Error("Failed to download file from S3",
			zap.String("key", key),
			zap.Error(err))
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	s.logger.Info("File downloaded successfully", zap.String("key", key))
	return buf.Bytes(), nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(key string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		s.logger.Error("Failed to delete file from S3",
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("failed to delete file: %w", err)
	}

	s.logger.Info("File deleted successfully", zap.String("key", key))
	return nil
}

// GetFileInfo gets information about a file in S3
func (s *S3Service) GetFileInfo(key string) (*s3.HeadObjectOutput, error) {
	result, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		s.logger.Error("Failed to get file info from S3",
			zap.String("key", key),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return result, nil
}

// GeneratePresignedURL generates a presigned URL for accessing a file
func (s *S3Service) GeneratePresignedURL(key string, expiration time.Duration) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		s.logger.Error("Failed to generate presigned URL",
			zap.String("key", key),
			zap.Error(err))
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url, nil
}

// GeneratePresignedPutURL generates a presigned URL for uploading a file
func (s *S3Service) GeneratePresignedPutURL(key string, contentType string, expiration time.Duration) (string, error) {
	req, _ := s.client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		s.logger.Error("Failed to generate presigned PUT URL",
			zap.String("key", key),
			zap.String("contentType", contentType),
			zap.Error(err))
		return "", fmt.Errorf("failed to generate presigned PUT URL: %w", err)
	}

	s.logger.Info("Generated presigned PUT URL",
		zap.String("key", key),
		zap.String("contentType", contentType),
		zap.Duration("expiration", expiration))

	return url, nil
}

// ListFiles lists files in the S3 bucket with optional prefix
func (s *S3Service) ListFiles(prefix string, maxKeys int64) ([]string, error) {
	if maxKeys <= 0 {
		maxKeys = 1000 // Default limit
	}

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s.bucket),
		MaxKeys: aws.Int64(maxKeys),
	}

	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	result, err := s.client.ListObjectsV2(input)
	if err != nil {
		s.logger.Error("Failed to list files from S3",
			zap.String("prefix", prefix),
			zap.Error(err))
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	var keys []string
	for _, obj := range result.Contents {
		keys = append(keys, aws.StringValue(obj.Key))
	}

	return keys, nil
}

// CopyFile copies a file within S3
func (s *S3Service) CopyFile(sourceKey, destinationKey string) error {
	copySource := fmt.Sprintf("%s/%s", s.bucket, sourceKey)

	_, err := s.client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(s.bucket),
		CopySource: aws.String(copySource),
		Key:        aws.String(destinationKey),
	})
	if err != nil {
		s.logger.Error("Failed to copy file in S3",
			zap.String("source", sourceKey),
			zap.String("destination", destinationKey),
			zap.Error(err))
		return fmt.Errorf("failed to copy file: %w", err)
	}

	s.logger.Info("File copied successfully",
		zap.String("source", sourceKey),
		zap.String("destination", destinationKey))
	return nil
}

// generateUniqueKey generates a unique key for file storage
func (s *S3Service) generateUniqueKey() string {
	return fmt.Sprintf("%d-%s", time.Now().Unix(), uuid.New().String())
}

// detectContentType detects content type based on file extension
func (s *S3Service) detectContentType(key string) string {
	ext := filepath.Ext(key)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream" // Default binary content type
	}
	return contentType
}

// UploadBinary is a convenience method for uploading binary data
func (s *S3Service) UploadBinary(data []byte, filename string, contentType string) (*UploadFileResponse, error) {
	return s.UploadFile(UploadFileRequest{
		Data:        bytes.NewReader(data),
		Key:         s.generateKeyWithFilename(filename),
		ContentType: contentType,
		Size:        int64(len(data)),
	})
}

// generateKeyWithFilename generates a key that includes the original filename
func (s *S3Service) generateKeyWithFilename(filename string) string {
	// Clean the filename to make it S3-safe
	cleanFilename := strings.ReplaceAll(filename, " ", "_")
	cleanFilename = strings.ReplaceAll(cleanFilename, "/", "_")

	timestamp := time.Now().Unix()
	uuid := uuid.New().String()

	return fmt.Sprintf("binaries/%d/%s/%s", timestamp, uuid, cleanFilename)
}

// GetFileURL returns the public URL for a file (if bucket allows public access)
func (s *S3Service) GetFileURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, "us-east-1", key)
}
