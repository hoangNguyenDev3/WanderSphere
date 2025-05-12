package service

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/hoangNguyenDev3/WanderSphere/backend/docs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb_aap "github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost"
)

// CheckUserAuthentication godoc
// @Summary Authenticate a user
// @Description Login with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param login body types.LoginRequest true "Login credentials"
// @Success 200 {object} types.LoginResponse "Login successful"
// @Failure 400 {object} types.ErrorResponse "Validation error or authentication failed"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /users/login [post]
func (svc *WebService) CheckUserAuthentication(ctx *gin.Context) {
	// Validate request
	var jsonRequest types.LoginRequest
	err := ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Call CheckUserAuthentication service
	authentication, err := svc.AuthenticateAndPostClient.CheckUserAuthentication(ctx, &pb_aap.CheckUserAuthenticationRequest{
		UserName:     jsonRequest.UserName,
		UserPassword: jsonRequest.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	if authentication.GetStatus() == pb_aap.CheckUserAuthenticationResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "auth_error",
			Message: "wrong username or password",
			Code:    http.StatusBadRequest,
		})
		return
	} else if authentication.GetStatus() == pb_aap.CheckUserAuthenticationResponse_WRONG_PASSWORD {
		ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "auth_error",
			Message: "wrong username or password",
			Code:    http.StatusBadRequest,
		})
		return
	} else if authentication.GetStatus() == pb_aap.CheckUserAuthenticationResponse_OK {
		// Set a sessionId for this session
		sessionId := uuid.New().String()

		// Get session configuration from config
		var expirationTime time.Duration

		// Use config value if available, otherwise use a reasonable default (24 hours)
		if svc.Config != nil && svc.Config.Auth.Session.ExpirationMinutes > 0 {
			expirationTime = time.Minute * time.Duration(svc.Config.Auth.Session.ExpirationMinutes)
		} else {
			expirationTime = time.Hour * 24 // Default to 24 hours
		}

		// Save current sessionID and expiration time in Redis
		err = svc.RedisClient.Set(svc.RedisClient.Context(), sessionId, authentication.GetUserId(), expirationTime).Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
				Error:   "session_error",
				Message: "Failed to create session: " + err.Error(),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Set cookie configuration
		cookieName := "session_id" // Default
		if svc.Config != nil && svc.Config.Auth.Session.CookieName != "" {
			cookieName = svc.Config.Auth.Session.CookieName
		}

		// Default to secure settings unless explicitly configured otherwise
		secure := true
		httpOnly := true
		sameSite := http.SameSiteStrictMode

		if svc.Config != nil {
			// Only override defaults if explicitly set in config
			if svc.Config.Auth.Session.Secure == false {
				secure = false
			}
			if svc.Config.Auth.Session.HTTPOnly == false {
				httpOnly = false
			}

			if svc.Config.Auth.Session.SameSite == "lax" {
				sameSite = http.SameSiteLaxMode
			} else if svc.Config.Auth.Session.SameSite == "none" {
				sameSite = http.SameSiteNoneMode
			}
		}

		// Set sessionID cookie with secure settings
		maxAge := int(expirationTime.Seconds())
		ctx.SetSameSite(sameSite)
		ctx.SetCookie(cookieName, sessionId, maxAge, "/", "", secure, httpOnly)

		// Get user details to include in response
		userInfo, err := svc.AuthenticateAndPostClient.GetUserDetailInfo(ctx, &pb_aap.GetUserDetailInfoRequest{
			UserId: authentication.GetUserId(),
		})

		if err != nil {
			ctx.JSON(http.StatusOK, types.MessageResponse{
				Message: "Login successful",
				Status:  "success",
			})
			return
		}

		// Return user details along with success message
		ctx.JSON(http.StatusOK, types.LoginResponse{
			Message: "Login successful",
			User: types.UserDetailInfo{
				UserID:         userInfo.GetUser().GetUserId(),
				UserName:       userInfo.GetUser().GetUserName(),
				FirstName:      userInfo.GetUser().GetFirstName(),
				LastName:       userInfo.GetUser().GetLastName(),
				DateOfBirth:    userInfo.GetUser().GetDateOfBirth().AsTime().Format(time.DateOnly),
				Email:          userInfo.GetUser().GetEmail(),
				ProfilePicture: "",
				CoverPicture:   "",
			},
		})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "unknown_error",
			Message: "An unexpected error occurred",
			Code:    http.StatusInternalServerError,
		})
		return
	}
}

// CreateUser godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body types.CreateUserRequest true "User registration information"
// @Success 200 {object} types.MessageResponse "User created successfully"
// @Failure 400 {object} types.MessageResponse "Validation error or user already exists"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /users/signup [post]
func (svc *WebService) CreateUser(ctx *gin.Context) {
	// Validate request
	var jsonRequest types.CreateUserRequest
	err := ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	// Call CreateUser service
	dob, err := time.Parse(time.DateOnly, jsonRequest.DateOfBirth)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	resp, err := svc.AuthenticateAndPostClient.CreateUser(ctx, &pb_aap.CreateUserRequest{
		UserName:     jsonRequest.UserName,
		UserPassword: jsonRequest.Password,
		FirstName:    jsonRequest.FirstName,
		LastName:     jsonRequest.LastName,
		DateOfBirth:  timestamppb.New(dob),
		Email:        jsonRequest.Email,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.CreateUserResponse_USERNAME_EXISTED {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "username existed"})
		return
	} else if resp.GetStatus() == pb_aap.CreateUserResponse_EMAIL_EXISTED {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "email existed"})
		return
	} else if resp.GetStatus() == pb_aap.CreateUserResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// EditUser godoc
// @Summary Edit user profile
// @Description Update user profile information
// @Tags users
// @Accept json
// @Produce json
// @Param user body types.EditUserRequest true "User information to update"
// @Success 200 {object} types.MessageResponse "User updated successfully"
// @Failure 400 {object} types.MessageResponse "Validation error or user not found"
// @Failure 401 {object} types.MessageResponse "Unauthorized"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /users/edit [post]
// @Security ApiKeyAuth
func (svc *WebService) EditUser(ctx *gin.Context) {
	// Check authorization
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Validate request
	var jsonRequest types.EditUserRequest
	err = ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	var password *string
	if jsonRequest.Password != "" {
		password = &jsonRequest.Password
	}
	var firstName *string
	if jsonRequest.FirstName != "" {
		firstName = &jsonRequest.FirstName
	}
	var lastName *string
	if jsonRequest.LastName != "" {
		lastName = &jsonRequest.LastName
	}
	var dateOfBirth *timestamppb.Timestamp
	if jsonRequest.DateOfBirth != "" {
		dob, err := time.Parse(time.DateOnly, jsonRequest.DateOfBirth)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
			return
		}
		dateOfBirth = timestamppb.New(dob)
	}

	// These are declared but not used currently until the proto is updated
	_ = jsonRequest.ProfilePicture
	_ = jsonRequest.CoverPicture

	// Call EditUser service
	resp, err := svc.AuthenticateAndPostClient.EditUser(ctx, &pb_aap.EditUserRequest{
		UserId:       int64(userId),
		UserPassword: password,
		FirstName:    firstName,
		LastName:     lastName,
		DateOfBirth:  dateOfBirth,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.EditUserResponse_USER_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.EditUserResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// GetUserDetailInfo godoc
// @Summary Get user details
// @Description Get detailed information about a user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} types.UserDetailInfoResponse "User details"
// @Failure 400 {object} types.MessageResponse "Invalid user ID or user not found"
// @Failure 500 {object} types.MessageResponse "Internal server error"
// @Router /users/{user_id} [get]
func (svc *WebService) GetUserDetailInfo(ctx *gin.Context) {
	// Check URL params
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	// Call gprc service
	resp, err := svc.AuthenticateAndPostClient.GetUserDetailInfo(ctx, &pb_aap.GetUserDetailInfoRequest{
		UserId: int64(userId),
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.GetUserDetailInfoResponse_USER_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.GetUserDetailInfoResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.UserDetailInfoResponse{
			UserID:         resp.GetUser().GetUserId(),
			UserName:       resp.GetUser().GetUserName(),
			FirstName:      resp.GetUser().GetFirstName(),
			LastName:       resp.GetUser().GetLastName(),
			DateOfBirth:    resp.GetUser().GetDateOfBirth().AsTime().Format(time.DateOnly),
			Email:          resp.GetUser().GetEmail(),
			ProfilePicture: "",
			CoverPicture:   "",
		})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

func (svc *WebService) checkSessionAuthentication(ctx *gin.Context) (sessionId string, userId int, err error) {
	// Get cookie name from config
	cookieName := "session_id" // Default
	if svc.Config != nil && svc.Config.Auth.Session.CookieName != "" {
		cookieName = svc.Config.Auth.Session.CookieName
	}

	// Try to get session cookie
	sessionId, err = ctx.Cookie(cookieName)
	if err != nil {
		svc.Logger.Debug("Session cookie not found",
			zap.String("cookie_name", cookieName),
			zap.Error(err))
		return "", 0, err
	}

	// Validate session in Redis
	val, err := svc.RedisClient.Get(svc.RedisClient.Context(), sessionId).Result()
	if err != nil {
		svc.Logger.Debug("Session not found in Redis or expired",
			zap.String("session_id", sessionId),
			zap.Error(err))
		return "", 0, err
	}

	// Parse user ID from session
	userId, err = strconv.Atoi(val)
	if err != nil {
		svc.Logger.Warn("Invalid user ID in session",
			zap.String("session_id", sessionId),
			zap.String("value", val),
			zap.Error(err))
		return "", 0, err
	}

	return sessionId, userId, nil
}
