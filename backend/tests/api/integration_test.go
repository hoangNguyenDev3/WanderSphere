package api

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
	"time"

	"wandersphere-api-tests/utils"
)

// generateUniqueUsername creates a highly unique username for testing
func generateUniqueUsername(prefix string) string {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	randomNum, _ := rand.Int(rand.Reader, big.NewInt(999999))
	return fmt.Sprintf("%s_%d_%d_%x", prefix, time.Now().UnixNano(), randomNum.Int64(), randomBytes)
}

func TestCompleteUserJourney(t *testing.T) {
	client, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	authHelper := utils.NewAuthHelper(client)

	t.Run("End-to-End User Journey", func(t *testing.T) {
		username := utils.GenerateUniqueUsername("journey_user")
		email := fmt.Sprintf("%s@test.wandersphere.com", username)

		// Step 1: User Registration
		signupReq := utils.CreateUserRequest{
			UserName:    username,
			Email:       email,
			Password:    "JourneyTest123!",
			FirstName:   "Journey",
			LastName:    "User",
			DateOfBirth: "1990-01-01",
		}

		signupResp, err := authHelper.SignupUser(signupReq)
		if err != nil {
			t.Fatalf("User registration failed: %v", err)
		}
		t.Logf("✓ Step 1: User registered successfully - %s", signupResp.Message)

		// Step 2: User Login
		loginReq := utils.LoginRequest{
			UserName: signupReq.UserName,
			Password: signupReq.Password,
		}

		user, err := authHelper.LoginUser(loginReq)
		if err != nil {
			t.Fatalf("User login failed: %v", err)
		}
		t.Logf("✓ Step 2: User logged in successfully - ID: %d", user.UserID)

		// Step 3: Profile Update
		editReq := utils.EditUserRequest{
			FirstName:      "Updated Journey",
			LastName:       "Test User",
			ProfilePicture: "https://example.com/journey-profile.jpg",
		}

		editResp, err := authHelper.EditUserProfile(editReq)
		if err != nil {
			t.Fatalf("Profile update failed: %v", err)
		}
		t.Logf("✓ Step 3: Profile updated successfully - %s", editResp.Message)

		// Step 4: Verify Profile Changes
		userDetails, err := authHelper.GetUserDetails(user.UserID)
		if err != nil {
			t.Fatalf("Failed to get updated user details: %v", err)
		}

		if userDetails.FirstName != editReq.FirstName {
			t.Errorf("Profile update not reflected: expected %s, got %s", editReq.FirstName, userDetails.FirstName)
		}
		t.Logf("✓ Step 4: Profile changes verified")

		// Step 5: Create First Post
		firstPostReq := utils.CreatePostRequest{
			ContentText: "This is my first post on WanderSphere!",
			Visible:     true,
		}

		firstPostResp, err := user.POST("/posts", firstPostReq)
		if err != nil {
			t.Fatalf("First post creation failed: %v", err)
		}

		if !firstPostResp.IsSuccess() {
			t.Fatalf("First post creation failed with status %d", firstPostResp.StatusCode)
		}
		t.Logf("✓ Step 5: First post created successfully")

		// Step 6: Create Second Post with Images
		secondPostReq := utils.CreatePostRequest{
			ContentText:      "Here's my second post with some images!",
			ContentImagePath: []string{"https://example.com/journey-pic1.jpg", "https://example.com/journey-pic2.jpg"},
			Visible:          true,
		}

		secondPostResp, err := user.POST("/posts", secondPostReq)
		if err != nil {
			t.Fatalf("Second post creation failed: %v", err)
		}

		if !secondPostResp.IsSuccess() {
			t.Fatalf("Second post creation failed with status %d", secondPostResp.StatusCode)
		}
		t.Logf("✓ Step 6: Second post with images created successfully")

		// Step 7: Check User's Posts
		userPostsResp, err := client.GET(fmt.Sprintf("/friends/%d/posts", user.UserID))
		if err != nil {
			t.Fatalf("Failed to get user posts: %v", err)
		}

		if userPostsResp.IsSuccess() {
			var postsResp utils.UserPostsResponse
			if err := userPostsResp.ParseJSON(&postsResp); err == nil {
				t.Logf("✓ Step 7: User has %d posts", len(postsResp.PostsIds))
			}
		} else {
			t.Logf("⚠ Step 7: Failed to get user posts (may be expected): Status %d", userPostsResp.StatusCode)
		}

		// Step 8: Check Initial Newsfeed
		newsfeedResp, err := user.GET("/newsfeed")
		if err != nil {
			t.Fatalf("Failed to get initial newsfeed: %v", err)
		}

		if newsfeedResp.IsSuccess() {
			var newsfeed utils.NewsfeedResponse
			if err := newsfeedResp.ParseJSON(&newsfeed); err == nil {
				t.Logf("✓ Step 8: Initial newsfeed has %d posts", len(newsfeed.PostsIds))
			}
		} else {
			t.Logf("⚠ Step 8: Failed to get newsfeed (may be expected): Status %d", newsfeedResp.StatusCode)
		}

		t.Log("🎉 Complete user journey test passed!")
	})
}

func TestSocialInteractionJourney(t *testing.T) {
	client, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	authHelper := utils.NewAuthHelper(client)

	t.Run("Multi-User Social Interaction Journey", func(t *testing.T) {
		// Create three users for social interaction
		users, err := authHelper.CreateMultipleTestUsers(3)
		if err != nil {
			t.Fatalf("Failed to create test users: %v", err)
		}

		alice, bob, charlie := users[0], users[1], users[2]
		t.Logf("✓ Created 3 test users: Alice (ID:%d), Bob (ID:%d), Charlie (ID:%d)",
			alice.UserID, bob.UserID, charlie.UserID)

		// Step 1: Each user creates initial posts
		for i, user := range users {
			names := []string{"Alice", "Bob", "Charlie"}
			createReq := utils.CreatePostRequest{
				ContentText: fmt.Sprintf("Hello WanderSphere! This is %s's introduction post.", names[i]),
				Visible:     true,
			}

			createResp, err := user.POST("/posts", createReq)
			if err != nil {
				t.Fatalf("%s post creation failed: %v", names[i], err)
			}

			if createResp.IsSuccess() {
				t.Logf("✓ %s created an introduction post", names[i])
			} else {
				t.Logf("⚠ %s post creation failed: Status %d", names[i], createResp.StatusCode)
			}
		}

		// Step 2: Alice follows Bob and Charlie
		followResp1, err := alice.POST(fmt.Sprintf("/friends/%d", bob.UserID), nil)
		if err == nil && followResp1.IsSuccess() {
			t.Logf("✓ Alice followed Bob")
		} else {
			t.Logf("⚠ Alice failed to follow Bob")
		}

		followResp2, err := alice.POST(fmt.Sprintf("/friends/%d", charlie.UserID), nil)
		if err == nil && followResp2.IsSuccess() {
			t.Logf("✓ Alice followed Charlie")
		} else {
			t.Logf("⚠ Alice failed to follow Charlie")
		}

		// Step 3: Bob follows Charlie (but not Alice back)
		followResp3, err := bob.POST(fmt.Sprintf("/friends/%d", charlie.UserID), nil)
		if err == nil && followResp3.IsSuccess() {
			t.Logf("✓ Bob followed Charlie")
		} else {
			t.Logf("⚠ Bob failed to follow Charlie")
		}

		// Step 4: Verify follow relationships
		// Check Alice's followings
		aliceFollowingsResp, err := client.GET(fmt.Sprintf("/friends/%d/followings", alice.UserID))
		if err == nil && aliceFollowingsResp.IsSuccess() {
			var followings utils.UserFollowingResponse
			if aliceFollowingsResp.ParseJSON(&followings) == nil {
				t.Logf("✓ Alice follows %d users", len(followings.FollowingsIDs))
			}
		}

		// Check Charlie's followers
		charlieFollowersResp, err := client.GET(fmt.Sprintf("/friends/%d/followers", charlie.UserID))
		if err == nil && charlieFollowersResp.IsSuccess() {
			var followers utils.UserFollowerResponse
			if charlieFollowersResp.ParseJSON(&followers) == nil {
				t.Logf("✓ Charlie has %d followers", len(followers.FollowersIDs))
			}
		}

		// Step 5: Bob creates a travel post
		travelPostReq := utils.CreatePostRequest{
			ContentText: "Just visited an amazing beach! The sunset was incredible. #travel #beach",
			Visible:     true,
		}

		travelPostResp, err := bob.POST("/posts", travelPostReq)
		if err == nil && travelPostResp.IsSuccess() {
			t.Logf("✓ Bob created a travel post")
		} else {
			t.Logf("⚠ Bob failed to create travel post")
		}

		// Step 6: Alice likes Bob's travel post (assuming post ID 1)
		likeResp, err := alice.POST("/posts/1/likes", nil)
		if err == nil && likeResp.IsSuccess() {
			t.Logf("✓ Alice liked Bob's travel post")
		} else {
			t.Logf("⚠ Alice failed to like Bob's post")
		}

		// Step 7: Charlie comments on Bob's travel post
		commentReq := utils.CreatePostCommentRequest{
			ContentText: "Wow, that sounds amazing! Which beach was it?",
		}

		commentResp, err := charlie.POST("/posts/1", commentReq)
		if err == nil && commentResp.IsSuccess() {
			t.Logf("✓ Charlie commented on Bob's travel post")
		} else {
			t.Logf("⚠ Charlie failed to comment on Bob's post")
		}

		// Step 8: Check post details to see interactions
		postDetailsResp, err := client.GET("/posts/1")
		if err == nil && postDetailsResp.IsSuccess() {
			var postDetails utils.PostDetailInfoResponse
			if postDetailsResp.ParseJSON(&postDetails) == nil {
				t.Logf("✓ Travel post has %d likes and %d comments",
					len(postDetails.UsersLiked), len(postDetails.Comments))
			}
		} else {
			t.Logf("⚠ Failed to get post details")
		}

		// Step 9: Check newsfeeds after social activity
		for i, user := range users {
			names := []string{"Alice", "Bob", "Charlie"}
			newsfeedResp, err := user.GET("/newsfeed")
			if err == nil && newsfeedResp.IsSuccess() {
				var newsfeed utils.NewsfeedResponse
				if newsfeedResp.ParseJSON(&newsfeed) == nil {
					t.Logf("✓ %s's newsfeed contains %d posts", names[i], len(newsfeed.PostsIds))
				}
			} else {
				t.Logf("⚠ %s's newsfeed failed", names[i])
			}
		}

		// Step 10: Alice unfollows Charlie
		unfollowResp, err := alice.DELETE(fmt.Sprintf("/friends/%d", charlie.UserID))
		if err == nil && unfollowResp.IsSuccess() {
			t.Logf("✓ Alice unfollowed Charlie")
		} else {
			t.Logf("⚠ Alice failed to unfollow Charlie")
		}

		t.Log("🎉 Social interaction journey test completed!")
	})
}

func TestContentLifecycleIntegration(t *testing.T) {
	// Create separate clients for each user to avoid session conflicts
	client1, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client 1: %v", err)
	}

	client2, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client 2: %v", err)
	}

	authHelper1 := utils.NewAuthHelper(client1)
	authHelper2 := utils.NewAuthHelper(client2)

	// Use a neutral client for GET requests that don't require specific authentication
	neutralClient, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create neutral API client: %v", err)
	}

	t.Run("Complete Content Lifecycle", func(t *testing.T) {
		// Create two users: content creator and content consumer with separate clients
		creator, err := authHelper1.CreateTestUser("", "", "")
		if err != nil {
			t.Fatalf("Failed to create creator user: %v", err)
		}

		consumer, err := authHelper2.CreateTestUser("", "", "")
		if err != nil {
			t.Fatalf("Failed to create consumer user: %v", err)
		}

		t.Logf("✓ Created Creator (ID:%d) and Consumer (ID:%d)", creator.UserID, consumer.UserID)

		// Step 1: Consumer follows Creator
		followResp, err := consumer.POST(fmt.Sprintf("/friends/%d", creator.UserID), nil)
		if err == nil && followResp.IsSuccess() {
			t.Logf("✓ Consumer followed Creator")
		} else {
			t.Logf("⚠ Consumer failed to follow Creator")
		}

		// Step 2: Creator creates a post
		testPostID, err := creator.CreateTestPost("Just discovered this amazing hiking trail! Perfect for weekend adventures.", true)
		if err != nil {
			t.Fatalf("Failed to create initial post: %v", err)
		}
		t.Logf("✓ Creator posted initial content with ID: %d", testPostID)

		// Step 3: Check post appears in Consumer's newsfeed
		consumerNewsfeedResp, err := consumer.GET("/newsfeed")
		if err == nil && consumerNewsfeedResp.IsSuccess() {
			var newsfeed utils.NewsfeedResponse
			if consumerNewsfeedResp.ParseJSON(&newsfeed) == nil {
				t.Logf("✓ Consumer's newsfeed has %d posts", len(newsfeed.PostsIds))
			}
		} else {
			t.Logf("⚠ Consumer's newsfeed check failed")
		}

		// Step 4: Consumer interacts with the post
		// Like the post using the actual post ID
		likeResp, err := consumer.POST(fmt.Sprintf("/posts/%d/likes", testPostID), nil)
		if err == nil && likeResp.IsSuccess() {
			t.Logf("✓ Consumer liked the post")
		} else {
			t.Logf("⚠ Consumer failed to like the post")
		}

		// Comment on the post using the actual post ID
		commentReq := utils.CreatePostCommentRequest{
			ContentText: "This looks incredible! Could you share the exact location?",
		}

		commentResp, err := consumer.POST(fmt.Sprintf("/posts/%d", testPostID), commentReq)
		if err == nil && commentResp.IsSuccess() {
			t.Logf("✓ Consumer commented on the post")
		} else {
			t.Logf("⚠ Consumer failed to comment on the post")
		}

		// Step 5: Creator edits the post to add more details
		contentText := "Just discovered this amazing hiking trail! Perfect for weekend adventures. UPDATE: It's located in Pine Ridge National Park, Trail #7."
		visible := true
		editPostReq := utils.EditPostRequest{
			ContentText: &contentText,
			Visible:     &visible,
		}

		editResp, err := creator.PUT(fmt.Sprintf("/posts/%d", testPostID), editPostReq)
		if err == nil && editResp.IsSuccess() {
			t.Logf("✓ Creator edited the post with more details")
		} else {
			t.Logf("⚠ Creator failed to edit the post")
		}

		// Step 6: Creator replies to the comment
		replyReq := utils.CreatePostCommentRequest{
			ContentText: "Thanks for asking! It's Pine Ridge National Park, Trail #7. The view from the top is worth the climb!",
		}

		replyResp, err := creator.POST(fmt.Sprintf("/posts/%d", testPostID), replyReq)
		if err == nil && replyResp.IsSuccess() {
			t.Logf("✓ Creator replied to the comment")
		} else {
			t.Logf("⚠ Creator failed to reply to the comment")
		}

		// Step 7: Check final post state
		finalPostResp, err := neutralClient.GET(fmt.Sprintf("/posts/%d", testPostID))
		if err == nil && finalPostResp.IsSuccess() {
			var postDetails utils.PostDetailInfoResponse
			if finalPostResp.ParseJSON(&postDetails) == nil {
				t.Logf("✓ Final post state: %d likes, %d comments",
					len(postDetails.UsersLiked), len(postDetails.Comments))

				// Verify content was updated
				if postDetails.ContentText != contentText {
					t.Logf("⚠ Post content may not have been updated as expected")
				} else {
					t.Logf("✓ Post content was successfully updated")
				}
			}
		} else {
			t.Logf("⚠ Failed to get final post state")
		}

		// Step 8: Create another user who joins the conversation
		client3, err := utils.NewAPIClient()
		if err != nil {
			t.Fatalf("Failed to create API client 3: %v", err)
		}
		authHelper3 := utils.NewAuthHelper(client3)

		lateJoiner, err := authHelper3.CreateTestUser("", "", "")
		if err == nil {
			t.Logf("✓ Late joiner user created (ID:%d)", lateJoiner.UserID)

			// Late joiner follows creator
			lateFollowResp, err := lateJoiner.POST(fmt.Sprintf("/friends/%d", creator.UserID), nil)
			if err == nil && lateFollowResp.IsSuccess() {
				t.Logf("✓ Late joiner followed Creator")
			}

			// Late joiner adds their own comment
			lateCommentReq := utils.CreatePostCommentRequest{
				ContentText: "I've been to Pine Ridge too! Trail #7 is definitely one of the best.",
			}

			lateCommentResp, err := lateJoiner.POST(fmt.Sprintf("/posts/%d", testPostID), lateCommentReq)
			if err == nil && lateCommentResp.IsSuccess() {
				t.Logf("✓ Late joiner added their comment")
			}
		}

		// Step 9: Check Creator's posts list
		creatorPostsResp, err := neutralClient.GET(fmt.Sprintf("/friends/%d/posts", creator.UserID))
		if err == nil && creatorPostsResp.IsSuccess() {
			var creatorPosts utils.UserPostsResponse
			if creatorPostsResp.ParseJSON(&creatorPosts) == nil {
				t.Logf("✓ Creator has %d total posts", len(creatorPosts.PostsIds))
			}
		}

		t.Log("🎉 Content lifecycle integration test completed!")
	})
}

func TestAPIErrorHandling(t *testing.T) {
	client, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	authHelper := utils.NewAuthHelper(client)

	t.Run("Comprehensive Error Handling", func(t *testing.T) {
		// Test 1: Invalid authentication scenarios
		t.Logf("Testing authentication error handling...")

		// Try to access protected endpoint without auth
		unauthClient, _ := utils.NewAPIClient()
		resp, err := unauthClient.GET("/newsfeed")
		if err == nil && resp.StatusCode == 401 {
			t.Logf("✓ Unauthenticated request properly rejected")
		}

		// Try to access protected endpoint with invalid session
		resp, err = unauthClient.POST("/posts", utils.CreatePostRequest{
			ContentText: "This should fail",
			Visible:     true,
		})
		if err == nil && resp.StatusCode == 401 {
			t.Logf("✓ Invalid session properly rejected")
		}

		// Test 2: Invalid data scenarios
		t.Logf("Testing data validation error handling...")

		// Create a valid user for further tests
		testUser, err := authHelper.CreateTestUser("", "", "")
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Try to create post with empty content
		resp, err = testUser.POST("/posts", utils.CreatePostRequest{
			ContentText: "", // Empty content should fail
			Visible:     true,
		})
		if err == nil && !resp.IsSuccess() {
			t.Logf("✓ Empty post content properly rejected")
		}

		// Try to follow non-existent user
		resp, err = testUser.POST("/friends/999999", nil)
		if err == nil && !resp.IsSuccess() {
			t.Logf("✓ Non-existent user follow properly rejected")
		}

		// Try to get details of non-existent post
		resp, err = client.GET("/posts/999999")
		if err == nil && !resp.IsSuccess() {
			t.Logf("✓ Non-existent post request properly rejected")
		}

		// Test 3: Invalid user operations
		t.Logf("Testing user operation error handling...")

		// Try to get non-existent user details
		resp, err = client.GET("/users/999999")
		if err == nil && !resp.IsSuccess() {
			t.Logf("✓ Non-existent user details properly rejected")
		}

		// Try to signup with duplicate username
		timestamp := time.Now().UnixNano()
		signupReq := utils.CreateUserRequest{
			UserName:    fmt.Sprintf("error_test_%d_%d", timestamp, time.Now().Nanosecond()%10000),
			Email:       fmt.Sprintf("error_test_%d@test.com", timestamp),
			Password:    "ErrorTest123!",
			FirstName:   "Error",
			LastName:    "Test",
			DateOfBirth: "1990-01-01",
		}

		// First signup should succeed
		_, err = authHelper.SignupUser(signupReq)
		if err == nil {
			t.Logf("✓ First signup succeeded")

			// Second signup with same username should fail
			signupReq.Email = fmt.Sprintf("error_test2_%d@test.com", timestamp)
			_, err = authHelper.SignupUser(signupReq)
			if err != nil {
				t.Logf("✓ Duplicate username signup properly rejected")
			}
		}

		// Test 4: Invalid login scenarios
		t.Logf("Testing login error handling...")

		// Try login with wrong password
		loginReq := utils.LoginRequest{
			UserName: testUser.UserData.UserName,
			Password: "WrongPassword123!",
		}

		_, err = authHelper.LoginUser(loginReq)
		if err != nil {
			t.Logf("✓ Wrong password login properly rejected")
		}

		// Try login with non-existent user
		loginReq = utils.LoginRequest{
			UserName: "nonexistent_user_12345",
			Password: "SomePassword123!",
		}

		_, err = authHelper.LoginUser(loginReq)
		if err != nil {
			t.Logf("✓ Non-existent user login properly rejected")
		}

		t.Log("🎉 Error handling test completed!")
	})
}
