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
	redisClient *redis.Client
	logger      *zap.Logger
}

func NewNewsfeedService(cfg *configs.NewsfeedConfig) (*NewsfeedService, error) {
	// Connect to redisClient
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password})
	if redisClient == nil {
		return nil, errors.New("redis connection failed")
	}

	// Establish logger
	logger, err := utils.NewLogger(&cfg.Logger)
	if err != nil {
		return nil, err
	}

	return &NewsfeedService{
		redisClient: redisClient,
		logger:      logger,
	}, nil
}

// GetNewsfeed retrieves the latest posts for a user's feed
// This implementation uses LRANGE instead of LPopCount to avoid destructive reads
func (svc *NewsfeedService) GetNewsfeed(ctx context.Context, request *pb_nf.GetNewsfeedRequest) (*pb_nf.GetNewsfeedResponse, error) {
	// Get the newsfeed key for this user
	newsfeedKey := fmt.Sprintf("newsfeed:%d", request.GetUserId())

	// Check if newsfeed exists
	exists, err := svc.redisClient.Exists(ctx, newsfeedKey).Result()
	if err != nil {
		svc.logger.Error("Error checking newsfeed existence", zap.Error(err))
		return nil, err
	}

	// Return empty response if newsfeed doesn't exist
	if exists == 0 {
		return &pb_nf.GetNewsfeedResponse{
			Status: pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY,
		}, nil
	}

	// Use LRANGE (non-destructive) instead of LPopCount (destructive)
	// This keeps the posts in the feed after retrieval
	postsIds, err := svc.redisClient.LRange(ctx, newsfeedKey, 0, DefaultFeedItemCount-1).Result()
	if err != nil {
		svc.logger.Error("Error retrieving newsfeed", zap.Error(err))
		return nil, err
	}

	// Set TTL on the feed to ensure cache freshness
	svc.redisClient.Expire(ctx, newsfeedKey, DefaultCacheExpirationTime)

	// If no posts found, return empty response
	if len(postsIds) == 0 {
		return &pb_nf.GetNewsfeedResponse{
			Status: pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY,
		}, nil
	}

	// Convert string IDs to int64
	var int64PostsIds []int64
	for _, id := range postsIds {
		intPostId, err := strconv.Atoi(id)
		if err != nil {
			svc.logger.Warn("Invalid post ID in cache", zap.String("post_id", id))
			continue
		}
		int64PostsIds = append(int64PostsIds, int64(intPostId))
	}

	// Return empty response if no valid post IDs
	if len(int64PostsIds) == 0 {
		return &pb_nf.GetNewsfeedResponse{
			Status: pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY,
		}, nil
	}

	return &pb_nf.GetNewsfeedResponse{
		Status:   pb_nf.GetNewsfeedResponse_OK,
		PostsIds: int64PostsIds,
	}, nil
}

// RemovePostFromNewsfeed removes a post from all newsfeeds when it's deleted
// This method should be called when a post is deleted to maintain cache consistency
func (svc *NewsfeedService) RemovePostFromNewsfeed(ctx context.Context, postID int64) error {
	// Get all newsfeed keys
	pattern := "newsfeed:*"
	var cursor uint64
	var err error

	// Remove the post ID from all newsfeeds
	postIDStr := strconv.FormatInt(postID, 10)

	for {
		var keys []string
		keys, cursor, err = svc.redisClient.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			svc.logger.Error("Error scanning Redis keys", zap.Error(err))
			return err
		}

		// Process each key
		for _, key := range keys {
			// Remove post ID from this feed
			_, err := svc.redisClient.LRem(ctx, key, 0, postIDStr).Result()
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
