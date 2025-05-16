package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// KafkaManager handles Kafka operations including topic management
type KafkaManager struct {
	Brokers []string
	Logger  *zap.Logger
}

// TopicConfig represents Kafka topic configuration
type TopicConfig struct {
	Name              string
	NumPartitions     int
	ReplicationFactor int
	RetentionMs       int64
}

// NewKafkaManager creates a new Kafka manager
func NewKafkaManager(cfg *configs.KafkaConfig, logger *zap.Logger) *KafkaManager {
	return &KafkaManager{
		Brokers: cfg.Brokers,
		Logger:  logger,
	}
}

// EnsureTopicExists ensures a Kafka topic exists, creating it if necessary
func (km *KafkaManager) EnsureTopicExists(ctx context.Context, topicConfig TopicConfig) error {
	km.Logger.Info("Ensuring Kafka topic exists",
		zap.String("topic", topicConfig.Name),
		zap.Strings("brokers", km.Brokers))

	// Check if topic already exists
	exists, err := km.TopicExists(ctx, topicConfig.Name)
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if exists {
		km.Logger.Info("Kafka topic already exists", zap.String("topic", topicConfig.Name))
		return nil
	}

	// Create the topic
	return km.CreateTopic(ctx, topicConfig)
}

// TopicExists checks if a Kafka topic exists
func (km *KafkaManager) TopicExists(ctx context.Context, topicName string) (bool, error) {
	conn, err := kafka.DialContext(ctx, "tcp", km.Brokers[0])
	if err != nil {
		return false, fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	// Set a reasonable timeout for the metadata request
	if err := conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		km.Logger.Warn("Failed to set connection deadline", zap.Error(err))
	}

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return false, fmt.Errorf("failed to read partitions: %w", err)
	}

	for _, partition := range partitions {
		if partition.Topic == topicName {
			return true, nil
		}
	}

	return false, nil
}

// CreateTopic creates a new Kafka topic
func (km *KafkaManager) CreateTopic(ctx context.Context, topicConfig TopicConfig) error {
	km.Logger.Info("Creating Kafka topic",
		zap.String("topic", topicConfig.Name),
		zap.Int("partitions", topicConfig.NumPartitions),
		zap.Int("replication_factor", topicConfig.ReplicationFactor))

	conn, err := kafka.DialContext(ctx, "tcp", km.Brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	// Set connection deadline
	if err := conn.SetDeadline(time.Now().Add(30 * time.Second)); err != nil {
		km.Logger.Warn("Failed to set connection deadline", zap.Error(err))
	}

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topicConfig.Name,
			NumPartitions:     topicConfig.NumPartitions,
			ReplicationFactor: topicConfig.ReplicationFactor,
			ConfigEntries: []kafka.ConfigEntry{
				{
					ConfigName:  "retention.ms",
					ConfigValue: fmt.Sprintf("%d", topicConfig.RetentionMs),
				},
				{
					ConfigName:  "cleanup.policy",
					ConfigValue: "delete",
				},
				{
					ConfigName:  "compression.type",
					ConfigValue: "gzip",
				},
			},
		},
	}

	err = conn.CreateTopics(topicConfigs...)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	km.Logger.Info("Successfully created Kafka topic", zap.String("topic", topicConfig.Name))
	return nil
}

// HealthCheck performs a health check on Kafka connectivity
func (km *KafkaManager) HealthCheck(ctx context.Context) error {
	conn, err := kafka.DialContext(ctx, "tcp", km.Brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	// Set a timeout for the health check
	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		km.Logger.Warn("Failed to set connection deadline", zap.Error(err))
	}

	// Try to read metadata to verify connectivity
	_, err = conn.ReadPartitions()
	if err != nil {
		return fmt.Errorf("failed to read Kafka metadata: %w", err)
	}

	return nil
}

// GetDefaultTopicConfig returns default topic configuration for WanderSphere
func (km *KafkaManager) GetDefaultTopicConfig(topicName string) TopicConfig {
	return TopicConfig{
		Name:              topicName,
		NumPartitions:     3,                       // Good for load distribution
		ReplicationFactor: 1,                       // Single node setup
		RetentionMs:       7 * 24 * 60 * 60 * 1000, // 7 days retention
	}
}

// InitializeWanderSphereTopics creates all required topics for WanderSphere
func (km *KafkaManager) InitializeWanderSphereTopics(ctx context.Context, mainTopic string) error {
	topics := []TopicConfig{
		km.GetDefaultTopicConfig(mainTopic),
		km.GetDefaultTopicConfig(mainTopic + "_dlq"),   // Dead letter queue
		km.GetDefaultTopicConfig(mainTopic + "_retry"), // Retry topic
	}

	for _, topicConfig := range topics {
		if err := km.EnsureTopicExists(ctx, topicConfig); err != nil {
			return fmt.Errorf("failed to ensure topic %s exists: %w", topicConfig.Name, err)
		}
	}

	km.Logger.Info("All WanderSphere Kafka topics initialized successfully",
		zap.String("main_topic", mainTopic),
		zap.Int("total_topics", len(topics)))

	return nil
}

// ListTopics returns a list of all available Kafka topics
func (km *KafkaManager) ListTopics(ctx context.Context) ([]string, error) {
	conn, err := kafka.DialContext(ctx, "tcp", km.Brokers[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		km.Logger.Warn("Failed to set connection deadline", zap.Error(err))
	}

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, fmt.Errorf("failed to read partitions: %w", err)
	}

	topicSet := make(map[string]bool)
	for _, partition := range partitions {
		topicSet[partition.Topic] = true
	}

	topics := make([]string, 0, len(topicSet))
	for topic := range topicSet {
		topics = append(topics, topic)
	}

	return topics, nil
}
