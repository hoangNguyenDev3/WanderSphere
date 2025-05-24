package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Define command line flags
	var (
		configPath = flag.String("config", "/app/config.yaml", "Path to config file")
		command    = flag.String("cmd", "help", "Command to execute: help, migration-status, migration-up, migration-down, migration-reset, kafka-topics, kafka-create-topic, redis-status")
		topicName  = flag.String("topic", "", "Kafka topic name for topic operations")
		service    = flag.String("service", "authpost", "Service to operate on: authpost, newsfeed, newsfeed_publishing, webapp")
	)
	flag.Parse()

	fmt.Println("🚀 WanderSphere System Administration Tool")
	fmt.Println("==========================================")

	switch *command {
	case "help":
		printHelp()
	case "migration-status":
		handleMigrationStatus(*configPath)
	case "migration-up":
		handleMigrationUp(*configPath)
	case "migration-down":
		handleMigrationDown(*configPath)
	case "migration-reset":
		handleMigrationReset(*configPath)
	case "kafka-topics":
		handleKafkaTopics(*configPath, *service)
	case "kafka-create-topic":
		if *topicName == "" {
			log.Fatal("Topic name is required for kafka-create-topic command")
		}
		handleKafkaCreateTopic(*configPath, *service, *topicName)
	case "redis-status":
		handleRedisStatus(*configPath, *service)
	default:
		fmt.Printf("Unknown command: %s\n", *command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`
Usage: system_admin -cmd <command> [options]

Commands:
  help               Show this help message
  migration-status   Show current migration status
  migration-up       Run pending migrations
  migration-down     Rollback last migration
  migration-reset    Reset database (DANGER: drops all data)
  kafka-topics       List all Kafka topics
  kafka-create-topic Create a new Kafka topic (requires -topic flag)
  redis-status       Show Redis connection pool status

Options:
  -config <path>     Path to config file (default: /app/config.yaml)
  -service <name>    Service name for Redis/Kafka operations (authpost, newsfeed, newsfeed_publishing, webapp)
  -topic <name>      Topic name for Kafka operations

Examples:
  # Check migration status
  system_admin -cmd migration-status

  # Run migrations
  system_admin -cmd migration-up

  # List Kafka topics
  system_admin -cmd kafka-topics -service newsfeed_publishing

  # Create new topic
  system_admin -cmd kafka-create-topic -service newsfeed_publishing -topic test_topic

  # Check Redis status
  system_admin -cmd redis-status -service webapp`)
}

func handleMigrationStatus(configPath string) {
	fmt.Println("📊 Checking migration status...")

	cfg, err := configs.GetAuthenticateAndPostConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := connectToDatabase(&cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer closeDatabase(db)

	logger, _ := utils.NewLogger(&cfg.Logger)
	migrationManager := utils.NewMigrationManager(db, &cfg.Postgres, logger)

	status, err := migrationManager.GetStatus("migrations")
	if err != nil {
		log.Fatalf("Failed to get migration status: %v", err)
	}

	fmt.Printf("Total migrations: %d\n", status.TotalMigrations)
	fmt.Printf("Applied migrations: %d\n", status.AppliedMigrations)
	fmt.Printf("Last applied: %s\n", status.LastApplied)

	if len(status.PendingMigrations) > 0 {
		fmt.Printf("Pending migrations: %s\n", strings.Join(status.PendingMigrations, ", "))
	} else {
		fmt.Println("✅ No pending migrations")
	}
}

func handleMigrationUp(configPath string) {
	fmt.Println("⬆️  Running migrations...")

	cfg, err := configs.GetAuthenticateAndPostConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := connectToDatabase(&cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer closeDatabase(db)

	logger, _ := utils.NewLogger(&cfg.Logger)
	migrationManager := utils.NewMigrationManager(db, &cfg.Postgres, logger)

	if err := migrationManager.Migrate("migrations"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("✅ Migrations completed successfully")
}

func handleMigrationDown(configPath string) {
	fmt.Println("⬇️  Rolling back last migration...")

	cfg, err := configs.GetAuthenticateAndPostConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := connectToDatabase(&cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer closeDatabase(db)

	logger, _ := utils.NewLogger(&cfg.Logger)
	migrationManager := utils.NewMigrationManager(db, &cfg.Postgres, logger)

	if err := migrationManager.Rollback(); err != nil {
		log.Fatalf("Rollback failed: %v", err)
	}

	fmt.Println("✅ Migration rolled back successfully")
}

func handleMigrationReset(configPath string) {
	fmt.Println("⚠️  DANGER: Resetting database...")
	fmt.Print("Are you sure? This will delete ALL data! (yes/no): ")

	var confirmation string
	fmt.Scanln(&confirmation)

	if confirmation != "yes" {
		fmt.Println("Operation cancelled")
		return
	}

	cfg, err := configs.GetAuthenticateAndPostConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := connectToDatabase(&cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer closeDatabase(db)

	logger, _ := utils.NewLogger(&cfg.Logger)
	migrationManager := utils.NewMigrationManager(db, &cfg.Postgres, logger)

	if err := migrationManager.Reset("migrations"); err != nil {
		log.Fatalf("Reset failed: %v", err)
	}

	fmt.Println("✅ Database reset completed")
}

func handleKafkaTopics(configPath, service string) {
	fmt.Printf("📋 Listing Kafka topics for %s service...\n", service)

	kafkaConfig, err := getKafkaConfig(configPath, service)
	if err != nil {
		log.Fatalf("Failed to get Kafka config: %v", err)
	}

	logger, _ := utils.NewLogger(&configs.LoggerConfig{Level: "info"})
	kafkaManager := utils.NewKafkaManager(kafkaConfig, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	topics, err := kafkaManager.ListTopics(ctx)
	if err != nil {
		log.Fatalf("Failed to list topics: %v", err)
	}

	fmt.Printf("Found %d topics:\n", len(topics))
	for _, topic := range topics {
		fmt.Printf("  - %s\n", topic)
	}
}

func handleKafkaCreateTopic(configPath, service, topicName string) {
	fmt.Printf("🔧 Creating Kafka topic '%s' for %s service...\n", topicName, service)

	kafkaConfig, err := getKafkaConfig(configPath, service)
	if err != nil {
		log.Fatalf("Failed to get Kafka config: %v", err)
	}

	logger, _ := utils.NewLogger(&configs.LoggerConfig{Level: "info"})
	kafkaManager := utils.NewKafkaManager(kafkaConfig, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	topicConfig := kafkaManager.GetDefaultTopicConfig(topicName)
	if err := kafkaManager.CreateTopic(ctx, topicConfig); err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}

	fmt.Printf("✅ Topic '%s' created successfully\n", topicName)
}

func handleRedisStatus(configPath, service string) {
	fmt.Printf("🔍 Checking Redis status for %s service...\n", service)

	redisConfig, err := getRedisConfig(configPath, service)
	if err != nil {
		log.Fatalf("Failed to get Redis config: %v", err)
	}

	logger, _ := utils.NewLogger(&configs.LoggerConfig{Level: "info"})
	redisPool, err := utils.NewRedisPool(redisConfig, logger)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisPool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisPool.HealthCheck(ctx); err != nil {
		log.Fatalf("Redis health check failed: %v", err)
	}

	stats := redisPool.Client.PoolStats()
	fmt.Printf("✅ Redis connection pool status:\n")
	fmt.Printf("  Total connections: %d\n", stats.TotalConns)
	fmt.Printf("  Idle connections: %d\n", stats.IdleConns)
	fmt.Printf("  Stale connections: %d\n", stats.StaleConns)
	fmt.Printf("  Hits: %d\n", stats.Hits)
	fmt.Printf("  Misses: %d\n", stats.Misses)
	fmt.Printf("  Timeouts: %d\n", stats.Timeouts)
}

// Helper functions
func connectToDatabase(cfg *configs.PostgresConfig) (*gorm.DB, error) {
	postgresConfig := postgres.Config{DSN: cfg.DSN}
	return gorm.Open(postgres.New(postgresConfig), &gorm.Config{})
}

func closeDatabase(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
}

func getKafkaConfig(configPath, service string) (*configs.KafkaConfig, error) {
	switch service {
	case "newsfeed_publishing":
		cfg, err := configs.GetNewsfeedPublishingConfig(configPath)
		if err != nil {
			return nil, err
		}
		return &cfg.Kafka, nil
	case "newsfeed":
		cfg, err := configs.GetNewsfeedConfig(configPath)
		if err != nil {
			return nil, err
		}
		return &cfg.Kafka, nil
	default:
		return nil, fmt.Errorf("service %s does not use Kafka", service)
	}
}

func getRedisConfig(configPath, service string) (*configs.RedisConfig, error) {
	switch service {
	case "webapp":
		cfg, err := configs.GetWebConfig(configPath)
		if err != nil {
			return nil, err
		}
		return &cfg.Redis, nil
	case "newsfeed_publishing":
		cfg, err := configs.GetNewsfeedPublishingConfig(configPath)
		if err != nil {
			return nil, err
		}
		return &cfg.Redis, nil
	case "newsfeed":
		cfg, err := configs.GetNewsfeedConfig(configPath)
		if err != nil {
			return nil, err
		}
		return &cfg.Redis, nil
	default:
		return nil, fmt.Errorf("service %s does not use Redis", service)
	}
}
