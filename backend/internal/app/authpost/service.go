package authpost

import (
	"errors"
	"path/filepath"
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
	db               *gorm.DB
	migrationManager *utils.MigrationManager
	nfPubClient      client_nfp.Client
	logger           *zap.Logger
}

func NewAuthenticateAndPostService(cfg *configs.AuthenticateAndPostConfig) (*AuthenticateAndPostService, error) {
	// Create logger first for better error reporting
	logger, err := utils.NewLogger(&cfg.Logger)
	if err != nil {
		// Fall back to production logger if there's an error
		logger, _ = zap.NewProduction()
	}

	logger.Info("Initializing AuthenticateAndPostService with database migrations")

	// Connect to database
	postgresConfig := postgres.Config{
		DSN: cfg.Postgres.DSN,
	}
	db, err := gorm.Open(postgres.New(postgresConfig), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get underlying sql.DB", zap.Error(err))
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Initialize migration manager
	migrationManager := utils.NewMigrationManager(db, &cfg.Postgres, logger)

	// Run automatic database migrations
	migrationDir := filepath.Join("migrations")
	logger.Info("Running database migrations", zap.String("migration_dir", migrationDir))

	// Temporarily skip migrations since they're already applied
	logger.Info("Skipping migrations - they are already applied")
	/*
		if err := migrationManager.Migrate(migrationDir); err != nil {
			logger.Error("Database migration failed", zap.Error(err))
			return nil, err
		}
	*/

	// Get migration status for logging
	status, err := migrationManager.GetStatus(migrationDir)
	if err != nil {
		logger.Warn("Failed to get migration status", zap.Error(err))
	} else {
		logger.Info("Database migration status",
			zap.Int("total_migrations", status.TotalMigrations),
			zap.Int("applied_migrations", status.AppliedMigrations),
			zap.Strings("pending_migrations", status.PendingMigrations),
			zap.String("last_applied", status.LastApplied))
	}

	// Connect to NewsfeedPublishingClient if configured
	var nfPubClient client_nfp.Client
	if len(cfg.NewsfeedPublishing.Hosts) > 0 {
		nfPubClient, err = client_nfp.NewClient(cfg.NewsfeedPublishing.Hosts)
		if err != nil {
			logger.Error("Failed to connect to newsfeed publishing service", zap.Error(err))
			// Continue without newsfeed publishing client
		} else {
			logger.Info("Successfully connected to newsfeed publishing service")
		}
	}

	logger.Info("AuthenticateAndPostService initialized successfully")
	return &AuthenticateAndPostService{
		db:               db,
		migrationManager: migrationManager,
		nfPubClient:      nfPubClient,
		logger:           logger,
	}, nil
}

// Getter methods for health checks
func (a *AuthenticateAndPostService) GetDB() *gorm.DB {
	return a.db
}

func (a *AuthenticateAndPostService) GetLogger() *zap.Logger {
	return a.logger
}

func (a *AuthenticateAndPostService) GetMigrationManager() *utils.MigrationManager {
	return a.migrationManager
}

func (a *AuthenticateAndPostService) GetRedis() interface{} {
	// AuthPost service doesn't directly use Redis, return nil
	return nil
}

// GetMigrationStatus returns the current status of database migrations
func (a *AuthenticateAndPostService) GetMigrationStatus() (*utils.MigrationStatus, error) {
	return a.migrationManager.GetStatus("migrations")
}

// Close gracefully closes the AuthPost service resources
func (a *AuthenticateAndPostService) Close() error {
	if sqlDB, err := a.db.DB(); err == nil {
		return sqlDB.Close()
	}
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
