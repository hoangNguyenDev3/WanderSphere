package authpost

import (
	"context"

	"github.com/hoangNguyenDev3/WanderSphere/internal/auth"
	"github.com/hoangNguyenDev3/WanderSphere/internal/models"
	authpost "github.com/hoangNguyenDev3/WanderSphere/pkg/types/proto/pb/authpost"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AuthenticationAndPostService struct {
	authpost.UnimplementedAuthenticationAndPostServer
	db *gorm.DB
}

func (s *AuthenticationAndPostService) CheckUserAuthentication(ctx context.Context, info *authpost.UserInfo) (*authpost.UserResult, error) {
	userModel := models.User{}
	err := s.db.Raw("select * from user where user_name = ?", info.GetUserName()).Scan(&userModel).Error
	if err != nil {
		return &authpost.UserResult{}, err
	}

	if userModel.ID == 0 {
		return &authpost.UserResult{}, status.Errorf(codes.NotFound, "User not found")
	}

	err = auth.CheckPasswordHash(userModel.HashedPassword, info.GetUserPassword())
	if err != nil {
		return &authpost.UserResult{}, status.Errorf(codes.Unauthenticated, "Invalid password")
	}

	return &authpost.UserResult{
		Status: 1,
		Info: &authpost.UserDetailInfo{
			UserName:  userModel.UserName,
			Email:     userModel.Email,
			FirstName: userModel.FirstName,
			LastName:  userModel.LastName,
		},
	}, nil
}

func (s *AuthenticationAndPostService) CreateUser(ctx context.Context, info *authpost.UserDetailInfo) (*authpost.UserResult, error) {
	// check if user already exists
	userModel := models.User{}
	err := s.db.Raw("select * from user where user_name = ?", info.GetUserName()).Scan(&userModel).Error
	if err != nil {
		return &authpost.UserResult{}, err
	}

	if userModel.ID != 0 {
		return &authpost.UserResult{}, status.Errorf(codes.AlreadyExists, "User already exists")
	}

	salt, err := auth.GenerateRandomSalt()
	if err != nil {
		return &authpost.UserResult{}, err
	}

	hashedPassword, err := auth.HashPassword(info.GetUserPassword(), salt)
	if err != nil {
		return &authpost.UserResult{}, err
	}

	err = s.db.Raw("insert into user (id, first_name, last_name, dob, email, user_name, hashed_password, salt) values (null, ?, ?, ?, ?, ?, ?, ?)", info.GetFirstName(), info.GetLastName(), info.GetDob(), info.GetEmail(), info.GetUserName(), hashedPassword, string(salt)).Error
	if err != nil {
		return &authpost.UserResult{}, err
	}

	return &authpost.UserResult{
		Status: 1,
		Info:   info,
	}, nil
}

func (s *AuthenticationAndPostService) EditUser(ctx context.Context, info *authpost.UserDetailInfo) (*authpost.UserResult, error) {
	// check if user exists
	userModel := models.User{}
	err := s.db.Raw("select * from user where user_name = ?", info.GetUserName()).Scan(&userModel).Error
	if err != nil {
		return &authpost.UserResult{}, err
	}

	if userModel.ID == 0 {
		return &authpost.UserResult{}, status.Errorf(codes.NotFound, "User not found")
	}

	// check if user is the same user
	if userModel.ID != info.GetUserId() {
		return &authpost.UserResult{}, status.Errorf(codes.PermissionDenied, "User is not the same")
	}

	// update user password
	if info.GetUserPassword() != "" {
		salt, err := auth.GenerateRandomSalt()
		if err != nil {
			return &authpost.UserResult{}, err
		}

		hashed_password, err := auth.HashPassword(info.GetUserPassword(), salt)
		if err != nil {
			return &authpost.UserResult{}, err
		}

		err = s.db.Exec("update user set hashed_password = ?, salt = ? where id = ?", hashed_password, salt, info.GetUserId()).Error
		if err != nil {
			return &authpost.UserResult{}, err
		}
	}

	// Edit first_name
	if info.GetFirstName() != "" {
		err = s.db.Exec("update user set first_name = ? where id = ?", info.GetFirstName(), info.GetUserId()).Error
		if err != nil {
			return &authpost.UserResult{}, err
		}
	}

	// Edit last_name
	if info.GetLastName() != "" {
		err = s.db.Exec("update user set last_name = ? where id = ?", info.GetLastName(), info.GetUserId()).Error
		if err != nil {
			return &authpost.UserResult{}, err
		}
	}

	// Edit dob
	if info.GetDob() != 0 {
		err = s.db.Exec("update user set dob = FROM_UNIXTIME(?) where id = ?", info.GetDob(), info.GetUserId()).Error
		if err != nil {
			return &authpost.UserResult{}, err
		}
	}

	// Edit email
	if info.GetEmail() != "" {
		err = s.db.Exec("update user set email = ? where id = ?", info.GetEmail(), info.GetUserId()).Error
		if err != nil {
			return &authpost.UserResult{}, err
		}
	}

	// Edit user_name
	if info.GetUserName() != "" {
		err = s.db.Exec("update user set username = ? where id = ?", info.GetUserName(), info.GetUserId()).Error
		if err != nil {
			return &authpost.UserResult{}, err
		}
	}

	return &authpost.UserResult{
		Status: 1,
		Info:   info,
	}, nil
}

func (s *AuthenticationAndPostService) GetUserFollower(ctx context.Context, info *authpost.UserInfo) (*authpost.UserFollower, error) {
	userModel := models.User{}
	err := s.db.Raw("select * from user where user_name = ?", info.GetUserName()).Scan(&userModel).Error
	if err != nil {
		return &authpost.UserFollower{}, err
	}

	// If the user exists, return the followers
	var followers []models.User
	err = s.db.Raw("select follower_id from user_user where user_id = ?", userModel.ID).Scan(&followers).Error
	if err != nil {
		return &authpost.UserFollower{}, err
	}

	returnUserFolower := authpost.UserFollower{}
	for _, follower := range followers {
		followerInfo := authpost.UserInfo{UserId: follower.ID, UserName: follower.UserName}
		returnUserFolower.Followers = append(returnUserFolower.Followers, &followerInfo)
	}

	return &returnUserFolower, nil
}

func (s *AuthenticationAndPostService) GetPostDetail(ctx context.Context, info *authpost.GetPostRequest) (*authpost.Post, error) {
	postModel := models.Post{}
	err := s.db.Raw("select * from post where id = ?", info.GetPostId()).Scan(&postModel).Error
	if err != nil {
		return &authpost.Post{}, err
	}

	returnPost := authpost.Post{
		PostId:           postModel.ID,
		UserId:           postModel.UserID,
		ContentText:      postModel.ContentText,
		ContentImagePath: postModel.ContentImagePath,
		CreatedAt:        postModel.CreatedAt.Unix(),
		Visible:          postModel.Visible,
	}
	return &returnPost, nil
}

func NewAuthenticationAndPostService(db *gorm.DB) *AuthenticationAndPostService {
	return &AuthenticationAndPostService{db: db}
}
