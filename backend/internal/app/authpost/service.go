package authpost

import (
	"errors"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewAuthenticateAndPostService(cfg *configs.AuthenticateAndPostConfig) (*AuthenticateAndPostService, error) {
	// Connect to database
	postgresConfig := postgres.Config{
		DSN: cfg.Postgres.DSN,
	}
	db, err := gorm.Open(postgres.New(postgresConfig), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Connect to NewsfeedPublishingClient
	nfPubClient, err := client_nfp.NewClient(cfg.NewsfeedPublishing.Hosts)
	if err != nil {
		return nil, err
	}

	// Establish logger
	logger, err := utils.NewLogger(&cfg.Logger)
	if err != nil {
		return nil, err
	}

	return &AuthenticateAndPostService{
		db:          db,
		nfPubClient: nfPubClient,
		logger:      logger,
	}, nil
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
