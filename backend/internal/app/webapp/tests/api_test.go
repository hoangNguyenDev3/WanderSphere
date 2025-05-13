// Package tests provides API tests for the webapp
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
)

// Setup test router with API endpoints
func setupTestRouter(t *testing.T) (*gin.Engine, *MockWebappService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(MockWebappService)

	// Define routes and handlers similar to the actual API
	// This would normally be in your server setup code
	api := router.Group("/api/v1")
	{
		// User routes
		api.POST("/users/signup", func(c *gin.Context) {
			var req types.CreateUserRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid request",
					Status:  "error",
				})
				return
			}

			resp := mockService.CreateUser(req)
			c.JSON(http.StatusOK, resp)
		})

		api.POST("/users/login", func(c *gin.Context) {
			var req types.LoginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, types.ErrorResponse{
					Message: "Invalid request",
					Error:   err.Error(),
					Code:    http.StatusBadRequest,
				})
				return
			}

			resp, err := mockService.Login(req)
			if err != nil {
				c.JSON(http.StatusBadRequest, types.ErrorResponse{
					Message: "Login failed",
					Error:   err.Error(),
					Code:    http.StatusBadRequest,
				})
				return
			}
			c.JSON(http.StatusOK, resp)
		})

		api.GET("/users/:user_id", func(c *gin.Context) {
			userIDStr := c.Param("user_id")
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid user ID",
					Status:  "error",
				})
				return
			}

			resp := mockService.GetUserDetails(userID)
			c.JSON(http.StatusOK, resp)
		})

		api.POST("/users/edit", func(c *gin.Context) {
			var req types.EditUserRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid request",
					Status:  "error",
				})
				return
			}

			resp := mockService.EditUser(req)
			c.JSON(http.StatusOK, resp)
		})

		// Post routes
		api.POST("/posts", func(c *gin.Context) {
			var req types.CreatePostRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid request",
					Status:  "error",
				})
				return
			}

			resp := mockService.CreatePost(req)
			c.JSON(http.StatusOK, resp)
		})

		api.GET("/posts/:post_id", func(c *gin.Context) {
			postIDStr := c.Param("post_id")
			postID, err := strconv.ParseInt(postIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid post ID",
					Status:  "error",
				})
				return
			}

			resp := mockService.GetPostDetails(postID)
			c.JSON(http.StatusOK, resp)
		})

		api.PUT("/posts/:post_id", func(c *gin.Context) {
			postIDStr := c.Param("post_id")
			postID, err := strconv.ParseInt(postIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid post ID",
					Status:  "error",
				})
				return
			}

			var req types.EditPostRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid request",
					Status:  "error",
				})
				return
			}

			resp := mockService.EditPost(postID, req)
			c.JSON(http.StatusOK, resp)
		})

		api.POST("/posts/:post_id", func(c *gin.Context) {
			postIDStr := c.Param("post_id")
			postID, err := strconv.ParseInt(postIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid post ID",
					Status:  "error",
				})
				return
			}

			var req types.CreatePostCommentRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid request",
					Status:  "error",
				})
				return
			}

			resp := mockService.CommentOnPost(postID, req)
			c.JSON(http.StatusOK, resp)
		})

		api.DELETE("/posts/:post_id", func(c *gin.Context) {
			postIDStr := c.Param("post_id")
			postID, err := strconv.ParseInt(postIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid post ID",
					Status:  "error",
				})
				return
			}

			resp := mockService.DeletePost(postID)
			c.JSON(http.StatusOK, resp)
		})

		api.POST("/posts/:post_id/likes", func(c *gin.Context) {
			postIDStr := c.Param("post_id")
			postID, err := strconv.ParseInt(postIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid post ID",
					Status:  "error",
				})
				return
			}

			resp := mockService.LikePost(postID)
			c.JSON(http.StatusOK, resp)
		})

		api.GET("/posts/url", func(c *gin.Context) {
			var req types.GetS3PresignedUrlRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid request",
					Status:  "error",
				})
				return
			}

			resp := mockService.GetS3PresignedUrl(req)
			c.JSON(http.StatusOK, resp)
		})

		// Friend routes
		api.POST("/friends/:user_id", func(c *gin.Context) {
			userIDStr := c.Param("user_id")
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid user ID",
					Status:  "error",
				})
				return
			}

			resp := mockService.FollowUser(userID)
			c.JSON(http.StatusOK, resp)
		})

		api.DELETE("/friends/:user_id", func(c *gin.Context) {
			userIDStr := c.Param("user_id")
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid user ID",
					Status:  "error",
				})
				return
			}

			resp := mockService.UnfollowUser(userID)
			c.JSON(http.StatusOK, resp)
		})

		api.GET("/friends/:user_id/followers", func(c *gin.Context) {
			userIDStr := c.Param("user_id")
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid user ID",
					Status:  "error",
				})
				return
			}

			resp := mockService.GetUserFollowers(userID)
			c.JSON(http.StatusOK, resp)
		})

		api.GET("/friends/:user_id/followings", func(c *gin.Context) {
			userIDStr := c.Param("user_id")
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid user ID",
					Status:  "error",
				})
				return
			}

			resp := mockService.GetUserFollowings(userID)
			c.JSON(http.StatusOK, resp)
		})

		api.GET("/friends/:user_id/posts", func(c *gin.Context) {
			userIDStr := c.Param("user_id")
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, types.MessageResponse{
					Message: "Invalid user ID",
					Status:  "error",
				})
				return
			}

			resp := mockService.GetUserPosts(userID)
			c.JSON(http.StatusOK, resp)
		})

		// Newsfeed routes
		api.GET("/newsfeed", func(c *gin.Context) {
			resp := mockService.GetNewsfeed()
			c.JSON(http.StatusOK, resp)
		})
	}

	return router, mockService
}

// User API Tests

func TestUserSignup(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	expectedResponse := types.MessageResponse{
		Message: "User created successfully",
		Status:  "success",
	}

	mockService.On("CreateUser", mock.AnythingOfType("types.CreateUserRequest")).Return(expectedResponse)

	// Create test request
	request := types.CreateUserRequest{
		UserName:    "testuser",
		Email:       "test@example.com",
		Password:    "password123",
		DateOfBirth: "1990-01-01",
		FirstName:   "Test",
		LastName:    "User",
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/signup", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Message, response.Message)
	assert.Equal(t, expectedResponse.Status, response.Status)

	mockService.AssertExpectations(t)
}

func TestUserLogin(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	expectedUser := types.UserDetailInfo{
		UserID:      int64(1),
		UserName:    "testuser",
		Email:       "test@example.com",
		FirstName:   "Test",
		LastName:    "User",
		DateOfBirth: "1990-01-01",
	}

	expectedResponse := types.LoginResponse{
		Message: "Login successful",
		User:    expectedUser,
	}

	mockService.On("Login", mock.AnythingOfType("types.LoginRequest")).Return(expectedResponse, nil)

	// Create test request
	request := types.LoginRequest{
		UserName: "testuser",
		Password: "password123",
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Message, response.Message)
	assert.Equal(t, expectedResponse.User.UserID, response.User.UserID)
	assert.Equal(t, expectedResponse.User.UserName, response.User.UserName)

	mockService.AssertExpectations(t)
}

func TestGetUserDetails(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	userID := int64(1)
	expectedResponse := types.UserDetailInfoResponse{
		UserID:      userID,
		UserName:    "testuser",
		Email:       "test@example.com",
		FirstName:   "Test",
		LastName:    "User",
		DateOfBirth: "1990-01-01",
	}

	mockService.On("GetUserDetails", userID).Return(expectedResponse)

	// Create test request
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%d", userID), nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.UserDetailInfoResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.UserID, response.UserID)
	assert.Equal(t, expectedResponse.UserName, response.UserName)

	mockService.AssertExpectations(t)
}

func TestEditUserProfile(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	expectedResponse := types.MessageResponse{
		Message: "User updated successfully",
		Status:  "success",
	}

	mockService.On("EditUser", mock.AnythingOfType("types.EditUserRequest")).Return(expectedResponse)

	// Create test request
	request := types.EditUserRequest{
		FirstName:      "Updated",
		LastName:       "Name",
		DateOfBirth:    "1990-01-01",
		ProfilePicture: "http://example.com/profile.jpg",
		CoverPicture:   "http://example.com/cover.jpg",
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/edit", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Message, response.Message)
	assert.Equal(t, expectedResponse.Status, response.Status)

	mockService.AssertExpectations(t)
}

// Post API Tests

func TestCreatePost(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	expectedResponse := types.MessageResponse{
		Message: "Post created successfully",
		Status:  "success",
	}

	mockService.On("CreatePost", mock.AnythingOfType("types.CreatePostRequest")).Return(expectedResponse)

	// Create test request
	visible := true
	request := types.CreatePostRequest{
		ContentText:      "This is a test post",
		ContentImagePath: []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
		Visible:          &visible,
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Message, response.Message)
	assert.Equal(t, expectedResponse.Status, response.Status)

	mockService.AssertExpectations(t)
}

func TestGetPostDetails(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	postID := int64(1)
	userID := int64(1)

	expectedResponse := types.PostDetailInfoResponse{
		PostID:           postID,
		UserID:           userID,
		ContentText:      "This is a test post",
		ContentImagePath: []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
		CreatedAt:        time.Now().Format(time.RFC3339),
		UsersLiked:       []int64{2, 3},
		Comments:         []types.CommentResponse{},
	}

	mockService.On("GetPostDetails", postID).Return(expectedResponse)

	// Create test request
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/posts/%d", postID), nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.PostDetailInfoResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.PostID, response.PostID)
	assert.Equal(t, expectedResponse.UserID, response.UserID)
	assert.Equal(t, expectedResponse.ContentText, response.ContentText)

	mockService.AssertExpectations(t)
}

func TestEditPost(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	postID := int64(1)
	expectedResponse := types.MessageResponse{
		Message: "Post updated successfully",
		Status:  "success",
	}

	mockService.On("EditPost", postID, mock.AnythingOfType("types.EditPostRequest")).Return(expectedResponse)

	// Create test request
	contentText := "Updated post content"
	contentImagePath := []string{"http://example.com/updated_image.jpg"}
	visible := true

	request := types.EditPostRequest{
		ContentText:      &contentText,
		ContentImagePath: &contentImagePath,
		Visible:          &visible,
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/posts/%d", postID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Message, response.Message)
	assert.Equal(t, expectedResponse.Status, response.Status)

	mockService.AssertExpectations(t)
}

func TestDeletePost(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	postID := int64(1)
	expectedResponse := types.MessageResponse{
		Message: "Post deleted successfully",
		Status:  "success",
	}

	mockService.On("DeletePost", postID).Return(expectedResponse)

	// Create test request
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/posts/%d", postID), nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Message, response.Message)
	assert.Equal(t, expectedResponse.Status, response.Status)

	mockService.AssertExpectations(t)
}

func TestCommentOnPost(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	postID := int64(1)

	expectedResponse := types.PostDetailInfoResponse{
		PostID:           postID,
		UserID:           int64(1),
		ContentText:      "This is a test post",
		ContentImagePath: []string{"http://example.com/image1.jpg"},
		CreatedAt:        time.Now().Format(time.RFC3339),
		Comments: []types.CommentResponse{
			{
				CommentId:   int64(1),
				PostId:      postID,
				UserId:      int64(2),
				ContentText: "This is a comment",
			},
		},
	}

	mockService.On("CommentOnPost", postID, mock.AnythingOfType("types.CreatePostCommentRequest")).Return(expectedResponse)

	// Create test request
	request := types.CreatePostCommentRequest{
		ContentText: "This is a comment",
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/posts/%d", postID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.PostDetailInfoResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.PostID, response.PostID)
	assert.Greater(t, len(response.Comments), 0)

	mockService.AssertExpectations(t)
}

func TestLikePost(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	postID := int64(1)
	expectedResponse := types.MessageResponse{
		Message: "Post liked successfully",
		Status:  "success",
	}

	mockService.On("LikePost", postID).Return(expectedResponse)

	// Create test request
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/posts/%d/likes", postID), nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Message, response.Message)
	assert.Equal(t, expectedResponse.Status, response.Status)

	mockService.AssertExpectations(t)
}

func TestGetS3PresignedUrl(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	expectedResponse := types.GetS3PresignedUrlResponse{
		URL:            "https://s3-bucket.amazonaws.com/file.jpg?signature=xyz",
		ExpirationTime: time.Now().Add(time.Hour).Format(time.RFC3339),
	}

	mockService.On("GetS3PresignedUrl", mock.AnythingOfType("types.GetS3PresignedUrlRequest")).Return(expectedResponse)

	// Create test request
	request := types.GetS3PresignedUrlRequest{
		FileName: "file.jpg",
		FileType: "image/jpeg",
	}

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/posts/url", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.GetS3PresignedUrlResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.URL, response.URL)

	mockService.AssertExpectations(t)
}

// Friend API Tests

func TestFollowUser(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	userID := int64(2)
	expectedResponse := types.MessageResponse{
		Message: "User followed successfully",
		Status:  "success",
	}

	mockService.On("FollowUser", userID).Return(expectedResponse)

	// Create test request
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/friends/%d", userID), nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Message, response.Message)
	assert.Equal(t, expectedResponse.Status, response.Status)

	mockService.AssertExpectations(t)
}

func TestUnfollowUser(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	userID := int64(2)
	expectedResponse := types.MessageResponse{
		Message: "User unfollowed successfully",
		Status:  "success",
	}

	mockService.On("UnfollowUser", userID).Return(expectedResponse)

	// Create test request
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/friends/%d", userID), nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Message, response.Message)
	assert.Equal(t, expectedResponse.Status, response.Status)

	mockService.AssertExpectations(t)
}

func TestGetUserFollowers(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	userID := int64(1)
	expectedResponse := types.UserFollowerResponse{
		FollowersIds: []int64{2, 3, 4},
	}

	mockService.On("GetUserFollowers", userID).Return(expectedResponse)

	// Create test request
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/friends/%d/followers", userID), nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.UserFollowerResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.FollowersIds, response.FollowersIds)

	mockService.AssertExpectations(t)
}

func TestGetUserFollowings(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	userID := int64(1)
	expectedResponse := types.UserFollowingResponse{
		FollowingsIds: []int64{5, 6, 7},
	}

	mockService.On("GetUserFollowings", userID).Return(expectedResponse)

	// Create test request
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/friends/%d/followings", userID), nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.UserFollowingResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.FollowingsIds, response.FollowingsIds)

	mockService.AssertExpectations(t)
}

func TestGetUserPosts(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	userID := int64(1)
	expectedResponse := types.UserPostsResponse{
		PostsIds: []int64{1, 2, 3},
	}

	mockService.On("GetUserPosts", userID).Return(expectedResponse)

	// Create test request
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/friends/%d/posts", userID), nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.UserPostsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.PostsIds, response.PostsIds)

	mockService.AssertExpectations(t)
}

// Newsfeed API Tests

func TestGetNewsfeed(t *testing.T) {
	router, mockService := setupTestRouter(t)

	// Setup mock expectations
	expectedResponse := types.NewsfeedResponse{
		PostsIds: []int64{1, 2, 3, 4, 5},
	}

	mockService.On("GetNewsfeed").Return(expectedResponse)

	// Create test request
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/newsfeed", nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)

	var response types.NewsfeedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.PostsIds, response.PostsIds)

	mockService.AssertExpectations(t)
}

// Date parsing test
func TestParsableDate(t *testing.T) {
	dateStr := time.Now().Format(time.DateOnly)
	_, err := time.Parse(time.DateOnly, dateStr)
	assert.NoError(t, err, "Date format should be parsable")
}
