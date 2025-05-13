// Package tests provides API tests for the webapp
package tests

import (
	"errors"
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
)

// MockWebappService is a mock implementation of the webapp service
type MockWebappService struct {
	mock.Mock
}

// User Service Methods

// CreateUser mocks the user creation endpoint
func (m *MockWebappService) CreateUser(req types.CreateUserRequest) types.MessageResponse {
	args := m.Called(req)
	return args.Get(0).(types.MessageResponse)
}

// Login mocks the user login endpoint
func (m *MockWebappService) Login(req types.LoginRequest) (types.LoginResponse, error) {
	args := m.Called(req)

	// If the method is expected to return an error
	if args.Get(1) != nil {
		return types.LoginResponse{}, args.Error(1)
	}

	return args.Get(0).(types.LoginResponse), nil
}

// GetUserDetails mocks the get user details endpoint
func (m *MockWebappService) GetUserDetails(userID int64) types.UserDetailInfoResponse {
	args := m.Called(userID)
	return args.Get(0).(types.UserDetailInfoResponse)
}

// EditUser mocks the edit user profile endpoint
func (m *MockWebappService) EditUser(req types.EditUserRequest) types.MessageResponse {
	args := m.Called(req)
	return args.Get(0).(types.MessageResponse)
}

// Post Service Methods

// CreatePost mocks the post creation endpoint
func (m *MockWebappService) CreatePost(req types.CreatePostRequest) types.MessageResponse {
	args := m.Called(req)
	return args.Get(0).(types.MessageResponse)
}

// GetPostDetails mocks the get post details endpoint
func (m *MockWebappService) GetPostDetails(postID int64) types.PostDetailInfoResponse {
	args := m.Called(postID)
	return args.Get(0).(types.PostDetailInfoResponse)
}

// EditPost mocks the edit post endpoint
func (m *MockWebappService) EditPost(postID int64, req types.EditPostRequest) types.MessageResponse {
	args := m.Called(postID, req)
	return args.Get(0).(types.MessageResponse)
}

// DeletePost mocks the delete post endpoint
func (m *MockWebappService) DeletePost(postID int64) types.MessageResponse {
	args := m.Called(postID)
	return args.Get(0).(types.MessageResponse)
}

// CommentOnPost mocks the comment on post endpoint
func (m *MockWebappService) CommentOnPost(postID int64, req types.CreatePostCommentRequest) types.PostDetailInfoResponse {
	args := m.Called(postID, req)
	return args.Get(0).(types.PostDetailInfoResponse)
}

// LikePost mocks the like post endpoint
func (m *MockWebappService) LikePost(postID int64) types.MessageResponse {
	args := m.Called(postID)
	return args.Get(0).(types.MessageResponse)
}

// GetS3PresignedUrl mocks the get S3 presigned URL endpoint
func (m *MockWebappService) GetS3PresignedUrl(req types.GetS3PresignedUrlRequest) types.GetS3PresignedUrlResponse {
	args := m.Called(req)
	return args.Get(0).(types.GetS3PresignedUrlResponse)
}

// Friend Service Methods

// FollowUser mocks the follow user endpoint
func (m *MockWebappService) FollowUser(userID int64) types.MessageResponse {
	args := m.Called(userID)
	return args.Get(0).(types.MessageResponse)
}

// UnfollowUser mocks the unfollow user endpoint
func (m *MockWebappService) UnfollowUser(userID int64) types.MessageResponse {
	args := m.Called(userID)
	return args.Get(0).(types.MessageResponse)
}

// GetUserFollowers mocks the get user followers endpoint
func (m *MockWebappService) GetUserFollowers(userID int64) types.UserFollowerResponse {
	args := m.Called(userID)
	return args.Get(0).(types.UserFollowerResponse)
}

// GetUserFollowings mocks the get user followings endpoint
func (m *MockWebappService) GetUserFollowings(userID int64) types.UserFollowingResponse {
	args := m.Called(userID)
	return args.Get(0).(types.UserFollowingResponse)
}

// GetUserPosts mocks the get user posts endpoint
func (m *MockWebappService) GetUserPosts(userID int64) types.UserPostsResponse {
	args := m.Called(userID)
	return args.Get(0).(types.UserPostsResponse)
}

// Newsfeed Service Methods

// GetNewsfeed mocks the get newsfeed endpoint
func (m *MockWebappService) GetNewsfeed() types.NewsfeedResponse {
	args := m.Called()
	return args.Get(0).(types.NewsfeedResponse)
}

// Error handling helper methods

// MockError represents a custom error for testing
type MockError struct {
	Code    int
	Message string
}

// Error implements the error interface
func (e *MockError) Error() string {
	return fmt.Sprintf("%s (Code: %d)", e.Message, e.Code)
}

// NewMockError creates a new mock error
func NewMockError(message string, code int) error {
	return &MockError{
		Message: message,
		Code:    code,
	}
}

// CreateTestError is a helper function that creates a simple error
func CreateTestError(message string) error {
	return errors.New(message)
}
