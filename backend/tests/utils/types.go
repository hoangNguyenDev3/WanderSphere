package utils

import "time"

// Request types from swagger.json

type CreateUserRequest struct {
	UserName    string `json:"user_name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	DateOfBirth string `json:"date_of_birth" binding:"required"`
}

type LoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type EditUserRequest struct {
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	DateOfBirth    string `json:"date_of_birth,omitempty"`
	Password       string `json:"password,omitempty"`
	ProfilePicture string `json:"profile_picture,omitempty"`
	CoverPicture   string `json:"cover_picture,omitempty"`
}

type CreatePostRequest struct {
	ContentText      string   `json:"content_text" binding:"required"`
	ContentImagePath []string `json:"content_image_path,omitempty"`
	Visible          bool     `json:"visible,omitempty"`
}

type EditPostRequest struct {
	ContentText      string   `json:"content_text,omitempty"`
	ContentImagePath []string `json:"content_image_path,omitempty"`
	Visible          bool     `json:"visible,omitempty"`
}

type CreatePostCommentRequest struct {
	ContentText string `json:"content_text" binding:"required"`
}

type GetS3PresignedUrlRequest struct {
	FileName string `json:"file_name" binding:"required"`
	FileType string `json:"file_type" binding:"required"`
}

// Response types from swagger.json

type MessageResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type UserDetailInfo struct {
	UserID         int    `json:"user_id"`
	UserName       string `json:"user_name"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DateOfBirth    string `json:"date_of_birth"`
	ProfilePicture string `json:"profile_picture"`
	CoverPicture   string `json:"cover_picture"`
}

type UserDetailInfoResponse UserDetailInfo

type LoginResponse struct {
	Message string         `json:"message"`
	User    UserDetailInfo `json:"user"`
}

type CommentResponse struct {
	CommentID   int    `json:"comment_id"`
	PostID      int    `json:"post_id"`
	UserID      int    `json:"user_id"`
	ContentText string `json:"content_text"`
}

type PostDetailInfoResponse struct {
	PostID           int               `json:"post_id"`
	UserID           int               `json:"user_id"`
	ContentText      string            `json:"content_text"`
	ContentImagePath []string          `json:"content_image_path"`
	CreatedAt        string            `json:"created_at"`
	UsersLiked       []int             `json:"users_liked"`
	Comments         []CommentResponse `json:"comments"`
}

type UserFollowerResponse struct {
	FollowersIDs []int `json:"followers_ids"`
}

type UserFollowingResponse struct {
	FollowingsIDs []int `json:"followings_ids"`
}

type UserPostsResponse struct {
	PostsIds []int64 `json:"posts_ids"`
}

type NewsfeedResponse struct {
	PostsIds []int64 `json:"posts_ids"`
}

type GetS3PresignedUrlResponse struct {
	URL            string `json:"url"`
	ExpirationTime string `json:"expiration_time"`
}

// Test data structures

type TestUser struct {
	ID              int               `json:"id,omitempty"`
	CreateRequest   CreateUserRequest `json:"create_request"`
	LoginRequest    LoginRequest      `json:"login_request"`
	ExpectedProfile UserDetailInfo    `json:"expected_profile"`
	SessionID       string            `json:"session_id,omitempty"`
	CreatedAt       *time.Time        `json:"created_at,omitempty"`
}

type TestPost struct {
	ID            int               `json:"id,omitempty"`
	CreateRequest CreatePostRequest `json:"create_request"`
	AuthorID      int               `json:"author_id"`
	CreatedAt     *time.Time        `json:"created_at,omitempty"`
}

type TestComment struct {
	ID            int                      `json:"id,omitempty"`
	CreateRequest CreatePostCommentRequest `json:"create_request"`
	PostID        int                      `json:"post_id"`
	AuthorID      int                      `json:"author_id"`
	CreatedAt     *time.Time               `json:"created_at,omitempty"`
}

// Test scenario data

type TestScenario struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Users       []TestUser    `json:"users"`
	Posts       []TestPost    `json:"posts"`
	Comments    []TestComment `json:"comments"`
	Follows     []struct {
		FollowerID int `json:"follower_id"`
		FolloweeID int `json:"followee_id"`
	} `json:"follows"`
}
