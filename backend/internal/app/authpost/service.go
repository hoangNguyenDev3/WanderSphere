package authpost

import (
	"errors"
	"time"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/utils"
	client_nfp "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/client/newsfeed_publishing"
	pb "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// AuthenticateAndPostService implements the AuthenticateAndPost service
type AuthenticateAndPostService struct {
	pb.UnimplementedAuthenticateAndPostServer
	db          *gorm.DB
	nfPubClient client_nfp.Client
	logger      *zap.Logger
}

func NewAuthenticateAndPostService(cfg *configs.AuthenticateAndPostConfig) (*AuthenticateAndPostService, error) {
	// Connect to database
	postgresConfig := postgres.Config{
		DSN: cfg.Postgres.DSN,
	}
	db, err := gorm.Open(postgres.New(postgresConfig), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Create logger
	var logger *zap.Logger
	logger, err = utils.NewLogger(&cfg.Logger)
	if err != nil {
		// Fall back to production logger if there's an error
		logger, _ = zap.NewProduction()
	}

	// Connect to NewsfeedPublishingClient if configured
	var nfPubClient client_nfp.Client
	if len(cfg.NewsfeedPublishing.Hosts) > 0 {
		nfPubClient, err = client_nfp.NewClient(cfg.NewsfeedPublishing.Hosts)
		if err != nil {
			logger.Error("Failed to connect to newsfeed publishing service", zap.Error(err))
			// Continue without newsfeed publishing client
		}
	}

	logger.Info("AuthenticateAndPostService initialized")
	return &AuthenticateAndPostService{
		db:          db,
		nfPubClient: nfPubClient,
		logger:      logger,
	}, nil
}

// Getter methods for health checks
func (a *AuthenticateAndPostService) GetDB() *gorm.DB {
	return a.db
}

func (a *AuthenticateAndPostService) GetLogger() *zap.Logger {
	return a.logger
}

func (a *AuthenticateAndPostService) GetRedis() interface{} {
	// AuthPost service doesn't directly use Redis, return nil
	return nil
}

// findUserById checks if an user with provided userId exists in database
func (a *AuthenticateAndPostService) findUserById(userId int64) (exist bool, user types.User) {
	result := a.db.First(&user, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, types.User{}
	}
	return true, user
}

// findUserByUserName checks if an user with provided username exists in database
func (a *AuthenticateAndPostService) findUserByUserName(userName string) (exist bool, user types.User) {
	result := a.db.Where(&types.User{UserName: userName}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, types.User{}
	}
	return true, user
}

// findPostById checks if an user with provided userId exists in database
func (a *AuthenticateAndPostService) findPostById(postId int64) (exist bool, post types.Post) {
	result := a.db.First(&post, postId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, types.Post{}
	}
	return true, post
}
