package authpost

import (
	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
	pb_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthenticateAndPostService struct {
	pb_aap.UnimplementedAuthenticateAndPostServer
	db *gorm.DB
}

func NewAuthenticateAndPostService(cfg *configs.AuthenticateAndPostConfig) (*AuthenticateAndPostService, error) {
	// Connect to database
	postgresConfig := postgres.Config{
		DSN: cfg.Postgres.DSN,
	}
	db, err := gorm.Open(postgres.New(postgresConfig), &gorm.Config{})
	if err != nil {
		return &AuthenticateAndPostService{}, err
	}

	return &AuthenticateAndPostService{db: db}, err
}

// checkUserName checks if an user with provided username exists in database
func (a *AuthenticateAndPostService) checkUserName(username string) (bool, types.User) {
	var userModel = types.User{}
	a.db.Raw("select * from user where user_name = ?", username).Scan(&userModel)

	if userModel.ID == 0 {
		return false, types.User{}
	}
	return true, userModel
}

// checkUserId checks if an user with provided userId exists in database
func (a *AuthenticateAndPostService) checkUserId(userId int64) (bool, types.User) {
	var userModel = types.User{}
	a.db.Raw("select * from user where id = ?", userId).Scan(&userModel)

	if userModel.ID == 0 {
		return false, types.User{}
	}
	return true, userModel
}

// checkPostId checks if an user with provided userId exists in database
func (a *AuthenticateAndPostService) checkPostId(postId int64) (bool, types.Post) {
	var postModel = types.Post{}
	a.db.Raw("select * from `post` where id = ?", postId).Scan(&postModel)

	if postModel.ID == 0 {
		return false, types.Post{}
	}
	return true, postModel
}

func (a *AuthenticateAndPostService) NewUserResult(userModel types.User) *pb_aap.UserResult {
	return &pb_aap.UserResult{
		Status: pb_aap.UserStatus_OK,
		Info: &pb_aap.UserDetailInfo{
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

func (a *AuthenticateAndPostService) NewActionResult(status pb_aap.ActionStatus) *pb_aap.ActionResult {
	return &pb_aap.ActionResult{Status: status}
}
