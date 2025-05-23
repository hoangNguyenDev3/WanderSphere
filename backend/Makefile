USERNAME = "postgres"
PASSWORD = "123456"
DB_HOST = "localhost"
DB_PORT = "5434"
DB_NAME = "wander_sphere"

# Database management
create_db:
	PGPASSWORD=$(PASSWORD) createdb -h $(DB_HOST) -p $(DB_PORT) -U $(USERNAME) $(DB_NAME) || true

# Database migrations
new_migration:
	migrate create -ext sql -dir migrations/ -seq $(MESSAGE_NAME)
up_migration: create_db
	migrate -path migrations/ -database "postgresql://$(USERNAME):$(PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up
down_migration:
	migrate -path migrations/ -database "postgresql://$(USERNAME):$(PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down

# Protobuf generation
proto: 
	chmod +x scripts/regenerate_protos.sh
	./scripts/regenerate_protos.sh

proto_aap:
	protoc --go_out=./pkg/types/proto/pb/authpost/ --go_opt=paths=source_relative --go-grpc_out=./pkg/types/proto/pb/authpost/ --go-grpc_opt=paths=source_relative ./pkg/types/proto/authpost.proto 
proto_newsfeed:
	protoc --go_out=./pkg/types/proto/pb/newsfeed/ --go_opt=paths=source_relative --go-grpc_out=./pkg/types/proto/pb/newsfeed/ --go-grpc_opt=paths=source_relative ./pkg/types/proto/newsfeed.proto 
proto_newsfeed_publishing:
	protoc --go_out=./pkg/types/proto/pb/newsfeed_publishing/ --go_opt=paths=source_relative --go-grpc_out=./pkg/types/proto/pb/newsfeed_publishing/ --go-grpc_opt=paths=source_relative ./pkg/types/proto/newsfeed_publishing.proto 

# Dependencies
.PHONY: vendor
vendor:
	go mod vendor -v

# Docker management
docker_clear:
	docker volume rm $(docker volume ls -qf dangling=true) & docker rmi $(docker images -qf "dangling=true")

# Docker Compose - full stack
compose_up:
	docker-compose --profile all up
compose_up_rebuild:
	docker-compose --profile all up --build --force-recreate
compose_down:
	docker-compose down

# Docker Compose - individual services
docker_authpost:
	docker-compose --profile aap --profile infra up
docker_newsfeed:
	docker-compose --profile newsfeed --profile infra up
docker_newsfeed_publishing:
	docker-compose --profile nfp --profile infra up
docker_webapp:
	docker-compose --profile web up
docker_infra:
	docker-compose --profile infra up

# Local development (direct run)
dev_authpost:
	go run cmd/authpost/main.go
dev_newsfeed:
	go run cmd/newsfeed/main.go
dev_newsfeed_publishing:
	go run cmd/newsfeed_publishing/main.go
dev_webapp:
	go run cmd/webapp/main.go

# Documentation
gen_swagger:
	swag init -g cmd/webapp/main.go --parseDependency --parseInternal --exclude .git,logs,personal --output docs

# Serve swagger documentation (requires swagger-ui or similar tool)
serve_swagger:
	@echo "Swagger documentation generated at docs/"
	@echo "JSON: http://localhost:19003/swagger/doc.json"
	@echo "UI: http://localhost:19003/swagger/index.html"
	@echo "Make sure the web service is running: make docker_webapp"

# Health checks for all services
health_check:
	@echo "Checking service health..."
	@echo "AuthPost (19101):" && curl -s http://localhost:19101/health || echo "❌ Not responding"
	@echo "Newsfeed (19102):" && curl -s http://localhost:19102/health || echo "❌ Not responding"  
	@echo "Web (19003):" && curl -s http://localhost:19003/health || echo "❌ Not responding"
	@echo "NFP (19104):" && curl -s http://localhost:19104/health || echo "❌ Not responding"

# Test commands
.PHONY: test test-verbose test-coverage test-integration

# Run all tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run specific test suite
test-webapp:
	go test -v ./internal/app/webapp/tests

# Run tests with coverage report
test-coverage:
	go test -cover ./...
	
# Run tests with coverage report and output to HTML
test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run integration tests (requires running services)
test-integration:
	go test -v -tags=integration ./...

# System Administration Tool - Enhanced Management Commands
.PHONY: build-system-admin system-help system-migration-status system-migration-up system-migration-down system-kafka-topics system-kafka-create-topic system-redis-status

# Build the System Administration tool
build-system-admin:
	@echo "Building System Administration tool..."
	go build -o bin/system_admin cmd/system_admin/main.go
	@echo "System admin tool built at bin/system_admin"

# System admin help
system-help: build-system-admin
	./bin/system_admin -config config.yaml -cmd help

# Database migration management (automated)
system-migration-status: build-system-admin
	@echo "Checking migration status..."
	./bin/system_admin -config config.yaml -cmd migration-status

system-migration-up: build-system-admin
	@echo "Running database migrations..."
	./bin/system_admin -config config.yaml -cmd migration-up

system-migration-down: build-system-admin
	@echo "Rolling back last migration..."
	./bin/system_admin -config config.yaml -cmd migration-down

system-migration-reset: build-system-admin
	@echo "DANGER: Resetting database..."
	./bin/system_admin -config config.yaml -cmd migration-reset

# Kafka topic management
system-kafka-topics: build-system-admin
	@echo "Listing Kafka topics..."
	./bin/system_admin -config config.yaml -cmd kafka-topics -service newsfeed_publishing

system-kafka-create-topic: build-system-admin
	@echo "Creating Kafka topic (set TOPIC_NAME)..."
	@if [ -z "$(TOPIC_NAME)" ]; then echo "Please set TOPIC_NAME: make system-kafka-create-topic TOPIC_NAME=my_topic"; exit 1; fi
	./bin/system_admin -config config.yaml -cmd kafka-create-topic -service newsfeed_publishing -topic $(TOPIC_NAME)

# Redis connection pool monitoring
system-redis-status-webapp: build-system-admin
	@echo "Checking Redis pool status for Web service..."
	./bin/system_admin -config config.yaml -cmd redis-status -service webapp

system-redis-status-nfp: build-system-admin
	@echo "Checking Redis pool status for Newsfeed Publishing service..."
	./bin/system_admin -config config.yaml -cmd redis-status -service newsfeed_publishing

system-redis-status-newsfeed: build-system-admin
	@echo "Checking Redis pool status for Newsfeed service..."
	./bin/system_admin -config config.yaml -cmd redis-status -service newsfeed

# System comprehensive health check (includes Redis pools and Kafka topics)
system-health-check: build-system-admin
	@echo "System Enhanced Health Check"
	@echo "==============================="
	@echo ""
	@echo "Service Health:"
	@$(MAKE) health_check
	@echo ""
	@echo "Database Migration Status:"
	@$(MAKE) system-migration-status
	@echo ""
	@echo "Kafka Topics:"
	@$(MAKE) system-kafka-topics || echo "❌ Kafka not available"
	@echo ""
	@echo "Redis Pool Status:"
	@echo "Web Service:"
	@$(MAKE) system-redis-status-webapp || echo "❌ Web Redis not available"
	@echo "Newsfeed Publishing Service:"
	@$(MAKE) system-redis-status-nfp || echo "❌ NFP Redis not available"
	@echo "Newsfeed Service:"
	@$(MAKE) system-redis-status-newsfeed || echo "❌ Newsfeed Redis not available"

# Complete system setup (initializes everything)
system-setup: build-system-admin
	@echo "System Complete Setup"
	@echo "========================"
	@echo "1. Running database migrations..."
	@$(MAKE) system-migration-up
	@echo ""
	@echo "2. Creating default Kafka topics..."
	@$(MAKE) system-kafka-create-topic TOPIC_NAME=wander_sphere || echo "ℹ️  Topic may already exist"
	@$(MAKE) system-kafka-create-topic TOPIC_NAME=wander_sphere_dlq || echo "ℹ️  DLQ topic may already exist"
	@$(MAKE) system-kafka-create-topic TOPIC_NAME=wander_sphere_retry || echo "ℹ️  Retry topic may already exist"
	@echo ""
	@echo "3. Running comprehensive health check..."
	@$(MAKE) system-health-check
	@echo ""
	@echo "✅ System setup completed!"

# System development workflow
system-dev-start: system-setup
	@echo "Starting system development environment..."
	docker-compose --profile all up

# Clean system artifacts
system-clean:
	@echo "Cleaning system artifacts..."
	rm -f bin/system_admin
	@echo "System cleanup completed"

# Enhanced health checks (includes all system features)
health_check_enhanced: system-health-check