package service

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/utils"
	"go.uber.org/zap"

	client_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/client/authpost"
	client_nf "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/client/newsfeed"
	pb_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	pb_nf "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed"
)

var validate = types.NewValidator()

type WebService struct {
	AuthenticateAndPostClient pb_aap.AuthenticateAndPostClient
	NewsfeedClient            pb_nf.NewsfeedClient
	RedisClient               *redis.Client
	Logger                    *zap.Logger
	Config                    *configs.WebConfig
}

func NewWebService(cfg *configs.WebConfig) (*WebService, error) {
	aapClient, err := client_aap.NewClient(cfg.AuthenticateAndPost.Hosts)
	if err != nil {
		return nil, err
	}

	nfClient, err := client_nf.NewClient(cfg.Newsfeed.Hosts)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password})
	if redisClient == nil {
		return nil, errors.New("redis connection failed")
	}

	logger, err := utils.NewLogger(&cfg.Logger)
	if err != nil {
		return nil, err
	}

	return &WebService{
		AuthenticateAndPostClient: aapClient,
		NewsfeedClient:            nfClient,
		RedisClient:               redisClient,
		Logger:                    logger,
		Config:                    cfg,
	}, nil
}

// Getter methods for health checks
func (ws *WebService) GetLogger() *zap.Logger {
	return ws.Logger
}

func (ws *WebService) GetRedis() *redis.Client {
	return ws.RedisClient
}

func (ws *WebService) GetConfig() *configs.WebConfig {
	return ws.Config
}
