# WanderSphere API Tests

This directory contains comprehensive tests for all the WanderSphere API endpoints as defined in the swagger specification.

## Test Coverage

The test suite covers the following API categories:

1. **User Management**
   - User registration
   - User login
   - Get user details
   - Edit user profile

2. **Posts**
   - Create post
   - Get post details
   - Edit post
   - Delete post
   - Comment on post
   - Like post
   - Get S3 presigned URL for uploads

3. **Friend Functionality**
   - Follow user
   - Unfollow user
   - Get user followers
   - Get user followings
   - Get user posts

4. **Newsfeed**
   - Get user's newsfeed

## How to Run Tests

### Run All Tests

```bash
cd /path/to/WanderSphere/backend
go test ./...
```

### Run with Verbose Output

```bash
cd /path/to/WanderSphere/backend
go test -v ./...
```

### Run Specific Tests

```bash
cd /path/to/WanderSphere/backend/internal/app/webapp/tests
go test -v -run TestUserSignup  # Run only the user signup test
```

### Run with Coverage

```bash
cd /path/to/WanderSphere/backend
go test ./... -cover
```

### Generate HTML Coverage Report

```bash
cd /path/to/WanderSphere/backend
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## How the Tests Work

The tests use a mocked service implementation (`MockWebappService`) to simulate the behavior of real API endpoints without requiring an actual database or external services.

Each test follows this pattern:
1. Setup a test router with mock service
2. Define expected response from the mock
3. Create and send HTTP request
4. Assert that the response matches expectations

This approach ensures that:
- API endpoints are correctly routed
- Request parameters are properly parsed
- Response structures match the API specification
- Error handling works as expected

## Adding New Tests

To add tests for new endpoints:

1. Add the endpoint to the `setupTestRouter` function in `api_test.go`
2. Add a mock method in `MockWebappService` in `mock_service_test.go`
3. Create a new test function in `api_test.go` 