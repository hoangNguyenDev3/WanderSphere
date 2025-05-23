package docs

// This file contains type definitions for Swagger documentation only
// These types mirror the actual types defined in the application

// LoginRequest represents user login credentials
type LoginRequest struct {
	UserName string `json:"user_name" example:"johndoe"`
	Password string `json:"password" example:"p@ssw0rd"`
}

// CreateUserRequest represents user registration data
type CreateUserRequest struct {
	UserName    string `json:"user_name" example:"johndoe"`
	Password    string `json:"password" example:"p@ssw0rd"`
	FirstName   string `json:"first_name" example:"John"`
	LastName    string `json:"last_name" example:"Doe"`
	DateOfBirth string `json:"date_of_birth" example:"1990-01-01"`
	Email       string `json:"email" example:"john.doe@example.com"`
}

// EditUserRequest represents user profile update data
type EditUserRequest struct {
	Password       string `json:"password" example:"newP@ssw0rd"`
	FirstName      string `json:"first_name" example:"Johnny"`
	LastName       string `json:"last_name" example:"Doe"`
	DateOfBirth    string `json:"date_of_birth" example:"1990-01-01"`
	ProfilePicture string `json:"profile_picture" example:"https://example.com/profile.jpg"`
	CoverPicture   string `json:"cover_picture" example:"https://example.com/cover.jpg"`
}

// CreatePostRequest represents a post creation request
type CreatePostRequest struct {
	ContentText      string   `json:"content_text" example:"Hello world!"`
	ContentImagePath []string `json:"content_image_path" example:"[\"https://example.com/image.jpg\"]"`
	Visible          bool     `json:"visible" example:"true"`
}

// EditPostRequest represents a post update request
type EditPostRequest struct {
	ContentText string `json:"content_text" example:"Updated post content"`
	Visible     bool   `json:"visible" example:"true"`
}

// CreatePostCommentRequest represents a comment creation request
type CreatePostCommentRequest struct {
	ContentText string `json:"content_text" example:"Great post!"`
}

// GetS3PresignedUrlRequest represents a request to get a presigned S3 URL
type GetS3PresignedUrlRequest struct {
	FileName string `json:"file_name" example:"image.jpg"`
	FileType string `json:"file_type" example:"image/jpeg"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
	Status  string `json:"status,omitempty" example:"success"`
}

// LoginResponse represents a login response with user information
type LoginResponse struct {
	Message string         `json:"message" example:"Login successful"`
	User    UserDetailInfo `json:"user"`
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string `json:"error" example:"validation_error"`
	Message string `json:"message,omitempty" example:"Invalid input data"`
	Code    int    `json:"code,omitempty" example:"400"`
}

// PostDetailInfoResponse represents detailed information about a post
type PostDetailInfoResponse struct {
	PostID           int64             `json:"post_id" example:"123"`
	UserID           int64             `json:"user_id" example:"456"`
	ContentText      string            `json:"content_text" example:"This is a post"`
	ContentImagePath []string          `json:"content_image_path" example:"[\"https://example.com/image.jpg\"]"`
	CreatedAt        string            `json:"created_at" example:"2023-01-01T12:00:00Z"`
	Comments         []CommentResponse `json:"comments"`
	UsersLiked       []int64           `json:"users_liked" example:"[789,101]"`
}

// CommentResponse represents a comment on a post
type CommentResponse struct {
	CommentId   int64  `json:"comment_id" example:"123"`
	UserId      int64  `json:"user_id" example:"456"`
	PostId      int64  `json:"post_id" example:"789"`
	ContentText string `json:"content_text" example:"Great post!"`
}

// UserDetailInfo represents detailed user information
type UserDetailInfo struct {
	UserID         int64  `json:"user_id" example:"123"`
	UserName       string `json:"user_name" example:"johndoe"`
	FirstName      string `json:"first_name" example:"John"`
	LastName       string `json:"last_name" example:"Doe"`
	DateOfBirth    string `json:"date_of_birth" example:"1990-01-01"`
	Email          string `json:"email" example:"john.doe@example.com"`
	ProfilePicture string `json:"profile_picture,omitempty" example:"https://example.com/profile.jpg"`
	CoverPicture   string `json:"cover_picture,omitempty" example:"https://example.com/cover.jpg"`
}

// UserDetailInfoResponse represents a response with user details
type UserDetailInfoResponse struct {
	UserID         int64  `json:"user_id" example:"123"`
	UserName       string `json:"user_name" example:"johndoe"`
	FirstName      string `json:"first_name" example:"John"`
	LastName       string `json:"last_name" example:"Doe"`
	DateOfBirth    string `json:"date_of_birth" example:"1990-01-01"`
	Email          string `json:"email" example:"john.doe@example.com"`
	ProfilePicture string `json:"profile_picture,omitempty" example:"https://example.com/profile.jpg"`
	CoverPicture   string `json:"cover_picture,omitempty" example:"https://example.com/cover.jpg"`
}

// GetS3PresignedUrlResponse represents a response with a presigned S3 URL
type GetS3PresignedUrlResponse struct {
	URL            string `json:"url" example:"https://s3.example.com/bucket/path/image.jpg?signature=abc"`
	ExpirationTime string `json:"expiration_time" example:"2023-01-01T12:15:00Z"`
}

// UserFollowerResponse represents a response with a user's followers
type UserFollowerResponse struct {
	FollowersIds []int64 `json:"followers_ids" example:"[123,456]"`
}

// UserFollowingResponse represents a response with a user's followings
type UserFollowingResponse struct {
	FollowingsIds []int64 `json:"followings_ids" example:"[789,101]"`
}

// UserPostsResponse represents a response with a user's posts
type UserPostsResponse struct {
	PostsIds []int64 `json:"posts_ids" example:"[123,456]"`
}

// NewsfeedResponse represents a response with a user's newsfeed
type NewsfeedResponse struct {
	PostsIds []int64 `json:"posts_ids" example:"[123,456]"`
}
