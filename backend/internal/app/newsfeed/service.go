package newsfeed

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/utils"
	pb_nf "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed"
	"go.uber.org/zap"
)

const (
	// Default number of items to fetch
	DefaultFeedItemCount = 10

	// Cache expiration time (24 hours by default)
	DefaultCacheExpirationTime = 24 * time.Hour

	// Default pagination values
	DefaultPageSize = 10
	MaxPageSize     = 50
)

type NewsfeedService struct {
	pb_nf.UnimplementedNewsfeedServer
	redisPool *utils.RedisPool
	logger    *zap.Logger
}

func NewNewsfeedService(cfg *configs.NewsfeedConfig) (*NewsfeedService, error) {
	// Establish logger
	logger, err := utils.NewLogger(&cfg.Logger)
	if err != nil {
		return nil, err
	}

	// Connect to enhanced Redis pool
	redisPool, err := utils.NewRedisPool(&cfg.Redis, logger)
	if err != nil {
		logger.Error("Failed to create Redis connection pool", zap.Error(err))
		return nil, errors.New("redis connection pool creation failed")
	}

	logger.Info("Successfully initialized enhanced Redis connection pool for Newsfeed service")

	return &NewsfeedService{
		redisPool: redisPool,
		logger:    logger,
	}, nil
}

// Getter methods for health checks
func (svc *NewsfeedService) GetLogger() *zap.Logger {
	return svc.logger
}

func (svc *NewsfeedService) GetRedis() *redis.Client {
	if svc.redisPool != nil {
		return svc.redisPool.Client
	}
	return nil
}

func (svc *NewsfeedService) GetRedisPool() *utils.RedisPool {
	return svc.redisPool
}

// Close gracefully closes the newsfeed service resources
func (svc *NewsfeedService) Close() error {
	if svc.redisPool != nil {
		return svc.redisPool.Close()
	}
	return nil
}

// GetNewsfeed retrieves the latest posts for a user's feed
// This implementation uses LRANGE instead of LPopCount to avoid destructive reads
func (svc *NewsfeedService) GetNewsfeed(ctx context.Context, req *pb_nf.GetNewsfeedRequest) (*pb_nf.GetNewsfeedResponse, error) {
	userID := req.GetUserId()
	page := req.GetPage()
	pageSize := req.GetPageSize()

	// Validate input
	if userID <= 0 {
		svc.logger.Warn("Invalid user ID", zap.Int64("user_id", userID))
		return &pb_nf.GetNewsfeedResponse{
			Status: pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY,
		}, nil
	}

	// Set default pagination values
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	// Calculate offset for pagination
	offset := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	// Create Redis key for the user's newsfeed
	newsfeedKey := fmt.Sprintf("newsfeed:%d", userID)

	svc.logger.Debug("Retrieving newsfeed",
		zap.Int64("user_id", userID),
		zap.String("key", newsfeedKey),
		zap.Int32("page", page),
		zap.Int32("page_size", pageSize))

	// Get total count for pagination
	totalItems, err := svc.redisPool.Client.LLen(ctx, newsfeedKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// No newsfeed found, return empty response
			svc.logger.Info("No newsfeed found for user", zap.Int64("user_id", userID))
			return &pb_nf.GetNewsfeedResponse{
				Status:      pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY,
				PostsIds:    []int64{},
				TotalPages:  0,
				CurrentPage: page,
				TotalItems:  0,
			}, nil
		}

		svc.logger.Error("Failed to get newsfeed length from Redis",
			zap.Int64("user_id", userID),
			zap.Error(err))

		return &pb_nf.GetNewsfeedResponse{
			Status: pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY,
		}, nil
	}

	// Calculate total pages
	totalPages := int32((totalItems + int64(pageSize) - 1) / int64(pageSize))

	// Get posts from Redis using LRANGE with pagination
	postIds, err := svc.redisPool.Client.LRange(ctx, newsfeedKey, offset, offset+limit-1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// No newsfeed found, return empty response
			svc.logger.Info("No posts found for user page",
				zap.Int64("user_id", userID),
				zap.Int32("page", page))
			return &pb_nf.GetNewsfeedResponse{
				Status:      pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY,
				PostsIds:    []int64{},
				TotalPages:  totalPages,
				CurrentPage: page,
				TotalItems:  int32(totalItems),
			}, nil
		}

		svc.logger.Error("Failed to retrieve newsfeed from Redis",
			zap.Int64("user_id", userID),
			zap.Error(err))

		return &pb_nf.GetNewsfeedResponse{
			Status: pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY,
		}, nil
	}

	// Convert string IDs to int64
	var postIdsInt64 []int64
	for _, idStr := range postIds {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			postIdsInt64 = append(postIdsInt64, id)
		} else {
			svc.logger.Warn("Invalid post ID in newsfeed",
				zap.String("post_id", idStr),
				zap.Error(err))
		}
	}

	svc.logger.Info("Retrieved newsfeed",
		zap.Int64("user_id", userID),
		zap.Int("post_count", len(postIdsInt64)),
		zap.Int32("page", page),
		zap.Int32("total_pages", totalPages))

	return &pb_nf.GetNewsfeedResponse{
		Status:      pb_nf.GetNewsfeedResponse_OK,
		PostsIds:    postIdsInt64,
		TotalPages:  totalPages,
		CurrentPage: page,
		TotalItems:  int32(totalItems),
	}, nil
}

func (svc *NewsfeedService) RemovePostFromNewsfeed(ctx context.Context, postID int64) error {
	// Get all newsfeed keys
	pattern := "newsfeed:*"
	var cursor uint64
	var err error

	// Remove the post ID from all newsfeeds
	postIDStr := strconv.FormatInt(postID, 10)

	for {
		var keys []string
		keys, cursor, err = svc.redisPool.Client.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			svc.logger.Error("Error scanning Redis keys", zap.Error(err))
			return err
		}

		// Process each key
		for _, key := range keys {
			// Remove post ID from this feed
			_, err := svc.redisPool.Client.LRem(ctx, key, 0, postIDStr).Result()
			if err != nil {
				svc.logger.Error("Error removing post from newsfeed",
					zap.String("key", key),
					zap.Int64("post_id", postID),
					zap.Error(err))
			}
		}

		// Exit when we've processed all keys
		if cursor == 0 {
			break
		}
	}

	return nil
}
