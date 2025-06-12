# WanderSphere - Social Media Platform

A modern social media platform built with Go (backend) and React (frontend), featuring real-time communication, file uploads, and social networking capabilities.

## 🚀 Quick Start with Docker Compose

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) (v20.10+)
- [Docker Compose](https://docs.docker.com/compose/install/) (v2.0+)
- At least 4GB RAM available for containers

### Running the Application

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd WanderSphere
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - **Frontend (React App)**: http://localhost:3000
   - **Backend API**: http://localhost:19003
   - **MinIO Console (File Storage)**: http://localhost:9001
   - **PostgreSQL**: localhost:5434 (user: postgres, password: 123456)

### 🛠️ Development Commands

#### Basic Operations
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Rebuild and restart
docker-compose up -d --build

# View running containers
docker-compose ps
```

#### Service-specific Operations
```bash
# Start only core services (no Kafka)
docker-compose up -d frontend backend postgres redis minio aap newsfeed

# Start with Kafka for event streaming
docker-compose --profile kafka up -d

# Restart a specific service
docker-compose restart frontend

# View logs for specific service
docker-compose logs -f backend
```

#### Data Management
```bash
# Reset all data (WARNING: This will delete all data!)
docker-compose down -v

# Backup database
docker-compose exec postgres pg_dump -U postgres wander_sphere > backup.sql

# Restore database
docker-compose exec -T postgres psql -U postgres wander_sphere < backup.sql
```

## 🏗️ Architecture Overview

### Services

| Service | Port | Description |
|---------|------|-------------|
| **frontend** | 3000 | React application served by Nginx |
| **backend** | 19003 | Main Go web API (webapp service) |
| **aap** | 19001, 19101 | Authentication and Posts service |
| **newsfeed** | 19002, 19102 | Newsfeed service |
| **postgres** | 5434 | PostgreSQL database |
| **redis** | 6379 | Redis cache |
| **minio** | 9000, 9001 | S3-compatible object storage |
| **kafka** | 9092 | Event streaming (optional) |
| **zookeeper** | 2181 | Kafka dependency (optional) |

### Network Communication

- All services communicate via the `wandersphere-network` Docker network
- Frontend proxies API requests to the backend through Nginx
- Backend services communicate internally using service names
- External access only through exposed ports

## 🔧 Configuration

### Environment Variables

The application uses the following key environment variables:

```yaml
# Frontend
NODE_ENV=production

# Backend
GO_ENV=production

# Database
POSTGRES_DB=wander_sphere
POSTGRES_USER=postgres
POSTGRES_PASSWORD=123456

# MinIO
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
```

### Backend Configuration

The backend services use `config.yaml` mounted from the host. Key settings:

- Database connection strings
- Redis configuration
- MinIO/S3 settings
- Service discovery endpoints

## 📱 Features

### Frontend (React)
- ✅ **User Authentication** - Login/signup with session management
- ✅ **Post Creation** - Rich text posts with image uploads (up to 4 images)
- ✅ **Social Features** - Follow/unfollow users, view followers/following
- ✅ **Newsfeed** - Personalized content feed
- ✅ **Comments** - Comment on posts with threaded discussions
- ✅ **Profile Management** - Edit profile, upload profile/cover photos
- ✅ **Search** - Discover users and content
- ✅ **Responsive Design** - Mobile-first UI with Tailwind CSS

### Backend (Go)
- ✅ **Microservices Architecture** - Modular service design
- ✅ **RESTful API** - Comprehensive API endpoints
- ✅ **File Upload** - S3-compatible storage with MinIO
- ✅ **Real-time Features** - Event-driven architecture
- ✅ **Caching** - Redis for performance optimization
- ✅ **Database** - PostgreSQL with proper migrations

## 🔒 Security Features

- **CORS Protection** - Configured for cross-origin requests
- **Input Validation** - Comprehensive input sanitization
- **Session Management** - Secure user sessions
- **File Upload Security** - Validation and size limits
- **Database Security** - Prepared statements and validation

## 🚨 Troubleshooting

### Common Issues

1. **Port Conflicts**
   ```bash
   # Check if ports are in use
   netstat -tulpn | grep :3000
   
   # Modify docker-compose.yml ports if needed
   ```

2. **Database Connection Issues**
   ```bash
   # Check database health
   docker-compose exec postgres pg_isready -U postgres
   
   # Reset database
   docker-compose down -v
   docker-compose up postgres -d
   ```

3. **Frontend Build Issues**
   ```bash
   # Rebuild frontend
   docker-compose build --no-cache frontend
   docker-compose up -d frontend
   ```

4. **Memory Issues**
   ```bash
   # Check container memory usage
   docker stats
   
   # Increase Docker memory limit in Docker Desktop settings
   ```

### Logs and Debugging

```bash
# View all logs
docker-compose logs

# Follow logs for specific service
docker-compose logs -f backend

# Access container shell
docker-compose exec backend sh
docker-compose exec frontend sh

# Check service health
docker-compose ps
```

## 🔄 Updates and Maintenance

### Updating the Application
```bash
# Pull latest changes
git pull origin main

# Rebuild and restart services
docker-compose up -d --build

# Clean up old images
docker image prune -f
```

### Database Migrations
```bash
# Run migrations (if applicable)
docker-compose exec backend ./webapp.linux migrate

# Or use the backend's migration tools
cd backend && make migrate-up
```

## 📈 Monitoring

### Health Checks
All services include health checks:
```bash
# Check service health
docker-compose ps

# Manual health check
curl http://localhost:19003/health
curl http://localhost:3000
```

### Resource Usage
```bash
# Monitor resource usage
docker stats

# Check disk usage
docker system df
```

## 🎯 Production Considerations

For production deployment:

1. **Change default passwords** in environment variables
2. **Configure SSL/TLS** certificates
3. **Set up proper backup** strategies
4. **Monitor resource usage** and scale accordingly
5. **Configure log aggregation**
6. **Set up health monitoring**

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test with Docker Compose
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Built with ❤️ using Go, React, PostgreSQL, Redis, and MinIO** 