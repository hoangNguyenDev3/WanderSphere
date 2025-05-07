package authpost

import (
	"context"
	"errors"

	"github.com/hoangNguyenDev3/WanderSphere/configs"
	"github.com/hoangNguyenDev3/WanderSphere/internal/auth"
	"github.com/hoangNguyenDev3/WanderSphere/internal/models"
	authpost "github.com/hoangNguyenDev3/WanderSphere/pkg/types/proto/pb/authpost"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthenticationAndPostService struct {
	authpost.UnimplementedAuthenticationAndPostServer
	db *gorm.DB
}

func (s *AuthenticationAndPostService) CreateUser(ctx context.Context, info *authpost.UserDetailInfo) (*authpost.UserResult, error) {
	existed, _ := s.checkUserName(info.GetUserName())
	if existed {
		return nil, errors.New("user already exist")
	}

	salt, err := auth.GenerateRandomSalt()
	if err != nil {
		return nil, err
	}

	hashed_password, err := auth.HashPassword(info.GetUserPassword(), salt)
	if err != nil {
		return nil, err
	}

	err = s.db.Exec(
		"insert into \"user\" (id, hashed_password, salt, first_name, last_name, dob, email, user_name) values (default, ?, ?, ?, ?, to_timestamp(?), ?, ?)",
		hashed_password,
		salt,
		info.GetFirstName(),
		info.GetLastName(),
		info.GetDob(),
		info.GetEmail(),
		info.GetUserName(),
	).Error
	if err != nil {
		return nil, err
	}

	// Return the necessary user information
	_, userModel := s.checkUserName(info.GetUserName())
	return s.NewUserResult(userModel), nil
}

func (s *AuthenticationAndPostService) CheckUserAuthentication(ctx context.Context, info *authpost.UserInfo) (*authpost.UserResult, error) {
	// Find user in the database using the provided information
	existed, userModel := s.checkUserName(info.GetUserName())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Check password matching
	err := auth.CheckPasswordHash(userModel.HashedPassword, info.GetUserPassword())
	if err != nil {
		return nil, err
	}

	// Return the necessary user information
	return s.NewUserResult(userModel), nil
}

func (s *AuthenticationAndPostService) EditUser(ctx context.Context, info *authpost.UserDetailInfo) (*authpost.UserResult, error) {
	// Check if the userId which is changing information exists in the database
	existed, userModel := s.checkUserName(info.GetUserName())
	if !existed {
		return &authpost.UserResult{}, errors.New("user does not exist")
	}

	// If the user exists, edit the information and return
	var err error

	// Edit password
	if info.GetUserPassword() != "" {
		salt, err := auth.GenerateRandomSalt()
		if err != nil {
			return &authpost.UserResult{}, err
		}

		hashed_password, err := auth.HashPassword(info.GetUserPassword(), salt)
		if err != nil {
			return &authpost.UserResult{}, err
		}

		err = s.db.Exec("update \"user\" set hashed_password = ?, salt = ? where id = ?", hashed_password, salt, userModel.ID).Error
		if err != nil {
			return &authpost.UserResult{}, err
		}
	}

	// Edit first_name
	if info.GetFirstName() != "" {
		err = s.db.Exec("update \"user\" set first_name = ? where id = ?", info.GetFirstName(), userModel.ID).Error
		if err != nil {
			return &authpost.UserResult{}, err
		}
	}

	// Edit last_name
	if info.GetLastName() != "" {
		err = s.db.Exec("update \"user\" set last_name = ? where id = ?", info.GetLastName(), userModel.ID).Error
		if err != nil {
			return &authpost.UserResult{}, err
		}
	}

	// Edit dob
	if info.GetDob() >= -2208988800 { // 1900-01-01
		err = s.db.Exec("update \"user\" set dob = to_timestamp(?) where id = ?", info.GetDob(), userModel.ID).Error
		if err != nil {
			return &authpost.UserResult{}, err
		}
	}

	// Edit email
	if info.GetEmail() != "" {
		err = s.db.Exec("update \"user\" set email = ? where id = ?", info.GetEmail(), userModel.ID).Error
		if err != nil {
			return &authpost.UserResult{}, err
		}
	}

	// Return the necessary user information
	_, userModel = s.checkUserName(info.GetUserName())
	return s.NewUserResult(userModel), nil
}

func (s *AuthenticationAndPostService) GetUserFollower(ctx context.Context, info *authpost.UserInfo) (*authpost.UserFollower, error) {
	// Check if the user exists
	userModel := models.User{}
	err := s.db.Raw("select * from \"user\" where user_name = ?", info.GetUserName()).Scan(&userModel).Error
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

func (s *AuthenticationAndPostService) GetPostDetail(ctx context.Context, request *authpost.GetPostRequest) (*authpost.Post, error) {
	// Check if the post exists
	postModel := models.Post{}
	err := s.db.Raw("select * from post where id = ?", request.GetPostId()).Scan(&postModel).Error
	if err != nil {
		return &authpost.Post{}, err
	}

	// If the post exists, return the post
	returnPost := authpost.Post{
		PostId:           postModel.ID,
		UserId:           postModel.UserID,
		ContentText:      postModel.ContentText,
		ContentImagePath: postModel.ContentImagePath,
		Visible:          postModel.Visible,
		CreatedAt:        postModel.CreatedAt.Unix(),
	}
	return &returnPost, nil
}

func NewAuthenticationAndPostService(cfg *configs.AuthenticateAndPostConfig) (*AuthenticationAndPostService, error) {
	// Connect to database
	postgresConfig := postgres.Config{
		DSN: cfg.Postgres.DSN,
	}
	db, err := gorm.Open(postgres.New(postgresConfig), &gorm.Config{})
	if err != nil {
		return &AuthenticationAndPostService{}, err
	}

	return &AuthenticationAndPostService{db: db}, err
}

// checkUserName checks if an user with provided username exists in database
func (s *AuthenticationAndPostService) checkUserName(username string) (bool, models.User) {
	var userModel = models.User{}
	s.db.Raw("select * from \"user\" where user_name = ?", username).Scan(&userModel)

	if userModel.ID == 0 {
		return false, models.User{}
	}
	return true, userModel
}

func (s *AuthenticationAndPostService) NewUserResult(userModel models.User) *authpost.UserResult {
	return &authpost.UserResult{
		Status: 0, // OK status
		Info: &authpost.UserDetailInfo{
			UserId:       userModel.ID,
			UserName:     userModel.UserName,
			UserPassword: "",
			FirstName:    userModel.FirstName,
			LastName:     userModel.LastName,
			Dob:          userModel.DOB.Unix(),
			Email:        userModel.Email,
		},
	}
}
