package webapp

import (
	"fmt"

	"github.com/hoangNguyenDev3/WanderSphere/configs"
	v1 "github.com/hoangNguyenDev3/WanderSphere/internal/app/webapp/v1"

	client_authpost "github.com/hoangNguyenDev3/WanderSphere/pkg/client/authpost"

	"github.com/gin-gonic/gin"
	pb_authpost "github.com/hoangNguyenDev3/WanderSphere/pkg/types/proto/pb/authpost"
	pb_newsfeed "github.com/hoangNguyenDev3/WanderSphere/pkg/types/proto/pb/newsfeed"
)

type WebService struct {
	authpostClient pb_authpost.AuthenticationAndPostClient
	newsfeedClient pb_newsfeed.NewsfeedClient
}

type WebController struct {
	webService *WebService
	router     *gin.Engine
	port       int
}

func (wc *WebController) Run() {
	wc.router.Run(fmt.Sprintf(":%d", wc.port))
}

func NewWebController(cfg *configs.WebConfig) (*WebController, error) {
	authpostCient, err := client_authpost.NewClient(cfg.AuthenticateAndPost.Hosts)
	if err != nil {
		return nil, err
	}

	webService := WebService{
		authpostClient: authpostCient,
	}

	router := gin.Default()
	for _, version := range cfg.APIVersions {
		verXRouter := router.Group(version)
		if version == "v1" {
			v1.AddAllRouter(verXRouter, webService.authpostClient, webService.newsfeedClient)
		}
	}

	webController := WebController{
		webService: webService,
		router:     router,
		port:       cfg.Port,
	}

	return &webController, nil
}
