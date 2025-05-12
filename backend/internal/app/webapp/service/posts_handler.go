package service

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/hoangNguyenDev3/WanderSphere/backend/docs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"

	pb_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
)

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post with text and optional images
// @Tags posts
// @Accept json
// @Produce json
// @Param request body types.CreatePostRequest true "Post creation parameters"
// @Success 200 {object} types.MessageResponse "Post created successfully"
// @Failure 400 {object} types.MessageResponse "Validation error"
// @Failure 401 {object} types.MessageResponse "Unauthorized"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /posts [post]
// @Security ApiKeyAuth
func (svc *WebService) CreatePost(ctx *gin.Context) {
	// Check authorization
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Validate request
	var jsonRequest types.CreatePostRequest
	err = ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	// Set visibility default if not provided
	visible := true
	if jsonRequest.Visible != nil {
		visible = *jsonRequest.Visible
	}

	// Call CreatePost service
	resp, err := svc.AuthenticateAndPostClient.CreatePost(ctx, &pb_aap.CreatePostRequest{
		UserId:           int64(userId),
		ContentText:      jsonRequest.ContentText,
		ContentImagePath: jsonRequest.ContentImagePath,
		Visible:          visible,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.CreatePostResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.CreatePostResponse_OK {
		ctx.JSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// GetPostDetail godoc
// @Summary Get post details
// @Description Get detailed information about a post
// @Tags posts
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Success 200 {object} types.PostDetailInfoResponse "Post details"
// @Failure 400 {object} types.MessageResponse "Invalid post ID or post not found"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /posts/{post_id} [get]
func (svc *WebService) GetPostDetail(ctx *gin.Context) {
	// Check URL params
	postId, err := strconv.Atoi(ctx.Param("post_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	// Call grpc service
	resp, err := svc.AuthenticateAndPostClient.GetPostDetailInfo(ctx, &pb_aap.GetPostDetailInfoRequest{
		PostId: int64(postId),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.GetPostDetailInfoResponse_POST_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	} else if resp.GetStatus() == pb_aap.GetPostDetailInfoResponse_OK {
		// Convert comments
		comments := make([]types.CommentResponse, 0)
		for _, comment := range resp.GetPost().GetComments() {
			comments = append(comments, types.CommentResponse{
				CommentId:   comment.GetCommentId(),
				UserId:      comment.GetUserId(),
				PostId:      comment.GetPostId(),
				ContentText: comment.GetContentText(),
			})
		}

		// Convert likes
		usersLiked := make([]int64, 0)
		for _, like := range resp.GetPost().GetLikedUsers() {
			usersLiked = append(usersLiked, like.GetUserId())
		}

		ctx.JSON(http.StatusOK, types.PostDetailInfoResponse{
			PostID:           resp.GetPost().GetPostId(),
			UserID:           resp.GetPost().GetUserId(),
			ContentText:      resp.GetPost().GetContentText(),
			ContentImagePath: resp.GetPost().GetContentImagePath(),
			CreatedAt:        resp.GetPost().GetCreatedAt().AsTime().Format(time.RFC3339),
			Comments:         comments,
			UsersLiked:       usersLiked,
		})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// EditPost godoc
// @Summary Edit post
// @Description Edit an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Param request body types.EditPostRequest true "Post edit parameters"
// @Success 200 {object} types.MessageResponse "Post updated successfully"
// @Failure 400 {object} types.MessageResponse "Validation error or post not found"
// @Failure 401 {object} types.MessageResponse "Unauthorized"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /posts/{post_id} [put]
// @Security ApiKeyAuth
func (svc *WebService) EditPost(ctx *gin.Context) {
	// Check authorization
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Check URL params
	postId, err := strconv.Atoi(ctx.Param("post_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	// Validate request
	var jsonRequest types.EditPostRequest
	err = ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	// Call EditPost service
	resp, err := svc.AuthenticateAndPostClient.EditPost(ctx, &pb_aap.EditPostRequest{
		UserId:      int64(userId),
		PostId:      int64(postId),
		ContentText: jsonRequest.ContentText,
		Visible:     jsonRequest.Visible,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.EditPostResponse_POST_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	} else if resp.GetStatus() == pb_aap.EditPostResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.EditPostResponse_NOT_ALLOWED {
		ctx.JSON(http.StatusUnauthorized, types.MessageResponse{Message: "not allowed to edit this post"})
		return
	} else if resp.GetStatus() == pb_aap.EditPostResponse_OK {
		ctx.JSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// DeletePost godoc
// @Summary Delete post
// @Description Delete an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Success 200 {object} types.MessageResponse "Post deleted successfully"
// @Failure 400 {object} types.MessageResponse "Validation error or post not found"
// @Failure 401 {object} types.MessageResponse "Unauthorized"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /posts/{post_id} [delete]
// @Security ApiKeyAuth
func (svc *WebService) DeletePost(ctx *gin.Context) {
	// Check authorization
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Check URL params
	postId, err := strconv.Atoi(ctx.Param("post_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	// Call DeletePost service
	resp, err := svc.AuthenticateAndPostClient.DeletePost(ctx, &pb_aap.DeletePostRequest{
		UserId: int64(userId),
		PostId: int64(postId),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.DeletePostResponse_POST_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	} else if resp.GetStatus() == pb_aap.DeletePostResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.DeletePostResponse_NOT_ALLOWED {
		ctx.JSON(http.StatusUnauthorized, types.MessageResponse{Message: "not allowed to delete this post"})
		return
	} else if resp.GetStatus() == pb_aap.DeletePostResponse_OK {
		ctx.JSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// CommentPost godoc
// @Summary Comment on a post
// @Description Add a comment to an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Param request body types.CreatePostCommentRequest true "Comment content"
// @Success 200 {object} types.PostDetailInfoResponse "Updated post with new comment"
// @Failure 400 {object} types.MessageResponse "Validation error or post not found"
// @Failure 401 {object} types.MessageResponse "Unauthorized"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /posts/{post_id} [post]
// @Security ApiKeyAuth
func (svc *WebService) CommentPost(ctx *gin.Context) {
	// Check authorization
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Check URL params
	postId, err := strconv.Atoi(ctx.Param("post_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	// Validate request
	var jsonRequest types.CreatePostCommentRequest
	err = ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	// Call CommentPost service
	resp, err := svc.AuthenticateAndPostClient.CommentPost(ctx, &pb_aap.CommentPostRequest{
		UserId:      int64(userId),
		PostId:      int64(postId),
		ContentText: jsonRequest.ContentText,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.CommentPostResponse_POST_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	} else if resp.GetStatus() == pb_aap.CommentPostResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.CommentPostResponse_OK {
		// Get updated post details
		postResp, err := svc.AuthenticateAndPostClient.GetPostDetailInfo(ctx, &pb_aap.GetPostDetailInfoRequest{
			PostId: int64(postId),
		})
		if err != nil {
			ctx.JSON(http.StatusOK, types.MessageResponse{Message: "Comment added successfully"})
			return
		}

		// Convert comments
		comments := make([]types.CommentResponse, 0)
		for _, comment := range postResp.GetPost().GetComments() {
			comments = append(comments, types.CommentResponse{
				CommentId:   comment.GetCommentId(),
				UserId:      comment.GetUserId(),
				PostId:      comment.GetPostId(),
				ContentText: comment.GetContentText(),
			})
		}

		// Convert likes
		usersLiked := make([]int64, 0)
		for _, like := range postResp.GetPost().GetLikedUsers() {
			usersLiked = append(usersLiked, like.GetUserId())
		}

		ctx.JSON(http.StatusOK, types.PostDetailInfoResponse{
			PostID:           postResp.GetPost().GetPostId(),
			UserID:           postResp.GetPost().GetUserId(),
			ContentText:      postResp.GetPost().GetContentText(),
			ContentImagePath: postResp.GetPost().GetContentImagePath(),
			CreatedAt:        postResp.GetPost().GetCreatedAt().AsTime().Format(time.RFC3339),
			Comments:         comments,
			UsersLiked:       usersLiked,
		})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// LikePost godoc
// @Summary Like a post
// @Description Like an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Success 200 {object} types.MessageResponse "Post liked successfully"
// @Failure 400 {object} types.MessageResponse "Validation error or post not found"
// @Failure 401 {object} types.MessageResponse "Unauthorized"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /posts/{post_id}/likes [post]
// @Security ApiKeyAuth
func (svc *WebService) LikePost(ctx *gin.Context) {
	// Check authorization
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Check URL params
	postId, err := strconv.Atoi(ctx.Param("post_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	// Call LikePost service
	resp, err := svc.AuthenticateAndPostClient.LikePost(ctx, &pb_aap.LikePostRequest{
		UserId: int64(userId),
		PostId: int64(postId),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.LikePostResponse_POST_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	} else if resp.GetStatus() == pb_aap.LikePostResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.LikePostResponse_OK {
		ctx.JSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// GetS3PresignedUrl godoc
// @Summary Get a presigned S3 URL for file upload
// @Description Get a presigned URL to directly upload a file to S3
// @Tags posts
// @Accept json
// @Produce json
// @Param request body types.GetS3PresignedUrlRequest true "File details"
// @Success 200 {object} types.GetS3PresignedUrlResponse "Presigned URL"
// @Failure 400 {object} types.MessageResponse "Validation error"
// @Failure 401 {object} types.MessageResponse "Unauthorized"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /posts/url [get]
// @Security ApiKeyAuth
func (svc *WebService) GetS3PresignedUrl(ctx *gin.Context) {
	// Check authorization
	_, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Validate request
	var jsonRequest types.GetS3PresignedUrlRequest
	err = ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	// TODO: Implement S3 presigned URL generation
	// This would typically involve using the AWS SDK to generate a presigned URL
	// For now, we'll return a placeholder response

	expirationTime := time.Now().Add(15 * time.Minute)
	ctx.JSON(http.StatusOK, types.GetS3PresignedUrlResponse{
		URL:            "https://s3.example.com/bucket/path/" + jsonRequest.FileName,
		ExpirationTime: expirationTime.Format(time.RFC3339),
	})
}
