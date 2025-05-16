package api

import (
	"fmt"
	"testing"
	"time"

	"wandersphere-api-tests/utils"
)

func TestUserSignup(t *testing.T) {
	client, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	authHelper := utils.NewAuthHelper(client)

	t.Run("Valid Signup", func(t *testing.T) {
		username := utils.GenerateUniqueUsername("testuser")
		email := fmt.Sprintf("%s@wandersphere.com", username)

		signupReq := utils.CreateUserRequest{
			UserName:    username,
			Email:       email,
			Password:    "TestPass123!",
			FirstName:   "Test",
			LastName:    "User",
			DateOfBirth: "1990-01-01",
		}

		resp, err := authHelper.SignupUser(signupReq)
		if err != nil {
			t.Fatalf("Valid signup failed: %v", err)
		}

		if resp.Message == "" {
			t.Error("Expected success message in response")
		}

		t.Logf("Signup successful: %s", resp.Message)
	})

	t.Run("Invalid Signup - Missing Required Fields", func(t *testing.T) {
		signupReq := utils.CreateUserRequest{
			UserName: "testuser_invalid",
			// Missing required fields: email, password, date_of_birth
		}

		_, err := authHelper.SignupUser(signupReq)
		if err == nil {
			t.Error("Expected signup to fail with missing required fields")
		}

		t.Logf("Expected failure: %v", err)
	})

	t.Run("Invalid Signup - Duplicate Username", func(t *testing.T) {
		username := utils.GenerateUniqueUsername("duplicate_user")
		email1 := fmt.Sprintf("%s_1@wandersphere.com", username)
		email2 := fmt.Sprintf("%s_2@wandersphere.com", username)

		signupReq := utils.CreateUserRequest{
			UserName:    username,
			Email:       email1,
			Password:    "TestPass123!",
			FirstName:   "Test",
			LastName:    "User",
			DateOfBirth: "1990-01-01",
		}

		// First signup should succeed
		_, err := authHelper.SignupUser(signupReq)
		if err != nil {
			t.Fatalf("First signup failed: %v", err)
		}

		// Second signup with same username should fail
		signupReq.Email = email2
		_, err = authHelper.SignupUser(signupReq)
		if err == nil {
			t.Error("Expected signup to fail with duplicate username")
		}

		t.Logf("Expected duplicate username failure: %v", err)
	})
}

func TestUserLogin(t *testing.T) {
	client, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	authHelper := utils.NewAuthHelper(client)

	// Create a test user first
	testUser, err := authHelper.CreateTestUser("", "", "")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("Valid Login", func(t *testing.T) {
		// Clear existing session
		authHelper.Logout()

		loginReq := utils.LoginRequest{
			UserName: testUser.UserData.UserName,
			Password: "TestPass123!",
		}

		authenticatedUser, err := authHelper.LoginUser(loginReq)
		if err != nil {
			t.Fatalf("Valid login failed: %v", err)
		}

		if authenticatedUser.UserID == 0 {
			t.Error("Expected valid user ID after login")
		}

		if authenticatedUser.SessionID == "" {
			t.Error("Expected session ID after login")
		}

		if authenticatedUser.UserData.UserName != testUser.UserData.UserName {
			t.Error("Returned user data doesn't match expected user")
		}

		t.Logf("Login successful for user: %s (ID: %d)", authenticatedUser.UserData.UserName, authenticatedUser.UserID)
	})

	t.Run("Invalid Login - Wrong Password", func(t *testing.T) {
		// Clear existing session
		authHelper.Logout()

		loginReq := utils.LoginRequest{
			UserName: testUser.UserData.UserName,
			Password: "WrongPassword123!",
		}

		_, err := authHelper.LoginUser(loginReq)
		if err == nil {
			t.Error("Expected login to fail with wrong password")
		}

		t.Logf("Expected failure with wrong password: %v", err)
	})

	t.Run("Invalid Login - Nonexistent User", func(t *testing.T) {
		// Clear existing session
		authHelper.Logout()

		loginReq := utils.LoginRequest{
			UserName: "nonexistent_user_12345",
			Password: "TestPass123!",
		}

		_, err := authHelper.LoginUser(loginReq)
		if err == nil {
			t.Error("Expected login to fail for nonexistent user")
		}

		t.Logf("Expected failure for nonexistent user: %v", err)
	})

	t.Run("Invalid Login - Missing Fields", func(t *testing.T) {
		// Clear existing session
		authHelper.Logout()

		loginReq := utils.LoginRequest{
			UserName: testUser.UserData.UserName,
			// Missing password
		}

		_, err := authHelper.LoginUser(loginReq)
		if err == nil {
			t.Error("Expected login to fail with missing password")
		}

		t.Logf("Expected failure with missing password: %v", err)
	})
}

func TestGetUserDetails(t *testing.T) {
	client, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	authHelper := utils.NewAuthHelper(client)

	// Create a test user
	testUser, err := authHelper.CreateTestUser("", "", "")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("Valid Get User Details", func(t *testing.T) {
		userDetails, err := authHelper.GetUserDetails(testUser.UserID)
		if err != nil {
			t.Fatalf("Failed to get user details: %v", err)
		}

		if userDetails.UserID != testUser.UserID {
			t.Errorf("Expected user ID %d, got %d", testUser.UserID, userDetails.UserID)
		}

		if userDetails.UserName != testUser.UserData.UserName {
			t.Errorf("Expected username %s, got %s", testUser.UserData.UserName, userDetails.UserName)
		}

		if userDetails.Email != testUser.UserData.Email {
			t.Errorf("Expected email %s, got %s", testUser.UserData.Email, userDetails.Email)
		}

		t.Logf("Retrieved user details: %+v", userDetails)
	})

	t.Run("Invalid Get User Details - Nonexistent User", func(t *testing.T) {
		nonexistentUserID := 999999
		_, err := authHelper.GetUserDetails(nonexistentUserID)
		if err == nil {
			t.Error("Expected get user details to fail for nonexistent user")
		}

		t.Logf("Expected failure for nonexistent user: %v", err)
	})

	t.Run("Invalid Get User Details - Invalid User ID", func(t *testing.T) {
		invalidUserID := -1
		_, err := authHelper.GetUserDetails(invalidUserID)
		if err == nil {
			t.Error("Expected get user details to fail for invalid user ID")
		}

		t.Logf("Expected failure for invalid user ID: %v", err)
	})
}

func TestUserEdit(t *testing.T) {
	client, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	authHelper := utils.NewAuthHelper(client)

	// Create and login a test user
	testUser, err := authHelper.CreateTestUser("", "", "")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("Valid Profile Edit", func(t *testing.T) {
		editReq := utils.EditUserRequest{
			FirstName:      "Updated",
			LastName:       "Name",
			ProfilePicture: "https://example.com/new-profile.jpg",
			CoverPicture:   "https://example.com/new-cover.jpg",
		}

		resp, err := authHelper.EditUserProfile(editReq)
		if err != nil {
			t.Fatalf("Profile edit failed: %v", err)
		}

		if resp.Message == "" {
			t.Error("Expected success message in edit response")
		}

		// Verify changes by getting user details
		userDetails, err := authHelper.GetUserDetails(testUser.UserID)
		if err != nil {
			t.Fatalf("Failed to get updated user details: %v", err)
		}

		if userDetails.FirstName != editReq.FirstName {
			t.Errorf("Expected first name %s, got %s", editReq.FirstName, userDetails.FirstName)
		}

		if userDetails.LastName != editReq.LastName {
			t.Errorf("Expected last name %s, got %s", editReq.LastName, userDetails.LastName)
		}

		t.Logf("Profile edit successful: %s", resp.Message)
	})

	t.Run("Edit Without Authentication", func(t *testing.T) {
		// Create a new client without authentication
		unauthClient, err := utils.NewAPIClient()
		if err != nil {
			t.Fatalf("Failed to create unauthenticated client: %v", err)
		}

		unauthHelper := utils.NewAuthHelper(unauthClient)

		editReq := utils.EditUserRequest{
			FirstName: "Should",
			LastName:  "Fail",
		}

		_, err = unauthHelper.EditUserProfile(editReq)
		if err == nil {
			t.Error("Expected profile edit to fail without authentication")
		}

		t.Logf("Expected failure without authentication: %v", err)
	})

	t.Run("Valid Password Change", func(t *testing.T) {
		newPassword := "NewTestPass456!"
		editReq := utils.EditUserRequest{
			Password: newPassword,
		}

		resp, err := authHelper.EditUserProfile(editReq)
		if err != nil {
			t.Fatalf("Password change failed: %v", err)
		}

		if resp.Message == "" {
			t.Error("Expected success message in password change response")
		}

		// Verify new password works by logging in with it
		authHelper.Logout()
		loginReq := utils.LoginRequest{
			UserName: testUser.UserData.UserName,
			Password: newPassword,
		}

		_, err = authHelper.LoginUser(loginReq)
		if err != nil {
			t.Fatalf("Failed to login with new password: %v", err)
		}

		t.Logf("Password change successful: %s", resp.Message)
	})
}

func TestAuthenticationFlow(t *testing.T) {
	client, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	authHelper := utils.NewAuthHelper(client)

	t.Run("Complete Authentication Flow", func(t *testing.T) {
		timestamp := time.Now().UnixNano()
		username := fmt.Sprintf("flowtest_%d_%d", timestamp, time.Now().Nanosecond()%10000)
		email := fmt.Sprintf("flowtest_%d@wandersphere.com", timestamp)
		password := "FlowTest123!"

		// Step 1: Signup
		signupReq := utils.CreateUserRequest{
			UserName:    username,
			Email:       email,
			Password:    password,
			FirstName:   "Flow",
			LastName:    "Test",
			DateOfBirth: "1990-01-01",
		}

		signupResp, err := authHelper.SignupUser(signupReq)
		if err != nil {
			t.Fatalf("Signup failed in flow test: %v", err)
		}
		t.Logf("✓ Signup successful: %s", signupResp.Message)

		// Step 2: Login
		loginReq := utils.LoginRequest{
			UserName: username,
			Password: password,
		}

		user, err := authHelper.LoginUser(loginReq)
		if err != nil {
			t.Fatalf("Login failed in flow test: %v", err)
		}
		t.Logf("✓ Login successful: User ID %d", user.UserID)

		// Step 3: Verify session
		if !authHelper.IsAuthenticated() {
			t.Error("Expected to be authenticated after login")
		}
		t.Logf("✓ Session validated")

		// Step 4: Get user details
		userDetails, err := authHelper.GetUserDetails(user.UserID)
		if err != nil {
			t.Fatalf("Get user details failed in flow test: %v", err)
		}
		t.Logf("✓ User details retrieved: %s", userDetails.UserName)

		// Step 5: Edit profile
		editReq := utils.EditUserRequest{
			FirstName: "Updated Flow",
			LastName:  "Test User",
		}

		editResp, err := authHelper.EditUserProfile(editReq)
		if err != nil {
			t.Fatalf("Profile edit failed in flow test: %v", err)
		}
		t.Logf("✓ Profile edit successful: %s", editResp.Message)

		// Step 6: Verify changes
		updatedDetails, err := authHelper.GetUserDetails(user.UserID)
		if err != nil {
			t.Fatalf("Failed to get updated details in flow test: %v", err)
		}

		if updatedDetails.FirstName != editReq.FirstName {
			t.Errorf("Profile update not persisted: expected %s, got %s", editReq.FirstName, updatedDetails.FirstName)
		}
		t.Logf("✓ Profile changes verified")

		// Step 7: Logout
		err = authHelper.Logout()
		if err != nil {
			t.Fatalf("Logout failed in flow test: %v", err)
		}

		if authHelper.IsAuthenticated() {
			t.Error("Expected to be unauthenticated after logout")
		}
		t.Logf("✓ Logout successful")

		t.Log("🎉 Complete authentication flow test passed!")
	})
}
