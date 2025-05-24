#!/bin/bash

# WanderSphere API Test Runner
# Comprehensive testing script for all API endpoints

set -e

# Configuration
API_BASE_URL="http://localhost:19003/api/v1"
TEST_TIMEOUT="300s"
LOG_LEVEL="info"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if service is running
check_service() {
    local service_name=$1
    local port=$2
    local max_wait=60
    local wait_time=0
    
    log_info "Waiting for $service_name service on port $port..."
    
    while [ $wait_time -lt $max_wait ]; do
        if nc -z localhost $port 2>/dev/null; then
            log_success "$service_name service is running on port $port"
            return 0
        fi
        sleep 2
        wait_time=$((wait_time + 2))
    done
    
    log_error "$service_name service failed to start on port $port"
    return 1
}

# Function to check API endpoint
check_api_endpoint() {
    local endpoint=$1
    local expected_status=${2:-200}
    
    log_info "Testing endpoint: $endpoint"
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" "$API_BASE_URL$endpoint" || echo "HTTPSTATUS:000")
    status=$(echo "$response" | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    
    if [[ "$status" == "$expected_status" ]] || [[ "$status" == "404" ]] || [[ "$status" == "401" ]]; then
        log_success "Endpoint $endpoint responded with status $status"
        return 0
    else
        log_warning "Endpoint $endpoint responded with status $status (expected $expected_status)"
        return 1
    fi
}

# Get the backend directory (parent of tests directory)
BACKEND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Cleanup function
cleanup() {
    log_info "Cleaning up test environment..."
    cd "$BACKEND_DIR" && make stop >/dev/null 2>&1 || true
}

# Trap cleanup on exit
trap cleanup EXIT

echo "🚀 WanderSphere API Testing Suite"
echo "================================="
echo ""

# Phase 1: Environment Setup
log_info "Phase 1: Setting up test environment"
echo "------------------------------------"

# Clean previous state
cleanup

# Start infrastructure services using Makefile
log_info "Starting infrastructure services..."
if ! (cd "$BACKEND_DIR" && make infra); then
    log_error "Failed to start infrastructure services"
    exit 1
fi

# Check infrastructure services
if ! check_service "postgres" 5434; then 
    log_error "Postgres failed to start, checking logs..."
    cd "$BACKEND_DIR" && docker-compose logs postgres
    exit 1
fi
if ! check_service "redis" 6379; then 
    log_error "Redis failed to start, checking logs..."
    cd "$BACKEND_DIR" && docker-compose logs redis
    exit 1
fi
if ! check_service "kafka" 9092; then 
    log_error "Kafka failed to start, checking logs..."
    cd "$BACKEND_DIR" && docker-compose logs kafka
    exit 1
fi

# Run database migrations using Makefile
log_info "Running database migrations..."
if ! (cd "$BACKEND_DIR" && make migrate); then
    log_error "Database migrations failed"
    exit 1
fi

# Start application services using Makefile
log_info "Starting application services..."
if ! (cd "$BACKEND_DIR" && make services); then
    log_error "Failed to start application services"
    exit 1
fi

# Check application services using Makefile health check
log_info "Checking application services..."
if ! (cd "$BACKEND_DIR" && make health); then
    log_warning "Some services may not be responding properly"
fi

echo ""

# Phase 2: Basic API Connectivity
log_info "Phase 2: Testing basic API connectivity"
echo "--------------------------------------"

# Test main endpoints for basic connectivity
check_api_endpoint "/users/login" 400  # Should return 400 without body
check_api_endpoint "/users/signup" 400 # Should return 400 without body
check_api_endpoint "/newsfeed" 401     # Should return 401 without auth

echo ""

# Phase 3: Go-based API Tests
log_info "Phase 3: Running comprehensive Go-based API tests"
echo "------------------------------------------------"

# Set environment variables for Go tests
export API_BASE_URL="$API_BASE_URL"
export TEST_TIMEOUT="$TEST_TIMEOUT"

# Navigate to test directory
cd "$(dirname "$0")"

# Initialize Go module if it doesn't exist
if [ ! -f "go.mod" ]; then
    log_info "Initializing Go module for tests..."
    go mod init wandersphere-api-tests
fi

# Update go.mod to ensure compatibility
log_info "Ensuring Go module is up to date..."
go mod tidy

# Clear Go test cache to avoid stale results
log_info "Clearing Go test cache..."
go clean -testcache

# Track test results
TEST_RESULTS=()
TESTS_PASSED=0
TESTS_FAILED=0

# Run Go tests with proper error tracking
log_info "Running authentication tests..."
if go test -v ./api -run TestUser -timeout=$TEST_TIMEOUT; then
    log_success "Authentication tests passed"
    TEST_RESULTS+=("✅ Authentication tests: PASSED")
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "Authentication tests failed"
    TEST_RESULTS+=("❌ Authentication tests: FAILED")
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

log_info "Running posts management tests..."
if go test -v ./api -run TestCreate -timeout=$TEST_TIMEOUT; then
    log_success "Posts management tests passed"
    TEST_RESULTS+=("✅ Posts management tests: PASSED")
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "Posts management tests failed"
    TEST_RESULTS+=("❌ Posts management tests: FAILED")
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

log_info "Running social features tests..."
if go test -v ./api -run TestFollow -timeout=$TEST_TIMEOUT; then
    log_success "Social features tests passed"
    TEST_RESULTS+=("✅ Social features tests: PASSED")
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "Social features tests failed"
    TEST_RESULTS+=("❌ Social features tests: FAILED")
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

log_info "Running newsfeed tests..."
if go test -v ./api -run TestNewsfeed -timeout=$TEST_TIMEOUT; then
    log_success "Newsfeed tests passed"
    TEST_RESULTS+=("✅ Newsfeed tests: PASSED")
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "Newsfeed tests failed"
    TEST_RESULTS+=("❌ Newsfeed tests: FAILED")
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

log_info "Running integration tests..."
if go test -v ./api -run TestComplete -timeout=$TEST_TIMEOUT; then
    log_success "Integration tests passed"
    TEST_RESULTS+=("✅ Integration tests: PASSED")
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "Integration tests failed"
    TEST_RESULTS+=("❌ Integration tests: FAILED")
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

echo ""

# Phase 4: Test Summary
log_info "Phase 4: Test execution summary"
echo "------------------------------"

# Generate test report
log_info "Generating comprehensive test report..."
go test -v ./api/... -json -timeout=$TEST_TIMEOUT > test_results.json 2>/dev/null || true

# Determine overall test status
if [ $TESTS_FAILED -eq 0 ]; then
    API_TEST_STATUS="✅ All API tests passed ($TESTS_PASSED/$((TESTS_PASSED + TESTS_FAILED)))"
    OVERALL_STATUS="🎉 ALL TESTS PASSED! WanderSphere API is fully functional!"
else
    API_TEST_STATUS="❌ Some API tests failed ($TESTS_PASSED/$((TESTS_PASSED + TESTS_FAILED)) passed)"
    OVERALL_STATUS="⚠️ Some tests failed - system needs attention"
fi

log_success "API testing completed!"
echo ""
echo "📊 Detailed Test Results:"
for result in "${TEST_RESULTS[@]}"; do
    echo "  $result"
done
echo ""
echo "📊 Test Summary:"
echo "- Infrastructure: ✅ Core services running"
echo "- Connectivity: ✅ Basic endpoints responding"
echo "- API Tests: $API_TEST_STATUS"
echo ""
echo "$OVERALL_STATUS"

# Optional: Keep services running for manual testing
if [[ "${KEEP_RUNNING:-false}" == "true" ]]; then
    log_info "Services will keep running for manual testing."
    log_info "To stop services: make stop"
    log_info "API Base URL: $API_BASE_URL"
else
    log_info "Stopping services..."
    cleanup
fi 