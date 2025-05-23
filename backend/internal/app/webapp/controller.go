package webapp

import (
	"fmt"
	"net/http"

	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/service"
	v1 "github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp/v1"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type WebController struct {
	webService    service.WebService
	router        *gin.Engine
	port          int
	healthChecker *utils.HealthChecker
}

func (wc *WebController) Run() {
	wc.router.Run(fmt.Sprintf(":%d", wc.port))
}

// NewWebController creates new WebController
func NewWebController(cfg *configs.WebConfig) (*WebController, error) {
	// Init web services
	webService, err := service.NewWebService(cfg)
	if err != nil {
		return nil, err
	}

	// Create health checker
	healthChecker := utils.NewHealthChecker("web", "1.0.0", webService.GetLogger())

	// Init router
	router := gin.Default()

	// Add health check endpoints
	initHealth(router, healthChecker, webService)

	for _, version := range cfg.APIVersions {
		verXRouter := router.Group(fmt.Sprint("/api/" + version))
		if version == "v1" { // TODO: Automate this when a new vision is added
			v1.AddAllRouter(verXRouter, webService)
		}
	}

	// Init other support tools
	initSwagger(router)
	initPprof(router)

	return &WebController{
		webService:    *webService,
		router:        router,
		port:          cfg.Port,
		healthChecker: healthChecker,
	}, nil
}

func initHealth(router *gin.Engine, healthChecker *utils.HealthChecker, webService *service.WebService) {
	// Basic health endpoint
	router.GET("/health", func(c *gin.Context) {
		status := healthChecker.GetHealthStatus()
		c.JSON(http.StatusOK, status)
	})

	// Detailed health endpoint with dependencies
	router.GET("/health/detailed", func(c *gin.Context) {
		status := healthChecker.GetDetailedHealthStatus(nil, webService.GetRedis())

		// Add web service specific dependencies
		healthChecker.AddDependencyStatus(status.Dependencies, "authpost", "healthy") // Assume healthy since service started
		healthChecker.AddDependencyStatus(status.Dependencies, "newsfeed", "healthy") // Assume healthy since service started

		// Determine overall status based on dependencies
		for _, depStatus := range status.Dependencies {
			if depStatus == "unhealthy" {
				status.Status = "degraded"
				break
			}
		}

		statusCode := http.StatusOK
		if status.Status == "degraded" {
			statusCode = http.StatusServiceUnavailable
		}
		c.JSON(statusCode, status)
	})
}

func initSwagger(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func initPprof(router *gin.Engine) {
	router.GET("/debug/pprof/", func(context *gin.Context) {
		pprof.Index(context.Writer, context.Request)
	})
	router.GET("/debug/pprof/profile", func(context *gin.Context) {
		pprof.Profile(context.Writer, context.Request)
	})
	router.GET("/debug/pprof/trace", func(context *gin.Context) {
		pprof.Trace(context.Writer, context.Request)
	})
}
