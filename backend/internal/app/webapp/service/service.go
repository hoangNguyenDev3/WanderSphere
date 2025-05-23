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
	RedisPool                 *utils.RedisPool
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

	logger, err := utils.NewLogger(&cfg.Logger)
	if err != nil {
		return nil, err
	}

	// Use enhanced Redis connection pool
	redisPool, err := utils.NewRedisPool(&cfg.Redis, logger)
	if err != nil {
		logger.Error("Failed to create Redis connection pool", zap.Error(err))
		return nil, errors.New("redis connection pool creation failed")
	}

	logger.Info("Successfully initialized enhanced Redis connection pool for Web service")

	return &WebService{
		AuthenticateAndPostClient: aapClient,
		NewsfeedClient:            nfClient,
		RedisPool:                 redisPool,
		Logger:                    logger,
		Config:                    cfg,
	}, nil
}

// Getter methods for health checks
func (ws *WebService) GetLogger() *zap.Logger {
	return ws.Logger
}

func (ws *WebService) GetRedis() *redis.Client {
	if ws.RedisPool != nil {
		return ws.RedisPool.Client
	}
	return nil
}

func (ws *WebService) GetRedisPool() *utils.RedisPool {
	return ws.RedisPool
}

func (ws *WebService) GetConfig() *configs.WebConfig {
	return ws.Config
}

// Close gracefully closes the web service resources
func (ws *WebService) Close() error {
	if ws.RedisPool != nil {
		return ws.RedisPool.Close()
	}
	return nil
}
