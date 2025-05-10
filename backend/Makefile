USERNAME = "postgres"
PASSWORD = "123456"
DB_HOST = "localhost"
DB_PORT = "5432"
DB_NAME = "wander_sphere"

# Database migrations
new_migration:
	migrate create -ext sql -dir scripts/migration/ -seq $(MESSAGE_NAME)
up_migration:
	migrate -path scripts/migration/ -database "postgres://$(USERNAME):$(PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up
down_migration:
	migrate -path scripts/migration/ -database "postgres://$(USERNAME):$(PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down

# Protobuf generation
proto_aap:
	protoc --proto_path=pkg/types/proto/ --go_out=pkg/types/proto/pb/authen_and_post --go_opt=paths=source_relative \
	--go-grpc_out=pkg/types/proto/pb/authen_and_post --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional  \
	authen_and_post.proto
proto_newsfeed:
	protoc --proto_path=pkg/types/proto/ --go_out=pkg/types/proto/pb/newsfeed --go_opt=paths=source_relative \
	--go-grpc_out=pkg/types/proto/pb/newsfeed --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional \
	newsfeed.proto
proto_newsfeed_publishing:
	protoc --proto_path=pkg/types/proto/ --go_out=pkg/types/proto/pb/newsfeed_publishing --go_opt=paths=source_relative \
	--go-grpc_out=pkg/types/proto/pb/newsfeed_publishing --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional \
	newsfeed_publishing.proto

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
	swag init -g cmd/webapp/main.go