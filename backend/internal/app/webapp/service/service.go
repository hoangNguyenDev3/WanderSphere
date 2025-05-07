package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
	client_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/client/authpost"

	pb_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	pb_nf "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed"
)

var validate = types.NewValidator()

type WebService struct {
	AuthenticateAndPostClient pb_aap.AuthenticateAndPostClient
	NewsfeedClient            pb_nf.NewsfeedClient
	RedisClient               *redis.Client
}

func NewWebService(conf *configs.WebConfig) (*WebService, error) {
	aapClient, err := client_aap.NewClient(conf.AuthenticateAndPost.Hosts)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{Addr: conf.Redis.Addr, Password: conf.Redis.Password})
	if redisClient == nil {
		return nil, errors.New("redis connection failed")
	}

	return &WebService{
		AuthenticateAndPostClient: aapClient,
		RedisClient:               redisClient,
	}, nil
}

func (svc *WebService) checkSessionAuthentication(ctx *gin.Context) (sessionId string, userId int, userName string, err error) {
	sessionId, err = ctx.Cookie("session_id")
	if err != nil {
		return "", 0, "", err
	}

	userId, err = svc.RedisClient.HGet(svc.RedisClient.Context(), sessionId, "userId").Int()
	if err != nil {
		return "", 0, "", err
	}

	userName, err = svc.RedisClient.HGet(svc.RedisClient.Context(), sessionId, "userName").Result()
	if err != nil {
		return "", 0, "", err
	}

	return sessionId, userId, userName, nil
}
