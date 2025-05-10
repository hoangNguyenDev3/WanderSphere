package newsfeed

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/utils"
	pb_nf "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed"
	"go.uber.org/zap"
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

func (svc *NewsfeedService) GetNewsfeed(ctx context.Context, request *pb_nf.GetNewsfeedRequest) (*pb_nf.GetNewsfeedResponse, error) {
	// Query newsfeed from redis
	newsfeedKey := "newsfeed:" + fmt.Sprint(request.GetUserId())
	postsIds, err := svc.redisClient.LPopCount(svc.redisClient.Context(), newsfeedKey, 5).Result()
	if errors.Is(err, redis.Nil) {
		return &pb_nf.GetNewsfeedResponse{
			Status: pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY,
		}, nil
	} else if err != nil {
		return nil, err
	}

	var int64PostsIds []int64
	for _, id := range postsIds {
		intPostId, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		int64PostsIds = append(int64PostsIds, int64(intPostId))
	}
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
