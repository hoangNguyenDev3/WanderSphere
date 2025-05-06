package v1

import (
	"github.com/gin-gonic/gin"
	pb_authpost "github.com/hoangNguyenDev3/WanderSphere/pkg/types/proto/pb/authpost"
	pb_newsfeed "github.com/hoangNguyenDev3/WanderSphere/pkg/types/proto/pb/newsfeed"
)

var (
	authpostClient pb_authpost.AuthenticationAndPostClient
)

func AddAllRouter(r *gin.RouterGroup, in_authpost_client pb_authpost.AuthenticationAndPostClient, in_newsfeed_client pb_newsfeed.NewsfeedClient) {
	authpostClient = in_authpost_client

	AddUserRouter(r)
}
