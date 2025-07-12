package newsfeed

import (
	"context"
	"encoding/base64"
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

	// Sorted set key prefix for engagement-ranked newsfeeds
	newsfeedRankedKeyPrefix = "newsfeed_ranked:"
)

// Cursor format: base64("{postID}") — encodes only the last seen post ID.
// Rank-based pagination: we look up the rank of the cursor post and start after it.

func encodeCursor(postID int64) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(postID, 10)))
}

func decodeCursor(cursor string) (postID int64, err error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, fmt.Errorf("invalid cursor encoding: %w", err)
	}
	postID, err = strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid cursor postID: %w", err)
	}
	return postID, nil
}

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
// Supports cursor-based pagination for ranked sorted set feeds,
// and offset-based pagination for legacy list-based feeds.
func (svc *NewsfeedService) GetNewsfeed(ctx context.Context, req *pb_nf.GetNewsfeedRequest) (*pb_nf.GetNewsfeedResponse, error) {
	userID := req.GetUserId()
	pageSize := req.GetPageSize()
	cursor := req.GetCursor()

	// Validate input
	if userID <= 0 {
		svc.logger.Warn("Invalid user ID", zap.Int64("user_id", userID))
		return &pb_nf.GetNewsfeedResponse{
			Status: pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY,
		}, nil
	}

	// Set default pagination values
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	// Create Redis keys for the user's newsfeed
	newsfeedKey := fmt.Sprintf("newsfeed:%d", userID)
	rankedKey := fmt.Sprintf("%s%d", newsfeedRankedKeyPrefix, userID)

	svc.logger.Debug("Retrieving newsfeed",
		zap.Int64("user_id", userID),
		zap.String("list_key", newsfeedKey),
		zap.String("ranked_key", rankedKey),
		zap.Int32("page_size", pageSize),
		zap.String("cursor", cursor))

	// Try ranked sorted set first, fall back to list-based feed
	rankedExists, err := svc.redisPool.Client.Exists(ctx, rankedKey).Result()
	if err != nil {
		svc.logger.Warn("Failed to check ranked feed existence, falling back to list",
			zap.Int64("user_id", userID),
			zap.Error(err))
		rankedExists = 0
	}

	if rankedExists > 0 {
		return svc.getNewsfeedFromRankedSetWithCursor(ctx, rankedKey, userID, pageSize, cursor)
	}

	// Fall back to list-based feed with offset pagination (page=1 for first page)
	page := req.GetPage()
	if page <= 0 {
		page = 1
	}
	offset := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	svc.logger.Debug("Ranked feed not found, using list-based feed",
		zap.Int64("user_id", userID))
	return svc.getNewsfeedFromList(ctx, newsfeedKey, userID, page, pageSize, offset, limit)
}

// getNewsfeedFromRankedSetWithCursor reads the feed from the engagement-ranked sorted set
// using rank-based pagination with ZRevRangeWithScores.
func (svc *NewsfeedService) getNewsfeedFromRankedSetWithCursor(
	ctx context.Context, rankedKey string, userID int64,
	pageSize int32, cursor string,
) (*pb_nf.GetNewsfeedResponse, error) {

	var start int64 = 0

	if cursor != "" {
		// Decode cursor to get last seen postID
		cursorPostID, err := decodeCursor(cursor)
		if err != nil {
			svc.logger.Warn("Invalid cursor, starting from beginning",
				zap.Int64("user_id", userID),
				zap.Error(err))
			// Fall through with start = 0 (first page)
		} else {
			// Find the rank of the cursor post in the sorted set (descending order)
			rank, err := svc.redisPool.Client.ZRevRank(ctx, rankedKey, strconv.FormatInt(cursorPostID, 10)).Result()
			if err != nil {
				// Cursor post may have been deleted; start from beginning
				svc.logger.Warn("Cursor post not found in ranked set, starting from beginning",
					zap.Int64("user_id", userID),
					zap.Int64("cursor_post_id", cursorPostID),
					zap.Error(err))
				// Fall through with start = 0
			} else {
				start = rank + 1 // Start after the cursor post
			}
		}
	}

	// Fetch pageSize + 1 to determine has_more
	stop := start + int64(pageSize) // This fetches pageSize+1 items (ZRevRange is inclusive on both ends)

	posts, err := svc.redisPool.Client.ZRevRangeWithScores(ctx, rankedKey, start, stop).Result()
	if err != nil {
		svc.logger.Warn("Failed to retrieve ranked newsfeed, falling back to list",
			zap.Int64("user_id", userID),
			zap.Error(err))
		newsfeedKey := fmt.Sprintf("newsfeed:%d", userID)
		return svc.getNewsfeedFromList(ctx, newsfeedKey, userID, 1, pageSize, 0, int64(pageSize))
	}

	hasMore := len(posts) > int(pageSize)
	if hasMore {
		posts = posts[:pageSize] // Trim to pageSize
	}

	if len(posts) == 0 {
		return &pb_nf.GetNewsfeedResponse{
			Status:  pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY,
			HasMore: false,
		}, nil
	}

	// Convert to post IDs and build next cursor
	postIDs := make([]int64, 0, len(posts))
	var lastPostID int64
	for _, z := range posts {
		memberStr, ok := z.Member.(string)
		if !ok {
			continue
		}
		postID, err := strconv.ParseInt(memberStr, 10, 64)
		if err != nil {
			svc.logger.Error("Failed to parse post ID from ranked set member",
				zap.String("member", memberStr),
				zap.Error(err))
			continue
		}
		postIDs = append(postIDs, postID)
		lastPostID = postID
	}

	// Get total count for metadata
	totalItems, _ := svc.redisPool.Client.ZCard(ctx, rankedKey).Result()

	var nextCursor string
	if hasMore {
		nextCursor = encodeCursor(lastPostID)
	}

	return &pb_nf.GetNewsfeedResponse{
		Status:     pb_nf.GetNewsfeedResponse_OK,
		PostsIds:   postIDs,
		TotalItems: int32(totalItems),
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// getNewsfeedFromList reads the feed from the legacy chronological list
func (svc *NewsfeedService) getNewsfeedFromList(ctx context.Context, newsfeedKey string, userID int64, page, pageSize int32, offset, limit int64) (*pb_nf.GetNewsfeedResponse, error) {
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

	hasMore := page < totalPages
	return &pb_nf.GetNewsfeedResponse{
		Status:      pb_nf.GetNewsfeedResponse_OK,
		PostsIds:    postIdsInt64,
		TotalPages:  totalPages,
		CurrentPage: page,
		TotalItems:  int32(totalItems),
		HasMore:     hasMore,
	}, nil
}

func (svc *NewsfeedService) RemovePostFromNewsfeed(ctx context.Context, postID int64) error {
	// Remove from both list-based and ranked feeds
	listPattern := "newsfeed:*"
	rankedPattern := newsfeedRankedKeyPrefix + "*"
	var cursor uint64
	var err error

	// Remove the post ID from all newsfeeds
	postIDStr := strconv.FormatInt(postID, 10)

	// Remove from list-based feeds
	for {
		var keys []string
		keys, cursor, err = svc.redisPool.Client.Scan(ctx, cursor, listPattern, 10).Result()
		if err != nil {
			svc.logger.Error("Error scanning Redis list keys", zap.Error(err))
			return err
		}

		for _, key := range keys {
			_, err := svc.redisPool.Client.LRem(ctx, key, 0, postIDStr).Result()
			if err != nil {
				svc.logger.Error("Error removing post from list newsfeed",
					zap.String("key", key),
					zap.Int64("post_id", postID),
					zap.Error(err))
			}
		}

		if cursor == 0 {
			break
		}
	}

	// Remove from ranked sorted set feeds
	cursor = 0
	for {
		var keys []string
		keys, cursor, err = svc.redisPool.Client.Scan(ctx, cursor, rankedPattern, 10).Result()
		if err != nil {
			svc.logger.Error("Error scanning Redis ranked keys", zap.Error(err))
			return err
		}

		for _, key := range keys {
			_, err := svc.redisPool.Client.ZRem(ctx, key, postIDStr).Result()
			if err != nil {
				svc.logger.Error("Error removing post from ranked newsfeed",
					zap.String("key", key),
					zap.Int64("post_id", postID),
					zap.Error(err))
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
