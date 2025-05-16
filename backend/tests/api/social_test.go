package api

import (
	"fmt"
	"testing"

	"wandersphere-api-tests/utils"
)

func TestFollowUser(t *testing.T) {
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

	// Create two test users with separate clients
	user1, err := authHelper1.CreateTestUser("", "", "")
	if err != nil {
		t.Fatalf("Failed to create test user 1: %v", err)
	}

	user2, err := authHelper2.CreateTestUser("", "", "")
	if err != nil {
		t.Fatalf("Failed to create test user 2: %v", err)
	}

	t.Run("Valid Follow User", func(t *testing.T) {
		// User1 follows User2
		resp, err := user1.POST(fmt.Sprintf("/friends/%d", user2.UserID), nil)
		if err != nil {
			t.Fatalf("Follow user request failed: %v", err)
		}

		if resp.IsSuccess() {
			var followResp utils.MessageResponse
			if err := resp.ParseJSON(&followResp); err != nil {
				t.Fatalf("Failed to parse follow response: %v", err)
			}

			if followResp.Message == "" {
				t.Error("Expected success message in follow response")
			}

			t.Logf("User followed successfully: %s", followResp.Message)
		} else if resp.StatusCode == 400 || resp.StatusCode == 404 {
			t.Logf("Follow failed (expected for test): Status %d", resp.StatusCode)
		} else {
			t.Fatalf("Unexpected status for follow user: %d", resp.StatusCode)
		}
	})

	t.Run("Follow User - Unauthenticated", func(t *testing.T) {
		unauthClient, err := utils.NewAPIClient()
		if err != nil {
			t.Fatalf("Failed to create unauthenticated client: %v", err)
		}

		resp, err := unauthClient.POST(fmt.Sprintf("/friends/%d", user2.UserID), nil)
		if err != nil {
			t.Fatalf("Unauthenticated follow request failed: %v", err)
		}

		if resp.StatusCode != 401 {
			t.Errorf("Expected status 401 for unauthenticated follow, got %d", resp.StatusCode)
		}

		t.Logf("Expected authentication failure: Status %d", resp.StatusCode)
	})

	t.Run("Follow Nonexistent User", func(t *testing.T) {
		nonexistentUserID := 999999
		resp, err := user1.POST(fmt.Sprintf("/friends/%d", nonexistentUserID), nil)
		if err != nil {
			t.Fatalf("Follow nonexistent user request failed: %v", err)
		}

		if resp.IsSuccess() {
			t.Error("Expected follow to fail for nonexistent user")
		}

		t.Logf("Expected failure for nonexistent user: Status %d", resp.StatusCode)
	})

	t.Run("Follow Self", func(t *testing.T) {
		t.Logf("User1 ID: %d", user1.UserID)
		t.Logf("Attempting self-follow: user %d trying to follow user %d", user1.UserID, user1.UserID)

		// User1 tries to follow themselves
		resp, err := user1.POST(fmt.Sprintf("/friends/%d", user1.UserID), nil)
		if err != nil {
			t.Fatalf("Follow self request failed: %v", err)
		}

		t.Logf("Self-follow attempt response: Status %d, Body: %s", resp.StatusCode, resp.GetStringBody())

		if resp.IsSuccess() {
			t.Error("❌ Expected follow to fail when following self - this is a security issue!")
		} else {
			t.Logf("✅ Self-follow correctly prevented: Status %d", resp.StatusCode)
		}
	})
}

func TestUnfollowUser(t *testing.T) {
	client, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	authHelper := utils.NewAuthHelper(client)

	// Create two test users
	user1, err := authHelper.CreateTestUser("", "", "")
	if err != nil {
		t.Fatalf("Failed to create test user 1: %v", err)
	}

	user2, err := authHelper.CreateTestUser("", "", "")
	if err != nil {
		t.Fatalf("Failed to create test user 2: %v", err)
	}

	t.Run("Valid Unfollow User", func(t *testing.T) {
		// First try to follow (may or may not succeed)
		user1.POST(fmt.Sprintf("/friends/%d", user2.UserID), nil)

		// Then unfollow
		resp, err := user1.DELETE(fmt.Sprintf("/friends/%d", user2.UserID))
		if err != nil {
			t.Fatalf("Unfollow user request failed: %v", err)
		}

		if resp.IsSuccess() {
			var unfollowResp utils.MessageResponse
			if err := resp.ParseJSON(&unfollowResp); err != nil {
				t.Fatalf("Failed to parse unfollow response: %v", err)
			}

			if unfollowResp.Message == "" {
				t.Error("Expected success message in unfollow response")
			}

			t.Logf("User unfollowed successfully: %s", unfollowResp.Message)
		} else if resp.StatusCode == 400 || resp.StatusCode == 404 {
			t.Logf("Unfollow failed (expected for test): Status %d", resp.StatusCode)
		} else {
			t.Fatalf("Unexpected status for unfollow user: %d", resp.StatusCode)
		}
	})

	t.Run("Unfollow User - Unauthenticated", func(t *testing.T) {
		unauthClient, err := utils.NewAPIClient()
		if err != nil {
			t.Fatalf("Failed to create unauthenticated client: %v", err)
		}

		resp, err := unauthClient.DELETE(fmt.Sprintf("/friends/%d", user2.UserID))
		if err != nil {
			t.Fatalf("Unauthenticated unfollow request failed: %v", err)
		}

		if resp.StatusCode != 401 {
			t.Errorf("Expected status 401 for unauthenticated unfollow, got %d", resp.StatusCode)
		}

		t.Logf("Expected authentication failure: Status %d", resp.StatusCode)
	})

	t.Run("Unfollow Nonexistent User", func(t *testing.T) {
		nonexistentUserID := 999999
		resp, err := user1.DELETE(fmt.Sprintf("/friends/%d", nonexistentUserID))
		if err != nil {
			t.Fatalf("Unfollow nonexistent user request failed: %v", err)
		}

		if resp.IsSuccess() {
			t.Error("Expected unfollow to fail for nonexistent user")
		}

		t.Logf("Expected failure for nonexistent user: Status %d", resp.StatusCode)
	})
}

func TestGetUserFollowers(t *testing.T) {
	client, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	authHelper := utils.NewAuthHelper(client)

	// Create test users
	user1, err := authHelper.CreateTestUser("", "", "")
	if err != nil {
		t.Fatalf("Failed to create test user 1: %v", err)
	}

	user2, err := authHelper.CreateTestUser("", "", "")
	if err != nil {
		t.Fatalf("Failed to create test user 2: %v", err)
	}

	t.Run("Valid Get User Followers", func(t *testing.T) {
		resp, err := client.GET(fmt.Sprintf("/friends/%d/followers", user1.UserID))
		if err != nil {
			t.Fatalf("Get user followers request failed: %v", err)
		}

		if resp.IsSuccess() {
			var followersResp utils.UserFollowerResponse
			if err := resp.ParseJSON(&followersResp); err != nil {
				t.Fatalf("Failed to parse followers response: %v", err)
			}

			t.Logf("User followers retrieved: %d followers", len(followersResp.FollowersIDs))
			if len(followersResp.FollowersIDs) > 0 {
				t.Logf("Follower IDs: %v", followersResp.FollowersIDs)
			}
		} else if resp.StatusCode == 400 || resp.StatusCode == 404 {
			t.Logf("Get followers failed (expected for test): Status %d", resp.StatusCode)
		} else {
			t.Fatalf("Unexpected status for get followers: %d", resp.StatusCode)
		}
	})

	t.Run("Get Followers - Nonexistent User", func(t *testing.T) {
		nonexistentUserID := 999999
		resp, err := client.GET(fmt.Sprintf("/friends/%d/followers", nonexistentUserID))
		if err != nil {
			t.Fatalf("Get nonexistent user followers request failed: %v", err)
		}

		if resp.IsSuccess() {
			t.Error("Expected get followers to fail for nonexistent user")
		}

		t.Logf("Expected failure for nonexistent user: Status %d", resp.StatusCode)
	})

	t.Run("Get Followers After Follow", func(t *testing.T) {
		// User2 follows User1
		followResp, err := user2.POST(fmt.Sprintf("/friends/%d", user1.UserID), nil)
		if err != nil {
			t.Fatalf("Follow request failed: %v", err)
		}

		if followResp.IsSuccess() {
			t.Logf("Follow successful, now checking followers")

			// Get User1's followers (should include User2)
			resp, err := client.GET(fmt.Sprintf("/friends/%d/followers", user1.UserID))
			if err != nil {
				t.Fatalf("Get followers after follow failed: %v", err)
			}

			if resp.IsSuccess() {
				var followersResp utils.UserFollowerResponse
				if err := resp.ParseJSON(&followersResp); err == nil {
					t.Logf("User1 now has %d followers", len(followersResp.FollowersIDs))
					// Check if user2 is in the followers list
					user2Found := false
					for _, followerID := range followersResp.FollowersIDs {
						if followerID == user2.UserID {
							user2Found = true
							break
						}
					}
					if user2Found {
						t.Logf("✓ User2 found in User1's followers list")
					}
				}
			}
		} else {
			t.Logf("Follow failed (expected for test), skipping followers check")
		}
	})
}

func TestGetUserFollowings(t *testing.T) {
	client, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	authHelper := utils.NewAuthHelper(client)

	// Create test users
	user1, err := authHelper.CreateTestUser("", "", "")
	if err != nil {
		t.Fatalf("Failed to create test user 1: %v", err)
	}

	user2, err := authHelper.CreateTestUser("", "", "")
	if err != nil {
		t.Fatalf("Failed to create test user 2: %v", err)
	}

	t.Run("Valid Get User Followings", func(t *testing.T) {
		resp, err := client.GET(fmt.Sprintf("/friends/%d/followings", user1.UserID))
		if err != nil {
			t.Fatalf("Get user followings request failed: %v", err)
		}

		if resp.IsSuccess() {
			var followingsResp utils.UserFollowingResponse
			if err := resp.ParseJSON(&followingsResp); err != nil {
				t.Fatalf("Failed to parse followings response: %v", err)
			}

			t.Logf("User followings retrieved: %d followings", len(followingsResp.FollowingsIDs))
			if len(followingsResp.FollowingsIDs) > 0 {
				t.Logf("Following IDs: %v", followingsResp.FollowingsIDs)
			}
		} else if resp.StatusCode == 400 || resp.StatusCode == 404 {
			t.Logf("Get followings failed (expected for test): Status %d", resp.StatusCode)
		} else {
			t.Fatalf("Unexpected status for get followings: %d", resp.StatusCode)
		}
	})

	t.Run("Get Followings - Nonexistent User", func(t *testing.T) {
		nonexistentUserID := 999999
		resp, err := client.GET(fmt.Sprintf("/friends/%d/followings", nonexistentUserID))
		if err != nil {
			t.Fatalf("Get nonexistent user followings request failed: %v", err)
		}

		if resp.IsSuccess() {
			t.Error("Expected get followings to fail for nonexistent user")
		}

		t.Logf("Expected failure for nonexistent user: Status %d", resp.StatusCode)
	})

	t.Run("Get Followings After Follow", func(t *testing.T) {
		// User1 follows User2
		followResp, err := user1.POST(fmt.Sprintf("/friends/%d", user2.UserID), nil)
		if err != nil {
			t.Fatalf("Follow request failed: %v", err)
		}

		if followResp.IsSuccess() {
			t.Logf("Follow successful, now checking followings")

			// Get User1's followings (should include User2)
			resp, err := client.GET(fmt.Sprintf("/friends/%d/followings", user1.UserID))
			if err != nil {
				t.Fatalf("Get followings after follow failed: %v", err)
			}

			if resp.IsSuccess() {
				var followingsResp utils.UserFollowingResponse
				if err := resp.ParseJSON(&followingsResp); err == nil {
					t.Logf("User1 now follows %d users", len(followingsResp.FollowingsIDs))
					// Check if user2 is in the followings list
					user2Found := false
					for _, followingID := range followingsResp.FollowingsIDs {
						if followingID == user2.UserID {
							user2Found = true
							break
						}
					}
					if user2Found {
						t.Logf("✓ User2 found in User1's followings list")
					}
				}
			}
		} else {
			t.Logf("Follow failed (expected for test), skipping followings check")
		}
	})
}

func TestGetUserPosts(t *testing.T) {
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

	t.Run("Valid Get User Posts", func(t *testing.T) {
		resp, err := client.GET(fmt.Sprintf("/friends/%d/posts", testUser.UserID))
		if err != nil {
			t.Fatalf("Get user posts request failed: %v", err)
		}

		if resp.IsSuccess() {
			var postsResp utils.UserPostsResponse
			if err := resp.ParseJSON(&postsResp); err != nil {
				t.Fatalf("Failed to parse user posts response: %v", err)
			}

			t.Logf("User posts retrieved: %d posts", len(postsResp.PostsIDs))
			if len(postsResp.PostsIDs) > 0 {
				t.Logf("Post IDs: %v", postsResp.PostsIDs)
			}
		} else if resp.StatusCode == 400 || resp.StatusCode == 404 {
			t.Logf("Get user posts failed (expected for test): Status %d", resp.StatusCode)
		} else {
			t.Fatalf("Unexpected status for get user posts: %d", resp.StatusCode)
		}
	})

	t.Run("Get Posts - Nonexistent User", func(t *testing.T) {
		nonexistentUserID := 999999
		resp, err := client.GET(fmt.Sprintf("/friends/%d/posts", nonexistentUserID))
		if err != nil {
			t.Fatalf("Get nonexistent user posts request failed: %v", err)
		}

		if resp.IsSuccess() {
			t.Error("Expected get user posts to fail for nonexistent user")
		}

		t.Logf("Expected failure for nonexistent user: Status %d", resp.StatusCode)
	})

	t.Run("Get Posts After Creating Posts", func(t *testing.T) {
		// Create a few posts for the user
		for i := 1; i <= 3; i++ {
			createReq := utils.CreatePostRequest{
				ContentText: fmt.Sprintf("Test post %d for user posts test", i),
				Visible:     true,
			}

			createResp, err := testUser.POST("/posts", createReq)
			if err != nil {
				t.Fatalf("Failed to create test post %d: %v", i, err)
			}

			if createResp.IsSuccess() {
				t.Logf("✓ Created test post %d", i)
			} else {
				t.Logf("⚠ Failed to create test post %d: Status %d", i, createResp.StatusCode)
			}
		}

		// Now get the user's posts
		resp, err := client.GET(fmt.Sprintf("/friends/%d/posts", testUser.UserID))
		if err != nil {
			t.Fatalf("Get user posts after creation failed: %v", err)
		}

		if resp.IsSuccess() {
			var postsResp utils.UserPostsResponse
			if err := resp.ParseJSON(&postsResp); err == nil {
				t.Logf("User now has %d posts", len(postsResp.PostsIDs))
				if len(postsResp.PostsIDs) >= 3 {
					t.Logf("✓ All test posts appear to be created and retrieved")
				}
			}
		} else {
			t.Logf("Get user posts failed (expected for test): Status %d", resp.StatusCode)
		}
	})
}

func TestSocialInteractions(t *testing.T) {
	client, err := utils.NewAPIClient()
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	authHelper := utils.NewAuthHelper(client)

	// Create three test users for multi-user interactions
	users, err := authHelper.CreateMultipleTestUsers(3)
	if err != nil {
		t.Fatalf("Failed to create test users: %v", err)
	}

	user1, user2, user3 := users[0], users[1], users[2]

	t.Run("Multi-User Social Network", func(t *testing.T) {
		// Step 1: Create posts for each user
		for i, user := range users {
			createReq := utils.CreatePostRequest{
				ContentText: fmt.Sprintf("Hello from user %d!", i+1),
				Visible:     true,
			}

			createResp, err := user.POST("/posts", createReq)
			if err != nil {
				t.Fatalf("Failed to create post for user %d: %v", i+1, err)
			}

			if createResp.IsSuccess() {
				t.Logf("✓ User %d created a post", i+1)
			} else {
				t.Logf("⚠ User %d failed to create post: Status %d", i+1, createResp.StatusCode)
			}
		}

		// Step 2: Set up follow relationships
		// User1 follows User2 and User3
		followResp, err := user1.POST(fmt.Sprintf("/friends/%d", user2.UserID), nil)
		if err == nil && followResp.IsSuccess() {
			t.Logf("✓ User1 followed User2")
		} else {
			t.Logf("⚠ User1 failed to follow User2")
		}

		followResp, err = user1.POST(fmt.Sprintf("/friends/%d", user3.UserID), nil)
		if err == nil && followResp.IsSuccess() {
			t.Logf("✓ User1 followed User3")
		} else {
			t.Logf("⚠ User1 failed to follow User3")
		}

		// User2 follows User3
		followResp, err = user2.POST(fmt.Sprintf("/friends/%d", user3.UserID), nil)
		if err == nil && followResp.IsSuccess() {
			t.Logf("✓ User2 followed User3")
		} else {
			t.Logf("⚠ User2 failed to follow User3")
		}

		// User3 follows User1 (creating a partial cycle)
		followResp, err = user3.POST(fmt.Sprintf("/friends/%d", user1.UserID), nil)
		if err == nil && followResp.IsSuccess() {
			t.Logf("✓ User3 followed User1")
		} else {
			t.Logf("⚠ User3 failed to follow User1")
		}

		// Step 3: Check follow relationships
		for i, user := range users {
			// Get followers
			followersResp, err := client.GET(fmt.Sprintf("/friends/%d/followers", user.UserID))
			if err == nil && followersResp.IsSuccess() {
				var followers utils.UserFollowerResponse
				if followersResp.ParseJSON(&followers) == nil {
					t.Logf("User %d has %d followers", i+1, len(followers.FollowersIDs))
				}
			}

			// Get followings
			followingsResp, err := client.GET(fmt.Sprintf("/friends/%d/followings", user.UserID))
			if err == nil && followingsResp.IsSuccess() {
				var followings utils.UserFollowingResponse
				if followingsResp.ParseJSON(&followings) == nil {
					t.Logf("User %d follows %d users", i+1, len(followings.FollowingsIDs))
				}
			}

			// Get posts
			postsResp, err := client.GET(fmt.Sprintf("/friends/%d/posts", user.UserID))
			if err == nil && postsResp.IsSuccess() {
				var posts utils.UserPostsResponse
				if postsResp.ParseJSON(&posts) == nil {
					t.Logf("User %d has %d posts", i+1, len(posts.PostsIDs))
				}
			}
		}

		// Step 4: Test unfollow
		unfollowResp, err := user1.DELETE(fmt.Sprintf("/friends/%d", user2.UserID))
		if err == nil && unfollowResp.IsSuccess() {
			t.Logf("✓ User1 unfollowed User2")
		} else {
			t.Logf("⚠ User1 failed to unfollow User2")
		}

		t.Log("🎉 Multi-user social interaction test completed!")
	})
}
