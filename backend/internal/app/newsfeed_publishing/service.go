package newsfeed_publishing_svc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	client_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/client/authpost"
	pb_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	pb_nfp "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed_publishing"
	"github.com/segmentio/kafka-go"
)

type NewsfeedPublishingService struct {
	pb_nfp.UnimplementedNewsfeedPublishingServer
	kafkaWriter               *kafka.Writer
	kafkaReader               *kafka.Reader
	redisClient               *redis.Client
	authenticateAndPostClient pb_aap.AuthenticateAndPostClient
}

func NewNewsfeedPublishingService(cfg *configs.NewsfeedPublishingConfig) (*NewsfeedPublishingService, error) {
	// Connect to kafka writer
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
		Logger:  log.New(os.Stdout, "kafka writer: ", 0),
	})
	if kafkaWriter == nil {
		return nil, errors.New("failed connecting to kafka writer")
	}

	// Connect to kafka reader
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
		Logger:  log.New(os.Stdout, "kafka reader: ", 0),
	})
	if kafkaReader == nil {
		return nil, errors.New("kafka connection failed")
	}

	// Connect to redis
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password})
	if redisClient == nil {
		return nil, errors.New("redis connection failed")
	}

	// Connect to aap service
	aapClient, err := client_aap.NewClient(cfg.AuthenticateAndPost.Hosts)
	if err != nil {
		return nil, err
	}

	// Return
	return &NewsfeedPublishingService{
		kafkaWriter:               kafkaWriter,
		kafkaReader:               kafkaReader,
		redisClient:               redisClient,
		authenticateAndPostClient: aapClient,
	}, nil
}

func (svc *NewsfeedPublishingService) PublishPost(ctx context.Context, info *pb_nfp.PublishPostRequest) (*pb_nfp.PublishPostResponse, error) {
	value := map[string]int64{
		"user_id": info.GetUserId(),
		"post_id": info.GetPostId(),
	}
	jsonValue, _ := json.Marshal(value)
	err := svc.kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte("post"),
		Value: jsonValue,
		Headers: []kafka.Header{
			{Key: "Content-Type", Value: []byte("application/json")},
		},
	})
	if err != nil {
		return nil, err
	}

	return &pb_nfp.PublishPostResponse{
		Status: pb_nfp.PublishPostResponse_OK,
	}, nil
}

func (svc *NewsfeedPublishingService) Run() {
	for {
		message, err := svc.kafkaReader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("worker will sleep 0.1s then try again")
			time.Sleep(100 * time.Millisecond)
			continue
		}
		svc.processMessage(message)
	}
}

func (svc *NewsfeedPublishingService) processMessage(message kafka.Message) {
	msgType := string(message.Key)

	// Process message based on its key
	if msgType == "post" {
		svc.processPost(message.Value)
	}
}

func (svc *NewsfeedPublishingService) processPost(value []byte) {
	var message map[string]int64
	err := json.Unmarshal(value, &message)
	if err != nil {
		panic(err)
	}

	// 2. Find followers of user that created post
	followersKey := "followers:" + strconv.Itoa(int(message["user_id"]))
	numKey, _ := svc.redisClient.Exists(context.Background(), followersKey).Result()
	if numKey == 0 {
		resp, err := svc.authenticateAndPostClient.GetUserFollower(
			context.Background(),
			&pb_aap.GetUserFollowerRequest{
				UserId: message["user_id"],
			})
		if err != nil {
			panic(err)
		}

		followersIDs := resp.GetFollowersIds()
		if len(followersIDs) > 0 {
			_, err = svc.redisClient.RPush(context.Background(), followersKey, followersIDs).Result()
			if err != nil {
				panic(err)
			}
		}
	}
	followersIDs, err := svc.redisClient.LRange(context.Background(), followersKey, 0, -1).Result()
	if err != nil {
		panic("err")
	}

	// 3. Add this post_id into followers' newsfeed
	for _, id := range followersIDs {
		newsfeedKey := "newsfeed:" + id
		_, err := svc.redisClient.RPush(context.Background(), newsfeedKey, message["post_id"]).Result()
		if err != nil {
			panic(err)
		}
	}
}
