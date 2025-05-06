package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pb_authpost "github.com/hoangNguyenDev3/WanderSphere/pkg/types/proto/pb/authpost"
)

func AddUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("/users")

	userRouter.POST("/register", createUserHandler)
	userRouter.POST("/login", loginUserHandler)
	userRouter.POST("/edit", editUserHandler)
}

func createUserHandler(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	firstname := ctx.PostForm("firstname")
	lastname := ctx.PostForm("lastname")
	dob := ctx.PostForm("dob")
	email := ctx.PostForm("email")

	// Process the contents
	processedDOB, err := time.Parse(time.DateOnly, dob)
	if err != nil {
		ctx.IndentedJSON(
			http.StatusForbidden,
			gin.H{
				"message": "dob needs to have following format yyyy-mm-dd",
				"error":   fmt.Sprintf("create user failed: %v", err),
			},
		)
		return
	}

	_, err = authpostClient.CreateUser(
		context.Background(),
		&pb_authpost.UserDetailInfo{
			UserName:     username,
			UserPassword: password,
			FirstName:    firstname,
			LastName:     lastname,
			Dob:          processedDOB.Unix(),
			Email:        email,
		},
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("create user failed: %v", err)})
		return
	}

	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "create user succeeded!"})
}

func loginUserHandler(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	_, err := authpostClient.CheckUserAuthentication(context.Background(), &pb_authpost.UserInfo{UserId: 1, UserName: username, UserPassword: password})
	if err != nil {
		ctx.IndentedJSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("Login failed: %v", err)})
		return
	}

	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "Log in succeeded!"})
}

func editUserHandler(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	firstname := ctx.PostForm("firstname")
	lastname := ctx.PostForm("lastname")
	dob := ctx.PostForm("dob")
	email := ctx.PostForm("email")

	// Process the contents
	var dob_unix int64

	if dob != "" {
		processedDOB, err := time.Parse(time.DateOnly, dob)
		if err != nil {
			ctx.IndentedJSON(
				http.StatusForbidden,
				gin.H{
					"message": "dob needs to have following format yyyy-mm-dd",
					"error":   fmt.Sprintf("edit user information failed: %v", err),
				},
			)
			return
		}
		dob_unix = processedDOB.Unix()
	}

	_, err := authpostClient.EditUser(
		context.Background(),
		&pb_authpost.UserDetailInfo{
			UserName:     username,
			UserPassword: password,
			FirstName:    firstname,
			LastName:     lastname,
			Dob:          dob_unix,
			Email:        email,
		},
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("edit user information failed: %v", err)})
		return
	}

	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "edit user information succeeded!"})
}
