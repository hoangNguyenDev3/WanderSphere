package utils

import (
	"fmt"
	"strconv"
	"time"
)

// AuthenticatedUser represents a user with valid authentication
type AuthenticatedUser struct {
	UserID    int
	SessionID string
	UserData  UserDetailInfo
	Client    *APIClient
}

// AuthHelper provides authentication utility functions
type AuthHelper struct {
	client *APIClient
}

// NewAuthHelper creates a new authentication helper
func NewAuthHelper(client *APIClient) *AuthHelper {
	return &AuthHelper{client: client}
}

// SignupUser creates a new user account
func (a *AuthHelper) SignupUser(req CreateUserRequest) (*MessageResponse, error) {
	resp, err := a.client.POST("/users/signup", req)
	if err != nil {
		return nil, fmt.Errorf("signup request failed: %w", err)
	}

	if !resp.IsSuccess() {
		var errorResp ErrorResponse
		if err := resp.ParseJSON(&errorResp); err != nil {
			return nil, fmt.Errorf("signup failed with status %d: %s", resp.StatusCode, resp.GetStringBody())
		}
		return nil, fmt.Errorf("signup failed: %s", errorResp.Message)
	}

	var signupResp MessageResponse
	if err := resp.ParseJSON(&signupResp); err != nil {
		return nil, fmt.Errorf("failed to parse signup response: %w", err)
	}

	return &signupResp, nil
}

// LoginUser authenticates a user and returns authenticated user data
func (a *AuthHelper) LoginUser(req LoginRequest) (*AuthenticatedUser, error) {
	resp, err := a.client.POST("/users/login", req)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}

	if !resp.IsSuccess() {
		var errorResp ErrorResponse
		if err := resp.ParseJSON(&errorResp); err != nil {
			return nil, fmt.Errorf("login failed with status %d: %s", resp.StatusCode, resp.GetStringBody())
		}
		return nil, fmt.Errorf("login failed: %s", errorResp.Message)
	}

	var loginResp LoginResponse
	if err := resp.ParseJSON(&loginResp); err != nil {
		return nil, fmt.Errorf("failed to parse login response: %w", err)
	}

	sessionID := a.client.GetSessionID()
	if sessionID == "" {
		return nil, fmt.Errorf("no session cookie received after login")
	}

	return &AuthenticatedUser{
		UserID:    loginResp.User.UserID,
		SessionID: sessionID,
		UserData:  loginResp.User,
		Client:    a.client,
	}, nil
}

// CreateTestUser creates a user for testing purposes
func (a *AuthHelper) CreateTestUser(username, email, password string) (*AuthenticatedUser, error) {
	// Generate unique username with better entropy if not provided
	if username == "" {
		username = GenerateUniqueUsername("testuser")
	}
	if email == "" {
		email = fmt.Sprintf("%s@test.wandersphere.com", username)
	}
	if password == "" {
		password = "TestPass123!"
	}

	// Retry mechanism for username conflicts
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Create user
		signupReq := CreateUserRequest{
			UserName:    username,
			Email:       email,
			Password:    password,
			FirstName:   "Test",
			LastName:    "User",
			DateOfBirth: "1990-01-01",
		}

		_, err := a.SignupUser(signupReq)
		if err != nil {
			// If username already exists, try with a different one
			if attempt < maxRetries-1 {
				username = GenerateUniqueUsername(fmt.Sprintf("testuser_retry_%d", attempt))
				email = fmt.Sprintf("%s@test.wandersphere.com", username)
				continue
			}
			return nil, fmt.Errorf("failed to create test user after %d attempts: %w", maxRetries, err)
		}

		// Login the user
		loginReq := LoginRequest{
			UserName: username,
			Password: password,
		}

		return a.LoginUser(loginReq)
	}

	return nil, fmt.Errorf("failed to create test user after %d attempts", maxRetries)
}

// Logout logs out the current user by clearing cookies
func (a *AuthHelper) Logout() error {
	return a.client.ClearCookies()
}

// IsAuthenticated checks if the client has a valid session
func (a *AuthHelper) IsAuthenticated() bool {
	return a.client.HasValidSession()
}

// GetUserDetails retrieves user details by user ID
func (a *AuthHelper) GetUserDetails(userID int) (*UserDetailInfoResponse, error) {
	resp, err := a.client.GET(fmt.Sprintf("/users/%d", userID))
	if err != nil {
		return nil, fmt.Errorf("get user details request failed: %w", err)
	}

	if !resp.IsSuccess() {
		var errorResp ErrorResponse
		if err := resp.ParseJSON(&errorResp); err != nil {
			return nil, fmt.Errorf("get user details failed with status %d: %s", resp.StatusCode, resp.GetStringBody())
		}
		return nil, fmt.Errorf("get user details failed: %s", errorResp.Message)
	}

	var userResp UserDetailInfoResponse
	if err := resp.ParseJSON(&userResp); err != nil {
		return nil, fmt.Errorf("failed to parse user details response: %w", err)
	}

	return &userResp, nil
}

// EditUserProfile updates user profile information
func (a *AuthHelper) EditUserProfile(req EditUserRequest) (*MessageResponse, error) {
	resp, err := a.client.PUT("/users/edit", req)
	if err != nil {
		return nil, fmt.Errorf("edit user request failed: %w", err)
	}

	if !resp.IsSuccess() {
		var errorResp ErrorResponse
		if err := resp.ParseJSON(&errorResp); err != nil {
			return nil, fmt.Errorf("edit user failed with status %d: %s", resp.StatusCode, resp.GetStringBody())
		}
		return nil, fmt.Errorf("edit user failed: %s", errorResp.Message)
	}

	var editResp MessageResponse
	if err := resp.ParseJSON(&editResp); err != nil {
		return nil, fmt.Errorf("failed to parse edit user response: %w", err)
	}

	return &editResp, nil
}

// CreateMultipleTestUsers creates multiple test users for testing scenarios
func (a *AuthHelper) CreateMultipleTestUsers(count int) ([]*AuthenticatedUser, error) {
	users := make([]*AuthenticatedUser, 0, count)

	for i := 0; i < count; i++ {
		// Generate more unique usernames for multiple users
		username := GenerateUniqueUsername(fmt.Sprintf("testuser_multi_%d", i))

		user, err := a.CreateTestUser(username, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to create test user %d: %w", i, err)
		}
		users = append(users, user)

		// Small delay to ensure different timestamps
		time.Sleep(1 * time.Millisecond)
	}

	return users, nil
}

// AuthenticatedRequest makes an authenticated request (requires session)
func (u *AuthenticatedUser) AuthenticatedRequest(method, path string, body interface{}) (*APIResponse, error) {
	if !u.Client.HasValidSession() {
		return nil, fmt.Errorf("user is not authenticated")
	}
	return u.Client.Request(method, path, body)
}

// GET makes an authenticated GET request
func (u *AuthenticatedUser) GET(path string) (*APIResponse, error) {
	return u.AuthenticatedRequest("GET", path, nil)
}

// POST makes an authenticated POST request
func (u *AuthenticatedUser) POST(path string, body interface{}) (*APIResponse, error) {
	return u.AuthenticatedRequest("POST", path, body)
}

// PUT makes an authenticated PUT request
func (u *AuthenticatedUser) PUT(path string, body interface{}) (*APIResponse, error) {
	return u.AuthenticatedRequest("PUT", path, body)
}

// DELETE makes an authenticated DELETE request
func (u *AuthenticatedUser) DELETE(path string) (*APIResponse, error) {
	return u.AuthenticatedRequest("DELETE", path, nil)
}

// GetUserIDStr returns the user ID as a string for URL path parameters
func (u *AuthenticatedUser) GetUserIDStr() string {
	return strconv.Itoa(u.UserID)
}
