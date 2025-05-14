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
	@echo "📊 Swagger documentation generated at docs/"
	@echo "📋 JSON: http://localhost:19003/swagger/doc.json"
	@echo "🌐 UI: http://localhost:19003/swagger/index.html"
	@echo "💡 Make sure the web service is running: make docker_webapp"

# Health checks for all services
health_check:
	@echo "🔍 Checking service health..."
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