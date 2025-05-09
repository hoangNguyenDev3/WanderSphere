package authpost

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
	pb_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthenticateAndPostService struct {
	pb_aap.UnimplementedAuthenticateAndPostServer
	db          *gorm.DB
	kafkaWriter *kafka.Writer

	logger *zap.Logger
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

	// Connect to kafka
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
		Logger:  log.New(os.Stdout, "kafka writer: ", 0),
	})
	if kafkaWriter == nil {
		return nil, errors.New("failed connecting to kafka writer")
	}

	// Establish logger
	logger, err := newLogger()
	if err != nil {
		return nil, err
	}

	return &AuthenticateAndPostService{
		db:          db,
		kafkaWriter: kafkaWriter,
		logger:      logger,
	}, nil
}

func newLogger() (*zap.Logger, error) {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "/tmp/logs"],
		"errorOutputPaths": ["stderr"],
		"initialFields": {"foo": "bar"},
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return nil, err
	}
	logger := zap.Must(cfg.Build())
	return logger, nil
}

func (a *AuthenticateAndPostService) findUserById(userId int64) (exist bool, user types.User) {
	result := a.db.First(&user, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, types.User{}
	}
	return true, user
}

func (a *AuthenticateAndPostService) findUserByUserName(userName string) (exist bool, user types.User) {
	result := a.db.Where(&types.User{UserName: userName}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, types.User{}
	}
	return true, user
}

func (a *AuthenticateAndPostService) findPostById(postId int64) (exist bool, post types.Post) {
	result := a.db.First(&post, postId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, types.Post{}
	}
	return true, post
}
