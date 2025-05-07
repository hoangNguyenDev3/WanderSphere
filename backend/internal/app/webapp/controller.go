package webapp

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
	v1 "github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/v1"
)

type WebController struct {
	webService *service.WebService
	router     *gin.Engine
	port       int
}

func (wc *WebController) Run() {
	wc.router.Run(fmt.Sprintf(":%d", wc.port))
}

func NewWebController(cfg *configs.WebConfig) (*WebController, error) {
	webService := service.NewWebService(cfg)
	if err != nil {
		return nil, err
	}

	router := gin.Default()
	for _, version := range cfg.APIVersions {
		verXRouter := router.Group(version)
		if version == "v1" {
			v1.AddAllRouter(verXRouter, webService.authpostClient, webService.newsfeedClient)
		}
	}

	webController := WebController{
		webService: *webService,
		router:     router,
		port:       cfg.Port,
	}

	return &webController, nil
}
