#!/bin/bash

# WanderSphere Startup Script
# This script helps you get WanderSphere up and running quickly

set -e

echo "🌍 WanderSphere - Social Media Platform"
echo "======================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running. Please start Docker and try again.${NC}"
    exit 1
fi

# Check if Docker Compose is available
if ! command -v docker-compose > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker Compose is not installed. Please install Docker Compose and try again.${NC}"
    exit 1
fi

echo -e "${BLUE}🐳 Docker is running!${NC}"
echo ""

# Function to show help
show_help() {
    echo "Usage: ./start.sh [OPTION]"
    echo ""
    echo "Options:"
    echo "  start, up       Start all services"
    echo "  stop, down      Stop all services"
    echo "  restart         Restart all services"
    echo "  build           Rebuild all images"
    echo "  logs            Show logs for all services"
    echo "  status          Show status of all services"
    echo "  clean           Stop and remove all containers, networks, and volumes"
    echo "  help            Show this help message"
    echo ""
    echo "Examples:"
    echo "  ./start.sh start    # Start WanderSphere"
    echo "  ./start.sh logs     # View application logs"
    echo "  ./start.sh clean    # Reset everything"
}

# Function to start services
start_services() {
    echo -e "${YELLOW}🚀 Starting WanderSphere services...${NC}"
    echo ""
    
    # Build and start services (including infrastructure components)
    docker-compose --profile infra --profile all up -d --build
    
    echo ""
    echo -e "${GREEN}✅ Services started successfully!${NC}"
    echo ""
    echo -e "${BLUE}📍 Access your application:${NC}"
    echo "   🌐 Frontend:     http://localhost:5008"
    echo "   🔧 Backend API:  http://localhost:19003"
    echo "   💾 MinIO UI:     http://localhost:9001 (admin/minioadmin)"
    echo "   🗄️  Database:     localhost:5434 (postgres/123456)"
    echo ""
    echo -e "${YELLOW}💡 Tip: Run './start.sh logs' to view application logs${NC}"
    echo ""
}

# Function to stop services
stop_services() {
    echo -e "${YELLOW}🛑 Stopping WanderSphere services...${NC}"
    docker-compose --profile infra --profile all down
    echo -e "${GREEN}✅ Services stopped successfully!${NC}"
}

# Function to restart services
restart_services() {
    echo -e "${YELLOW}🔄 Restarting WanderSphere services...${NC}"
    docker-compose --profile infra --profile all restart
    echo -e "${GREEN}✅ Services restarted successfully!${NC}"
}

# Function to rebuild services
build_services() {
    echo -e "${YELLOW}🔨 Rebuilding WanderSphere services...${NC}"
    docker-compose --profile infra --profile all build --no-cache
    echo -e "${GREEN}✅ Services rebuilt successfully!${NC}"
}

# Function to show logs
show_logs() {
    echo -e "${BLUE}📋 WanderSphere Application Logs${NC}"
    echo "   Press Ctrl+C to exit"
    echo ""
    docker-compose --profile infra --profile all logs -f
}

# Function to show status
show_status() {
    echo -e "${BLUE}📊 WanderSphere Service Status${NC}"
    echo ""
    docker-compose --profile infra --profile all ps
    echo ""
}

# Function to clean everything
clean_all() {
    echo -e "${RED}🧹 WARNING: This will remove ALL containers, networks, and volumes!${NC}"
    echo -e "${RED}All data will be lost permanently.${NC}"
    echo ""
    read -p "Are you sure you want to continue? (y/N): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}🗑️  Cleaning up...${NC}"
        docker-compose --profile infra --profile all down -v --remove-orphans
        docker system prune -f
        echo -e "${GREEN}✅ Cleanup completed!${NC}"
    else
        echo -e "${BLUE}❌ Cleanup cancelled.${NC}"
    fi
}

# Main script logic
case "${1:-start}" in
    "start"|"up")
        start_services
        ;;
    "stop"|"down")
        stop_services
        ;;
    "restart")
        restart_services
        ;;
    "build")
        build_services
        ;;
    "logs")
        show_logs
        ;;
    "status")
        show_status
        ;;
    "clean")
        clean_all
        ;;
    "help"|"-h"|"--help")
        show_help
        ;;
    *)
        echo -e "${RED}❌ Unknown option: $1${NC}"
        echo ""
        show_help
        exit 1
        ;;
esac 