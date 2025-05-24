# WanderSphere Backend

A microservices-based social travel platform built with Go, featuring user authentication, post management, social features, and real-time newsfeed.

## 🚀 Quick Start

Get the entire system running with just one command:

```bash
git clone <repository-url>
cd backend
make start
```

That's it! This will automatically:
- Start all infrastructure services (PostgreSQL, Redis, Kafka)
- Run database migrations
- Start all application services  
- Verify everything is working

**🌐 Access Points:**
- API: http://localhost:19003/api/v1
- Swagger UI: http://localhost:19003/swagger/index.html
- Health Check: `make health`

## 📋 Prerequisites

Before you begin, ensure you have:

- **Docker & Docker Compose** - For running services
- **Go 1.19+** - For development and testing
- **migrate CLI** - For database migrations

```bash
# Install migrate CLI (Linux)
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

## 🏗️ Architecture

WanderSphere consists of 4 microservices and supporting infrastructure:

### Microservices
| Service | Port | Health Check | Description |
|---------|------|--------------|-------------|
| **Web API** | 19003 | :19103/health | Main REST API gateway |
| **AuthPost** | 19001 | :19101/health | Authentication & Posts |
| **Newsfeed** | 19002 | :19102/health | Timeline & Social Features |
| **Newsfeed Publishing** | 19004 | :19104/health | Event Processing |

### Infrastructure
- **PostgreSQL**: Port 5434 - Primary database
- **Redis**: Port 6379 - Caching and sessions  
- **Kafka**: Port 9092 - Event streaming

## 📚 Commands Reference

Run `make help` to see all available commands.

### 🚀 Essential Commands

```bash
make start      # Start the entire system from scratch
make stop       # Stop the entire system
make health     # Check if all services are healthy
make test-api   # Run comprehensive API tests
make dev        # Development mode (keeps services running)
```

### ⚙️ Step-by-step Commands

```bash
make infra      # Start infrastructure services (postgres, redis, kafka)
make migrate    # Run database migrations
make services   # Start application services
```

### 🧪 Testing Commands

```bash
make test               # Run unit tests
make test-verbose       # Verbose unit tests
make test-coverage      # With coverage report
make test-coverage-html # Generate HTML coverage report
```

### 📊 Database Commands

```bash
make new-migration MESSAGE_NAME=create_users_table  # Create new migration
make migrate-up         # Run migrations manually
make migrate-down       # Rollback migrations
```

### 🛠️ Development Commands

```bash
# Run individual services locally (without Docker)
make dev-authpost       # Run AuthPost service
make dev-newsfeed       # Run Newsfeed service  
make dev-webapp         # Run Web API service
make dev-nfp           # Run NFP service

# Documentation and code generation
make docs              # Generate API documentation
make proto             # Regenerate all protobuf files
make proto-authpost    # Generate AuthPost protobuf
make proto-newsfeed    # Generate Newsfeed protobuf  
make proto-nfp         # Generate NFP protobuf

# Utilities
make clean             # Clean Docker artifacts
make deps              # Update Go dependencies
```

### 🔧 Advanced Commands

```bash
make rebuild           # Rebuild and start with fresh images
make run-fg           # Run system in foreground (for debugging)
make logs             # Show all service logs
make logs-service SERVICE=web  # Show specific service logs
```

## 🧪 Testing

### Automated Testing

The `make test-api` command runs a comprehensive test suite via `tests/run_tests.sh` that:

```bash
# Run full API test suite
make test-api

# Keep services running for manual testing afterward
KEEP_RUNNING=true make test-api
```

The test script will:
1. Start all services automatically using `make infra`, `make migrate`, `make services`
2. Run database migrations
3. Test authentication endpoints (TestUser)
4. Test post management (TestCreate)
5. Test social features (TestFollow)
6. Test newsfeed functionality (TestNewsfeed)
7. Run integration tests (TestComplete)
8. Generate test reports

### Manual Testing

1. Start the system: `make start`
2. Open Swagger UI: http://localhost:19003/swagger/index.html
3. Test endpoints interactively

### Unit Tests

```bash
make test               # Run backend unit tests with: go test ./...
make test-integration   # Run integration tests in tests/ folder  
make test-all          # Run all tests (unit + integration)
make test-verbose       # Verbose output with: go test -v ./...
make test-coverage      # Coverage report with: go test -cover ./...
make test-coverage-html # Generate HTML coverage report
```

## 🐳 Docker Management

### Service Profiles

Start specific service groups:

```bash
# Infrastructure only
docker-compose --profile infra up

# Specific services  
docker-compose --profile web up        # Web API only
docker-compose --profile aap up        # AuthPost only
docker-compose --profile newsfeed up   # Newsfeed only
docker-compose --profile nfp up        # Publishing only

# Everything
docker-compose --profile all up
```

### Cleanup

```bash
make stop        # Stop and remove all containers/volumes
make clean       # Clean Docker artifacts
```

## 📊 Monitoring & Health Checks

### Quick Health Check

```bash
make health
```

Output example:
```
🏥 Checking service health...
AuthPost (19101): ✅ OK
Newsfeed (19102): ✅ OK  
Web API (19103): ✅ OK
NFP (19104): ✅ OK
```

### Individual Service Health

```bash
curl http://localhost:19101/health  # AuthPost
curl http://localhost:19102/health  # Newsfeed
curl http://localhost:19103/health  # Web API
curl http://localhost:19104/health  # Publishing
```

### Viewing Logs

```bash
make logs                           # All services
make logs-service SERVICE=web       # Specific service
docker-compose logs -f postgres     # Infrastructure service
```

## 🔧 Configuration

### Environment Configuration

Key settings are in `config.yaml`:

```yaml
database:
  host: localhost
  port: 5434
  name: wander_sphere
  user: postgres
  password: 123456

redis:
  host: localhost
  port: 6379

kafka:
  brokers: ["localhost:9092"]
```

### Service Ports

- **Application Services**: 19001-19004
- **Health Check Endpoints**: 19101-19104  
- **Main API**: 19003
- **Database**: 5434 (PostgreSQL)
- **Cache**: 6379 (Redis)
- **Messaging**: 9092 (Kafka)

## 🔄 Development Workflow

### Adding New Features

1. **Create database migration:**
   ```bash
   make new-migration MESSAGE_NAME=add_user_profiles
   ```

2. **Update protobuf definitions:**
   ```bash
   # Edit .proto files in pkg/types/proto/
   make proto
   ```

3. **Implement service logic**

4. **Test your changes:**
   ```bash
   make test-api
   ```

5. **Update documentation:**
   ```bash
   make docs
   ```

### Code Generation

```bash
make proto              # Regenerate all protobuf files
make proto-authpost     # Individual service protos
make docs              # Generate Swagger documentation
```

### Local Development

```bash
# Run services locally (outside Docker)
make dev-webapp         # Start web service locally
make dev-authpost       # Start auth service locally

# Or use development mode with Docker
make dev               # All services in Docker with auto-restart
```

## 🚨 Troubleshooting

### Common Issues

**Services not starting:**
```bash
# Check Docker resources
docker system df
docker system prune

# Check service logs
make logs-service SERVICE=postgres
make logs-service SERVICE=redis
```

**Database migration failures:**
```bash
# Check database connection
docker-compose exec postgres psql -U postgres -d wander_sphere -c "\l"

# Reset database (⚠️ DESTROYS DATA)
make stop
docker volume rm backend_postgres_data
make start
```

**Port conflicts:**
```bash
# Check what's using the ports
sudo netstat -tulpn | grep :5434
sudo netstat -tulpn | grep :6379

# Stop conflicting services
sudo systemctl stop postgresql
sudo systemctl stop redis
```

### Complete System Reset

If you need to start completely fresh:

```bash
make stop
docker system prune -a --volumes
make start
```

### Getting Help

```bash
make help              # Show all available commands
make logs              # Check service logs
make health            # Verify service status
```

## 📈 Production Considerations

### Security Checklist

- [ ] Change default passwords in `config.yaml`
- [ ] Use environment variables for secrets
- [ ] Enable PostgreSQL SSL in production
- [ ] Configure Redis AUTH
- [ ] Set up Kafka SASL/SSL
- [ ] Enable HTTPS for API endpoints

### Performance & Scaling

- [ ] Configure connection pooling
- [ ] Set up Redis clustering
- [ ] Configure PostgreSQL read replicas
- [ ] Use Kafka partitioning for load distribution
- [ ] Implement horizontal service scaling

### Monitoring & Logging

- [ ] Integrate Prometheus/Grafana
- [ ] Set up centralized logging (ELK stack)
- [ ] Configure alerting for service failures
- [ ] Monitor database performance
- [ ] Track Kafka consumer lag

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Run the test suite: `make test-api`
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Development Setup

```bash
# Clone and setup
git clone <your-fork-url>
cd backend
make start

# Make changes and test
make test-api
make test

# Generate docs if API changed
make docs
```
