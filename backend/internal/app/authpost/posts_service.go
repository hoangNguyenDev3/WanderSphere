package authpost

import (
	"context"
	"time"

	"strings"

	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
	pb_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	pb_nfp "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed_publishing"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *AuthenticateAndPostService) CreatePost(ctx context.Context, info *pb_aap.CreatePostRequest) (*pb_aap.CreatePostResponse, error) {
	a.logger.Debug("start creating post")
	defer a.logger.Debug("end creating post")

	exist, _ := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.CreatePostResponse{Status: pb_aap.CreatePostResponse_USER_NOT_FOUND}, nil
	}

	// Process image paths
	var contentImagePath string
	if len(info.GetContentImagePath()) > 0 {
		// For consistency with the database schema, we're joining multiple paths
		// with a space separator. When retrieving, we'll split this string back to a slice.
		contentImagePath = strings.Join(info.GetContentImagePath(), " ")
		a.logger.Debug("Using image paths", zap.String("paths", contentImagePath))
	}

	newPost := types.Post{
		UserID:           info.GetUserId(),
		ContentText:      info.GetContentText(),
		ContentImagePath: contentImagePath,
	}

	// Handle visibility - if not visible, set DeletedAt to current time
	if !info.GetVisible() {
		now := time.Now()
		newPost.DeletedAt = &now
	}

	result := a.db.Create(&newPost)
	if result.Error != nil {
		a.logger.Error("Error creating post", zap.Error(result.Error))
		return nil, result.Error
	}

	// Send user_id and post_id to NewsfeedPublishingClient to announce to followers
	if a.nfPubClient != nil {
		_, err := a.nfPubClient.PublishPost(ctx, &pb_nfp.PublishPostRequest{
			UserId: newPost.UserID,
			PostId: int64(newPost.ID),
		})
		if err != nil {
			a.logger.Error("Error publishing post to newsfeed", zap.Error(err))
			// Continue anyway, as the post is created - async event can be retried
		}
	}

	responsePostId := int64(newPost.ID)

	return &pb_aap.CreatePostResponse{
		Status: pb_aap.CreatePostResponse_OK,
		PostId: responsePostId,
	}, nil
}

func (a *AuthenticateAndPostService) EditPost(ctx context.Context, info *pb_aap.EditPostRequest) (*pb_aap.EditPostResponse, error) {
	a.logger.Debug("start editing post", zap.Int64("user_id", info.GetUserId()), zap.Int64("post_id", info.GetPostId()))
	defer a.logger.Debug("end editing post")

	exist, user := a.findUserById(info.GetUserId())
	if !exist {
		a.logger.Debug("user not found for edit post", zap.Int64("user_id", info.GetUserId()))
		return &pb_aap.EditPostResponse{Status: pb_aap.EditPostResponse_USER_NOT_FOUND}, nil
	}

	// Use the consistent findPostById method instead of direct DB query
	exist, post := a.findPostById(info.GetPostId())
	if !exist {
		a.logger.Debug("post not found for edit", zap.Int64("post_id", info.GetPostId()))
		return &pb_aap.EditPostResponse{Status: pb_aap.EditPostResponse_POST_NOT_FOUND}, nil
	}

	// Compare user ID as int64 to avoid type conversion issues
	if int64(user.ID) != post.UserID {
		a.logger.Debug("user not allowed to edit post",
			zap.Int64("user_id", int64(user.ID)),
			zap.Int64("post_owner_id", post.UserID),
			zap.Int64("post_id", info.GetPostId()))
		return &pb_aap.EditPostResponse{Status: pb_aap.EditPostResponse_NOT_ALLOWED}, nil
	}

	// Apply updates
	if info.ContentText != nil {
		post.ContentText = info.GetContentText()
		a.logger.Debug("updating post content text", zap.String("new_content", post.ContentText))
	}
	if info.ContentImagePath != nil {
		post.ContentImagePath = info.GetContentImagePath()
		a.logger.Debug("updating post image path", zap.String("new_path", post.ContentImagePath))
	}
	if info.Visible != nil {
		if info.GetVisible() {
			post.DeletedAt = nil // Not deleted - visible
			a.logger.Debug("making post visible")
		} else {
			now := time.Now()
			post.DeletedAt = &now // Deleted - not visible
			a.logger.Debug("making post invisible")
		}
	}

	err := a.db.Save(&post).Error
	if err != nil {
		a.logger.Error("error saving edited post", zap.Error(err), zap.Int64("post_id", info.GetPostId()))
		return nil, err
	}

	a.logger.Info("post edited successfully", zap.Int64("post_id", info.GetPostId()), zap.Int64("user_id", info.GetUserId()))
	return &pb_aap.EditPostResponse{
		Status: pb_aap.EditPostResponse_OK,
	}, nil
}

func (a *AuthenticateAndPostService) DeletePost(ctx context.Context, info *pb_aap.DeletePostRequest) (*pb_aap.DeletePostResponse, error) {
	a.logger.Debug("start deleting post")
	defer a.logger.Debug("end deleting post")

	exist, user := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.DeletePostResponse{Status: pb_aap.DeletePostResponse_USER_NOT_FOUND}, nil
	}
	exist, post := a.findPostById(info.GetPostId())
	if !exist {
		return &pb_aap.DeletePostResponse{Status: pb_aap.DeletePostResponse_POST_NOT_FOUND}, nil
	}
	// Compare user ID as int64 to avoid type conversion issues
	if int64(user.ID) != post.UserID {
		return &pb_aap.DeletePostResponse{Status: pb_aap.DeletePostResponse_NOT_ALLOWED}, nil
	}

	// Start a transaction to ensure data consistency
	tx := a.db.Begin()
	if tx.Error != nil {
		a.logger.Error("Error starting transaction for post deletion", zap.Error(tx.Error))
		return nil, tx.Error
	}

	a.logger.Debug("Starting cascading delete for post", zap.Int64("post_id", info.GetPostId()))

	// First delete all comments for this post
	commentsResult := tx.Where("post_id = ?", info.GetPostId()).Delete(&types.Comment{})
	if commentsResult.Error != nil {
		tx.Rollback()
		a.logger.Error("Error deleting comments for post", zap.Error(commentsResult.Error))
		return nil, commentsResult.Error
	}
	a.logger.Debug("Deleted comments", zap.Int64("count", commentsResult.RowsAffected))

	// Then delete all likes for this post
	likesResult := tx.Where("post_id = ?", info.GetPostId()).Delete(&types.Like{})
	if likesResult.Error != nil {
		tx.Rollback()
		a.logger.Error("Error deleting likes for post", zap.Error(likesResult.Error))
		return nil, likesResult.Error
	}
	a.logger.Debug("Deleted likes", zap.Int64("count", likesResult.RowsAffected))

	// Finally delete the post itself using the transaction
	postResult := tx.Where("id = ?", info.GetPostId()).Delete(&types.Post{})
	if postResult.Error != nil {
		tx.Rollback()
		a.logger.Error("Error deleting post", zap.Error(postResult.Error))
		return nil, postResult.Error
	}
	a.logger.Debug("Deleted post", zap.Int64("count", postResult.RowsAffected))

	// Commit the transaction
	err := tx.Commit().Error
	if err != nil {
		a.logger.Error("Error committing post deletion transaction", zap.Error(err))
		return nil, err
	}

	a.logger.Info("Post deleted successfully", zap.Int64("post_id", info.GetPostId()))

	// Note: In a real implementation, we would need to:
	// 1. Add a dependency on the newsfeed service client
	// 2. Call a method to invalidate the cache for this post
	// Something like: a.newsfeedClient.RemovePostFromNewsfeed(ctx, post.ID)
	// But for now we'll just log that we would do this
	a.logger.Info("Post deleted, should invalidate newsfeed cache",
		zap.Int64("post_id", int64(post.ID)),
		zap.Int64("user_id", post.UserID))

	return &pb_aap.DeletePostResponse{
		Status: pb_aap.DeletePostResponse_OK,
	}, nil
}

func (a *AuthenticateAndPostService) GetPostDetailInfo(ctx context.Context, info *pb_aap.GetPostDetailInfoRequest) (*pb_aap.GetPostDetailInfoResponse, error) {
	a.logger.Debug("start getting post")
	defer a.logger.Debug("end getting post")

	exist, _ := a.findPostById(info.GetPostId())
	if !exist {
		return &pb_aap.GetPostDetailInfoResponse{Status: pb_aap.GetPostDetailInfoResponse_POST_NOT_FOUND}, nil
	}

	var post types.Post
	result := a.db.Preload("Comments").Preload("LikedUsers").First(&post, info.GetPostId())
	if result.Error != nil {
		return nil, result.Error
	}

	var comments []*pb_aap.Comment
	for i := range post.Comments {
		comments = append(comments, &pb_aap.Comment{
			CommentId:   int64(post.Comments[i].ID),
			UserId:      post.Comments[i].UserID,
			ContentText: post.Comments[i].ContentText,
			PostId:      int64(post.ID),
		})
	}

	var likedUsers []*pb_aap.Like
	for i := range post.LikedUsers {
		likedUsers = append(likedUsers, &pb_aap.Like{
			UserId: int64(post.LikedUsers[i].ID),
			PostId: int64(post.ID),
		})
	}

	return &pb_aap.GetPostDetailInfoResponse{
		Status: pb_aap.GetPostDetailInfoResponse_OK,
		Post: &pb_aap.PostDetailInfo{
			PostId:           int64(post.ID),
			UserId:           post.UserID,
			ContentText:      post.ContentText,
			ContentImagePath: strings.Split(post.ContentImagePath, " "),
			CreatedAt:        timestamppb.New(post.CreatedAt),
			Comments:         comments,
			LikedUsers:       likedUsers,
		},
	}, nil
}

func (a *AuthenticateAndPostService) CommentPost(ctx context.Context, info *pb_aap.CommentPostRequest) (*pb_aap.CommentPostResponse, error) {
	a.logger.Debug("start commenting post")
	defer a.logger.Debug("end commenting post")

	exist, _ := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.CommentPostResponse{Status: pb_aap.CommentPostResponse_USER_NOT_FOUND}, nil
	}
	exist, _ = a.findPostById(info.GetPostId())
	if !exist {
		return &pb_aap.CommentPostResponse{Status: pb_aap.CommentPostResponse_POST_NOT_FOUND}, nil
	}

	var newComment = types.Comment{
		PostID:      info.GetPostId(),
		UserID:      info.GetUserId(),
		ContentText: info.GetContentText(),
	}
	err := a.db.Create(&newComment).Error
	if err != nil {
		return nil, err
	}

	return &pb_aap.CommentPostResponse{
		Status:    pb_aap.CommentPostResponse_OK,
		CommentId: int64(newComment.ID),
	}, nil
}

func (a *AuthenticateAndPostService) LikePost(ctx context.Context, info *pb_aap.LikePostRequest) (*pb_aap.LikePostResponse, error) {
	a.logger.Debug("start liking post")
	defer a.logger.Debug("end liking post")

	exist, user := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.LikePostResponse{Status: pb_aap.LikePostResponse_USER_NOT_FOUND}, nil
	}
	exist, _ = a.findPostById(info.GetPostId())
	if !exist {
		return &pb_aap.LikePostResponse{Status: pb_aap.LikePostResponse_POST_NOT_FOUND}, nil
	}

	var post types.Post
	err := a.db.Preload("LikedUsers").First(&post, info.GetPostId()).Error
	if err != nil {
		return nil, err
	}
	err = a.db.Model(&post).Association("LikedUsers").Append(&user)
	if err != nil {
		return nil, err
	}

	return &pb_aap.LikePostResponse{
		Status: pb_aap.LikePostResponse_OK,
	}, nil
}
