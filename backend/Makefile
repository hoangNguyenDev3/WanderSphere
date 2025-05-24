# WanderSphere Backend - Makefile
# Simple commands to run the system from scratch

# Database configuration
USERNAME = "postgres"
PASSWORD = "123456"
DB_HOST = "localhost"
DB_PORT = "5434"
DB_NAME = "wander_sphere"

#################################
# Quick Start Commands
#################################

# Show all available commands
help:
	@echo "🏗️  WanderSphere Backend - Available Commands"
	@echo "============================================="
	@echo ""
	@echo "🚀 Essential Commands:"
	@echo "  make start       - Start the entire system from scratch"
	@echo "  make stop        - Stop the entire system"
	@echo "  make health      - Check if all services are healthy"
	@echo "  make test-api    - Run comprehensive API tests"
	@echo "  make dev         - Development mode (keeps services running)"
	@echo ""
	@echo "⚙️  Step-by-step Commands:"
	@echo "  make infra       - Start infrastructure services (postgres, redis, kafka)"
	@echo "  make migrate     - Run database migrations"
	@echo "  make services    - Start application services"
	@echo ""
	@echo "🧪 Testing Commands:"
	@echo "  make test        - Run backend unit tests"
	@echo "  make test-integration - Run integration tests (tests/ folder)"
	@echo "  make test-all    - Run all tests (unit + integration)"
	@echo "  make test-verbose - Run backend tests with verbose output"
	@echo "  make test-coverage - Run backend tests with coverage report"
	@echo ""
	@echo "📊 Database Commands:"
	@echo "  make new-migration MESSAGE_NAME=name - Create new migration"
	@echo "  make migrate-up  - Run migrations manually"
	@echo "  make migrate-down - Rollback migrations"
	@echo ""
	@echo "🛠️  Development Commands:"
	@echo "  make dev-authpost - Run AuthPost service locally"
	@echo "  make dev-newsfeed - Run Newsfeed service locally"
	@echo "  make dev-webapp  - Run Web API service locally"
	@echo "  make dev-nfp     - Run NFP service locally"
	@echo ""
	@echo "📚 Documentation:"
	@echo "  make docs        - Generate API documentation"
	@echo ""
	@echo "🧹 Cleanup:"
	@echo "  make clean       - Clean Docker artifacts"
	@echo ""
	@echo "ℹ️  For more details, see the README.md"

# Start the entire system from scratch (recommended for new users)
start: infra migrate services health
	@echo ""
	@echo "🎉 WanderSphere system is now running!"
	@echo "API Base URL: http://localhost:19003/api/v1"
	@echo "Swagger UI: http://localhost:19003/swagger/index.html"
	@echo ""
	@echo "To run tests: make test-api"
	@echo "To stop: make stop"

# Stop the entire system
stop:
	@echo "🛑 Stopping WanderSphere system..."
	docker-compose down --volumes --remove-orphans
	@echo "✅ System stopped"

# Start infrastructure services only (postgres, redis, kafka)
infra:
	@echo "🔧 Starting infrastructure services..."
	docker-compose --profile infra up -d
	@echo "⏳ Waiting for services to be ready..."
	@sleep 30
	@echo "✅ Infrastructure services started"

# Start all application services 
services:
	@echo "🚀 Starting application services..."
	docker-compose --profile all up -d
	@echo "⏳ Waiting for services to start..."
	@sleep 45
	@echo "✅ Application services started"

# Run database migrations
migrate:
	@echo "📊 Running database migrations..."
	migrate -path migrations/ -database "postgresql://$(USERNAME):$(PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up
	@echo "✅ Database migrations completed"

# Run comprehensive API tests using run_tests.sh script
test-api:
	@echo "🧪 Running API tests..."
	cd tests && ./run_tests.sh
	@echo "✅ API tests completed"

# Quick health check for all services
health:
	@echo "🏥 Checking service health..."
	@echo -n "AuthPost (19101): " && curl -s --max-time 3 http://localhost:19101/health >/dev/null && echo "✅ OK" || echo "❌ Down"
	@echo -n "Newsfeed (19102): " && curl -s --max-time 3 http://localhost:19102/health >/dev/null && echo "✅ OK" || echo "❌ Down"
	@echo -n "Web API (19003): " && curl -s --max-time 3 http://localhost:19003/health >/dev/null && echo "✅ OK" || echo "❌ Down"
	@echo -n "NFP (19104): " && curl -s --max-time 3 http://localhost:19104/health >/dev/null && echo "✅ OK" || echo "❌ Down"

# Development mode - keeps services running for manual testing
dev: start
	@echo "🔧 Development mode: Services will keep running"
	@echo "To stop: make stop"

#################################
# Database Management
#################################

# Create a new migration
new-migration:
	@if [ -z "$(MESSAGE_NAME)" ]; then echo "❌ Please provide MESSAGE_NAME: make new-migration MESSAGE_NAME=create_users_table"; exit 1; fi
	migrate create -ext sql -dir migrations/ -seq $(MESSAGE_NAME)

# Run migrations manually
migrate-up:
	migrate -path migrations/ -database "postgresql://$(USERNAME):$(PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

# Rollback migrations
migrate-down:
	migrate -path migrations/ -database "postgresql://$(USERNAME):$(PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down

#################################
# Testing Commands
#################################

# Run backend unit tests only
test:
	go test ./...

# Run integration tests in tests/ folder only  
test-integration: start
	@echo "🧪 Running integration tests..."
	cd tests && go test -v ./...

# Run all tests (backend unit tests + integration tests)
test-all:
	@echo "🧪 Running all tests..."
	-go test ./...
	@echo "🧪 Running integration tests..."
	cd tests && go test -v ./...
	@echo "✅ All tests completed"

# Run backend tests with verbose output
test-verbose:
	go test -v ./...

# Run backend tests with coverage report
test-coverage:
	go test -cover ./...
	
# Run backend tests with coverage report and output to HTML
test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

#################################
# Local Development
#################################

# Run individual services locally (without Docker)
dev-authpost:
	go run cmd/authpost/main.go

dev-newsfeed:
	go run cmd/newsfeed/main.go

dev-webapp:
	go run cmd/webapp/main.go

dev-nfp:
	go run cmd/newsfeed_publishing/main.go

#################################
# Documentation & Protobuf
#################################

# Generate API documentation
docs:
	swag init -g cmd/webapp/main.go --parseDependency --parseInternal --exclude .git,logs,personal --output docs
	@echo "📚 Swagger documentation generated!"
	@echo "JSON: http://localhost:19003/swagger/doc.json"
	@echo "UI: http://localhost:19003/swagger/index.html"

# Regenerate all protobuf files
proto: 
	chmod +x scripts/regenerate_protos.sh
	./scripts/regenerate_protos.sh

# Individual protobuf generation
proto-authpost:
	protoc --go_out=./pkg/types/proto/pb/authpost/ --go_opt=paths=source_relative --go-grpc_out=./pkg/types/proto/pb/authpost/ --go-grpc_opt=paths=source_relative ./pkg/types/proto/authpost.proto 

proto-newsfeed:
	protoc --go_out=./pkg/types/proto/pb/newsfeed/ --go_opt=paths=source_relative --go-grpc_out=./pkg/types/proto/pb/newsfeed/ --go-grpc_opt=paths=source_relative ./pkg/types/proto/newsfeed.proto 

proto-nfp:
	protoc --go_out=./pkg/types/proto/pb/newsfeed_publishing/ --go_opt=paths=source_relative --go-grpc_out=./pkg/types/proto/pb/newsfeed_publishing/ --go-grpc_opt=paths=source_relative ./pkg/types/proto/newsfeed_publishing.proto 

#################################
# Cleanup & Utilities
#################################

# Clean Docker artifacts
clean:
	@echo "🧹 Cleaning Docker artifacts..."
	docker system prune -f
	docker volume prune -f
	@echo "✅ Cleanup completed"

# Update Go dependencies
deps:
	go mod tidy
	go mod vendor -v

#################################
# Advanced Commands (for CI/CD and power users)
#################################

# Build and start with fresh images
rebuild:
	@echo "🔨 Rebuilding and starting system..."
	docker-compose --profile all up --build --force-recreate -d
	@echo "✅ System rebuilt and started"

# Run system in foreground (for debugging)
run-fg:
	docker-compose --profile all up

# Show service logs
logs:
	docker-compose logs -f

# Show logs for specific service
logs-service:
	@if [ -z "$(SERVICE)" ]; then echo "❌ Please specify SERVICE: make logs-service SERVICE=web"; exit 1; fi
	docker-compose logs -f $(SERVICE)

.PHONY: help start stop infra services migrate test-api health dev new-migration migrate-up migrate-down test test-integration test-all test-verbose test-coverage test-coverage-html dev-authpost dev-newsfeed dev-webapp dev-nfp docs proto proto-authpost proto-newsfeed proto-nfp clean deps rebuild run-fg logs logs-service