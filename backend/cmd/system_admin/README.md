# WanderSphere System Administration Tool

A comprehensive CLI tool for managing infrastructure features in WanderSphere backend.

## 🚀 Quick Start

```bash
# Build the system administration tool
make build-system-admin

# Get help
./bin/system_admin -cmd help

# Quick setup everything
make system-setup
```

## 📋 Available Commands

### Database Migration Management

```bash
# Check current migration status
./bin/system_admin -cmd migration-status
./bin/system_admin -cmd migration-status -config /path/to/config.yaml

# Run pending migrations
./bin/system_admin -cmd migration-up

# Rollback last migration
./bin/system_admin -cmd migration-down

# DANGER: Reset database (development only)
./bin/system_admin -cmd migration-reset
```

### Kafka Topic Management

```bash
# List all Kafka topics
./bin/system_admin -cmd kafka-topics -service newsfeed_publishing

# Create a new topic
./bin/system_admin -cmd kafka-create-topic -service newsfeed_publishing -topic my_new_topic

# Health check Kafka connectivity
./bin/system_admin -cmd kafka-topics -service newsfeed
```

### Redis Connection Pool Monitoring

```bash
# Check Redis pool status for Web service
./bin/system_admin -cmd redis-status -service webapp

# Check Redis pool status for Newsfeed Publishing service
./bin/system_admin -cmd redis-status -service newsfeed_publishing

# Check Redis pool status for Newsfeed service
./bin/system_admin -cmd redis-status -service newsfeed
```

## 🛠️ Command Line Options

| Option | Default | Description |
|--------|---------|-------------|
| `-config` | `/app/config.yaml` | Path to configuration file |
| `-cmd` | `help` | Command to execute |
| `-service` | `authpost` | Service to operate on (`authpost`, `newsfeed`, `newsfeed_publishing`, `webapp`) |
| `-topic` | `` | Kafka topic name for topic operations |

## 📊 Example Outputs

### Migration Status
```bash
$ ./bin/system_admin -cmd migration-status
📊 Checking migration status...
Total migrations: 1
Applied migrations: 1
Last applied: 0001
✅ No pending migrations
```

### Kafka Topics
```bash
$ ./bin/system_admin -cmd kafka-topics -service newsfeed_publishing
📋 Listing Kafka topics for newsfeed_publishing service...
Found 3 topics:
  - wander_sphere
  - wander_sphere_dlq
  - wander_sphere_retry
```

### Redis Pool Status
```bash
$ ./bin/system_admin -cmd redis-status -service webapp
🔍 Checking Redis status for webapp service...
✅ Redis connection pool status:
  Total connections: 5
  Idle connections: 3
  Stale connections: 0
  Hits: 150
  Misses: 12
  Timeouts: 0
```

## 🔧 Makefile Integration

For easier usage, use the integrated Makefile commands:

### Database Management
```bash
make system-migration-status    # Check migration status
make system-migration-up        # Run migrations
make system-migration-down      # Rollback last migration
make system-migration-reset     # Reset database (DANGER)
```

### Kafka Management
```bash
make system-kafka-topics                              # List topics
make system-kafka-create-topic TOPIC_NAME=test_topic # Create topic
```

### Redis Monitoring
```bash
make system-redis-status-webapp      # Web service Redis status
make system-redis-status-nfp         # Newsfeed Publishing service Redis status
make system-redis-status-newsfeed    # Newsfeed service Redis status
```

### Comprehensive Operations
```bash
make system-setup               # Complete system initialization
make system-health-check        # Enhanced health check with all features
make system-dev-start           # Full development environment with setup
```

## 🚨 Troubleshooting

### Common Issues

#### 1. Database Connection Failed
```bash
Error: Failed to connect to database: dial tcp [::1]:5432: connect: connection refused
```
**Solution**: Ensure PostgreSQL is running and accessible
```bash
docker-compose --profile infra up postgres
```

#### 2. Kafka Not Available
```bash
Error: Failed to list topics: dial tcp [::1]:9092: connect: connection refused
```
**Solution**: Start Kafka services
```bash
docker-compose --profile infra up kafka
```

#### 3. Redis Connection Failed
```bash
Error: Failed to connect to Redis: dial tcp [::1]:6379: connect: connection refused
```
**Solution**: Start Redis service
```bash
docker-compose --profile infra up redis
```

#### 4. Migration Files Not Found
```bash
Error: migration directory does not exist: migrations
```
**Solution**: Run from the correct directory (backend root) or specify config path
```bash
cd /path/to/backend
./bin/system_admin -cmd migration-status
```

### Service-Specific Troubleshooting

#### Newsfeed Publishing Service Issues
- Check Kafka connectivity: `make system-kafka-topics`
- Check Redis pool: `make system-redis-status-nfp` 
- Verify topics exist: Auto-created on startup or use `make system-kafka-create-topic`

#### Web Service Issues
- Check Redis pool: `make system-redis-status-webapp`
- Verify service config has correct Redis settings

#### Newsfeed Service Issues
- Check Redis pool: `make system-redis-status-newsfeed`
- Verify pagination works with enhanced pool

## 🔐 Security Considerations

### Database Operations
- **Migration Reset**: Only use in development environments
- **Configuration**: Ensure database credentials are properly secured
- **Backup**: Always backup before running migrations in production

### Access Control
- Limit access to the system administration tool to authorized personnel only
- Use secure configuration files with appropriate file permissions
- Monitor administrative actions through logs

### Network Security
- Ensure Redis and Kafka connections use secure channels when in production
- Regularly rotate credentials and API keys
- Implement proper firewall rules for database access

## 📝 Development Notes

### Code Structure
- Main executable: `cmd/system_admin/main.go`
- Configuration management through centralized config package
- Service-specific configuration loading based on service parameter
- Comprehensive error handling and user feedback

### Adding New Commands
1. Add the command to the main switch statement
2. Implement the handler function following existing patterns
3. Update help text and documentation
4. Add corresponding Makefile target if needed

### Testing
- Test all commands in development environment before production use
- Verify database backups before running destructive operations
- Test Kafka and Redis connectivity separately if issues arise

## 🏗️ Architecture

### Tool Structure
```
cmd/system_admin/
├── main.go           # Main CLI application
└── README.md         # This documentation

Dependencies:
├── internal/utils/
│   ├── redis.go      # Enhanced Redis pool management
│   ├── kafka.go      # Kafka topic management
│   ├── migrations.go # Database migration management
│   └── health.go     # Health checking utilities
```

### Configuration Loading
The tool automatically loads configuration for each service:
- **AuthPost**: Database migration operations
- **NFP**: Kafka and Redis operations  
- **Newsfeed**: Kafka and Redis operations
- **Web**: Redis operations

## 📈 Performance Monitoring

### Redis Pool Metrics
Monitor these key metrics for Redis performance:
- **Hits/Misses**: Cache effectiveness
- **Total/Idle Connections**: Pool utilization
- **Timeouts**: Connection pressure
- **Stale Connections**: Pool health

### Kafka Health Indicators
- **Topic Existence**: Auto-creation success
- **Partition Count**: Load distribution
- **Connection Success**: Broker availability

### Database Migration Tracking
- **Applied Count**: Schema version
- **Pending Migrations**: Required updates
- **Last Applied**: Current schema state

## 🚀 Production Usage

### Deployment Checklist
1. ✅ Run `make system-migration-status` to check schema state
2. ✅ Verify Kafka topics with `make system-kafka-topics`  
3. ✅ Monitor Redis pools with `make system-redis-status-*`
4. ✅ Run comprehensive health check: `make system-health-check`

### Monitoring Integration
- **Metrics**: Redis pool statistics logged every 5 minutes
- **Health Checks**: All services report enhanced dependency status
- **Alerting**: Monitor pool timeout rates and migration failures

The system administration tool provides enterprise-grade infrastructure management with comprehensive monitoring and automation capabilities. 