package authpost

import (
	"context"
	"errors"

	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/auth"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
	pb "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// CheckUserAuthentication verifies user credentials
func (s *AuthenticateAndPostService) CheckUserAuthentication(ctx context.Context, req *pb.CheckUserAuthenticationRequest) (*pb.CheckUserAuthenticationResponse, error) {
	s.logger.Debug("CheckUserAuthentication request received", zap.String("username", req.UserName))

	// Find user by username (data layer interaction)
	var user types.User
	result := s.db.Where(&types.User{UserName: req.UserName}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &pb.CheckUserAuthenticationResponse{
			Status: pb.CheckUserAuthenticationResponse_USER_NOT_FOUND,
		}, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	// Verify password with salt (business logic)
	err := auth.CheckPasswordHash(user.HashedPassword, req.UserPassword, user.Salt)
	if err != nil {
		return &pb.CheckUserAuthenticationResponse{
			Status: pb.CheckUserAuthenticationResponse_WRONG_PASSWORD,
		}, nil
	}

	return &pb.CheckUserAuthenticationResponse{
		Status: pb.CheckUserAuthenticationResponse_OK,
		UserId: user.ID,
	}, nil
}

// CreateUser creates a new user
func (s *AuthenticateAndPostService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	s.logger.Debug("CreateUser request received", zap.String("username", req.UserName))

	// Check if username already exists
	exists, _ := s.findUserByUserName(req.UserName)
	if exists {
		return &pb.CreateUserResponse{
			Status: pb.CreateUserResponse_USERNAME_EXISTED,
		}, nil
	}

	// Generate salt and hash password
	salt, err := auth.GenerateRandomSalt()
	if err != nil {
		s.logger.Error("Error generating salt", zap.Error(err))
		return nil, err
	}

	hashedPassword, err := auth.HashPassword(req.UserPassword, salt)
	if err != nil {
		s.logger.Error("Error hashing password", zap.Error(err))
		return nil, err
	}

	// Create user
	user := types.User{
		UserName:       req.UserName,
		HashedPassword: hashedPassword,
		Salt:           salt,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		DateOfBirth:    req.DateOfBirth.AsTime(),
		Email:          req.Email,
	}

	result := s.db.Create(&user)
	if result.Error != nil {
		s.logger.Error("Error creating user", zap.Error(result.Error))
		return nil, result.Error
	}

	return &pb.CreateUserResponse{
		Status: pb.CreateUserResponse_OK,
		UserId: user.ID,
	}, nil
}

// GetUserDetailInfo retrieves user details
func (s *AuthenticateAndPostService) GetUserDetailInfo(ctx context.Context, req *pb.GetUserDetailInfoRequest) (*pb.GetUserDetailInfoResponse, error) {
	s.logger.Debug("GetUserDetailInfo request received", zap.Int64("user_id", req.UserId))

	// Find user
	exists, user := s.findUserById(req.UserId)
	if !exists {
		return &pb.GetUserDetailInfoResponse{
			Status: pb.GetUserDetailInfoResponse_USER_NOT_FOUND,
		}, nil
	}

	// Return user details
	return &pb.GetUserDetailInfoResponse{
		Status: pb.GetUserDetailInfoResponse_OK,
		User: &pb.UserDetailInfo{
			UserId:      user.ID,
			UserName:    user.UserName,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			DateOfBirth: timestamppb.New(user.DateOfBirth),
			Email:       user.Email,
		},
	}, nil
}

// EditUser updates user information
func (s *AuthenticateAndPostService) EditUser(ctx context.Context, req *pb.EditUserRequest) (*pb.EditUserResponse, error) {
	s.logger.Debug("EditUser called", zap.Int64("user_id", req.UserId))

	// Find user
	exists, user := s.findUserById(req.UserId)
	if !exists {
		return &pb.EditUserResponse{
			Status: pb.EditUserResponse_USER_NOT_FOUND,
		}, nil
	}

	// Update fields
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.DateOfBirth != nil {
		user.DateOfBirth = req.DateOfBirth.AsTime()
	}
	if req.UserPassword != nil {
		salt, err := auth.GenerateRandomSalt()
		if err != nil {
			s.logger.Error("Error generating salt", zap.Error(err))
			return nil, err
		}
		hashedPassword, err := auth.HashPassword(*req.UserPassword, salt)
		if err != nil {
			s.logger.Error("Error hashing password", zap.Error(err))
			return nil, err
		}
		user.Salt = salt
		user.HashedPassword = hashedPassword
	}

	// Save changes
	result := s.db.Save(&user)
	if result.Error != nil {
		s.logger.Error("Error updating user", zap.Error(result.Error))
		return nil, result.Error
	}

	return &pb.EditUserResponse{
		Status: pb.EditUserResponse_OK,
	}, nil
}
