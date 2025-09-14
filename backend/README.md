# WanderSphere Backend

This directory contains the Go backend for WanderSphere, a distributed social media platform with user authentication, post management, social relationships, ranked newsfeeds, and operational observability.

For the full project overview and architecture diagram, see the [root README](../README.md).

## Backend Overview

The backend is organized as 4 independently runnable Go services:

| Service | Port | Health | Responsibility |
|---------|------|--------|----------------|
| Web API | 19003 | `/health` | Public REST API, request middleware, auth session handling, uploads, rate limiting |
| AuthPost | 19001 | `:19101/health` | Users, authentication, posts, comments, likes, follows |
| Newsfeed | 19002 | `:19102/health` | Feed reads, Redis-backed ranked pagination, fallback feed reads |
| Newsfeed Publishing | 19004 | `:19104/health` | Kafka event processing, feed fan-out, engagement score updates, deduplication |

Supporting infrastructure:

| Component | Port | Purpose |
|-----------|------|---------|
| PostgreSQL | 5434 | Relational data for users, posts, comments, likes, and follows |
| Redis | 6379 | Sessions, materialized feeds, rate limits, circuit-breaker state, deduplication keys |
| Kafka | 9092 | Asynchronous feed and engagement events |
| MinIO | 9000/9001 | S3-compatible media storage for local development |
| Prometheus | 9090 | Metrics collection |
| Grafana | 3001 | Metrics dashboards |

## Core Capabilities

- gRPC service boundaries between the API gateway, AuthPost, Newsfeed, and Newsfeed Publishing services.
- Kafka-backed event flow for post creation and engagement updates.
- Redis-materialized timelines using sorted sets for ranked feeds and lists as fallback storage.
- Cursor-based pagination for stable feed reads as new posts arrive.
- Per-host gRPC circuit breakers with exponential backoff for downstream failures.
- Redis `SetNX` deduplication keys for Kafka consumers, with failed events routed to dead-letter topics.
- Redis-backed API rate limiting with request headers for limit, remaining capacity, and reset time.
- Health endpoints and Prometheus metrics across backend services.

## Local Development

### Prerequisites

- Docker and Docker Compose
- Go for running services or tests locally
- `migrate` CLI for manual migration commands
- `protoc` and Go protobuf plugins when regenerating gRPC code

### Start Everything

From the backend directory:

```bash
make start
```

This starts infrastructure, applies migrations, starts application services, and runs a health check.

### Stop Everything

```bash
make stop
```

### Start Step by Step

```bash
make infra
make migrate
make services
make health
```

### Run One Service Locally

```bash
make dev-webapp
make dev-authpost
make dev-newsfeed
make dev-nfp
```

When running services locally, make sure the required infrastructure is already running and that `config.yaml` points to the expected hostnames for your environment.

## Command Reference

| Command | Purpose |
|---------|---------|
| `make help` | Show available backend commands |
| `make start` | Start the complete backend stack |
| `make stop` | Stop services and remove local volumes |
| `make health` | Check service health endpoints |
| `make test` | Run backend unit tests |
| `make test-api` | Run API/integration tests through `tests/run_tests.sh` |
| `make test-integration` | Run integration tests in `tests/` |
| `make test-coverage` | Run Go tests with coverage |
| `make docs` | Regenerate Swagger documentation |
| `make proto` | Regenerate all protobuf files |
| `make new-migration MESSAGE_NAME=<name>` | Create a new SQL migration |
| `make migrate-up` | Apply migrations manually |
| `make migrate-down` | Roll back migrations manually |
| `make logs` | Show Docker Compose logs |
| `make logs-service SERVICE=<name>` | Show logs for one service |

## Testing Workflow

### Unit Tests

```bash
make test
```

### API and Integration Tests

```bash
make test-api
```

The API test suite starts the required services, runs migrations, tests authentication, post management, social graph operations, and newsfeed behavior, then produces test output under the test workflow.

### Coverage

```bash
make test-coverage
make test-coverage-html
```

### Optional Load-Test Scripts

The repository includes k6 scripts under `tests/load/` for local experiments. Treat generated load-test results as environment-specific and do not compare runs unless the profile, data set, machine, and service state are the same.

```bash
cd tests/load
./run_load_tests.sh smoke
./run_load_tests.sh load all
```

Generated files are written under `tests/load/results/<timestamp>/`.

## API Documentation

Swagger documentation is generated from the Web API service:

```bash
make docs
```

When the stack is running, open:

- Swagger UI: http://localhost:19003/swagger/index.html
- Swagger JSON: http://localhost:19003/swagger/doc.json

## Protobuf and gRPC

Protocol buffer definitions live under `pkg/types/proto/`.

```bash
make proto
```

Service-specific generation is also available:

```bash
make proto-authpost
make proto-newsfeed
make proto-nfp
```

## Configuration

Local service configuration is stored in `config.yaml`. The file defines service ports, database connection settings, Redis and Kafka addresses, S3-compatible storage settings, session behavior, rate limiting, CORS, request limits, pprof exposure, and Swagger exposure.

Keep local defaults scoped to development. For deployed environments, move secrets and environment-specific values out of static configuration and into an appropriate secret/configuration system.

## Operations

### Health Checks

```bash
make health
```

Individual health endpoints:

```bash
curl http://localhost:19003/health
curl http://localhost:19101/health
curl http://localhost:19102/health
curl http://localhost:19104/health
```

### Logs

```bash
make logs
make logs-service SERVICE=web
docker-compose logs -f postgres
docker-compose logs -f redis
docker-compose logs -f kafka
```

### Observability

Prometheus and Grafana are included in the local Docker Compose stack:

- Prometheus: http://localhost:9090
- Grafana: http://localhost:3001

Backend services expose metrics for HTTP requests, gRPC requests, rate limiting, circuit-breaker state, and Kafka message processing.

## Development Workflow

1. Create or update migrations when persistence changes.
2. Update protobuf definitions when service contracts change.
3. Regenerate protobuf code with `make proto`.
4. Implement service logic and middleware changes.
5. Run `make test` for unit coverage.
6. Run `make test-api` for cross-service behavior.
7. Regenerate Swagger docs with `make docs` when HTTP handlers change.

## Deployment Considerations

Before using this stack outside local development, review these items:

- Move passwords, object-storage credentials, and session settings out of local config files.
- Configure TLS for public endpoints and internal service communication where required.
- Enable authentication and network restrictions for Redis, PostgreSQL, Kafka, MinIO, Prometheus, and Grafana.
- Define backup and restore workflows for PostgreSQL, Redis persistence, MinIO objects, and Kafka topics.
- Add alerting for service health, request latency, error rates, circuit-breaker state, Kafka lag, and storage capacity.
- Tune connection pools, Kafka partitions, Redis memory policy, and Docker resource limits for the target environment.

## Troubleshooting

| Issue | Suggested Fix |
|-------|---------------|
| Services do not start | Run `make logs` and inspect infrastructure services first |
| Database migrations fail | Confirm PostgreSQL is healthy and the migration CLI is installed |
| Kafka is unavailable | Check `docker-compose logs zookeeper kafka` |
| Redis-backed features fail | Check `docker-compose logs redis` |
| API gateway cannot reach gRPC services | Confirm service containers are running and ports `19001-19004` are available |
| Port conflicts | Stop the conflicting local process or change the mapped port in Docker Compose |

## Repository Pointers

| Path | Purpose |
|------|---------|
| `cmd/` | Service entry points and Dockerfiles |
| `internal/app/` | Application service implementations |
| `internal/middleware/` | HTTP/gRPC middleware, rate limiting, security, circuit breakers |
| `internal/metrics/` | Prometheus metric definitions and middleware |
| `internal/utils/` | Shared infrastructure utilities |
| `pkg/types/proto/` | Protocol buffer definitions and generated gRPC code |
| `migrations/` | SQL migrations |
| `tests/` | API and integration tests |
| `tests/load/` | k6 load-test scripts |
| `prometheus/` | Prometheus configuration |
| `grafana/` | Grafana provisioning and dashboards |
