// @title WanderSphere API
// @version 1.0
// @description API Server for WanderSphere Application
// @contact.name API Support
// @contact.email support@wandersphere.com
// @host localhost:19003
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in cookie
// @name session_id

package main

import (
	"flag"
	"log"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	_ "github.com/hoangNguyenDev3/WanderSphere/backend/docs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp"
)

func main() {
	// Flags
	cfgPath := flag.String("conf", "configs/files/local_webapp.yml", "Path to config file for this service")
	flag.Parse()

	// Load configs
	cfg, err := configs.GetWebConfigDirect(*cfgPath)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
		return
	}

	// Create new web controller
	web_controller, err := webapp.NewWebController(cfg)
	if err != nil {
		log.Fatalf("failed to create web controller: %v", err)
		return
	}

	log.Printf("Starting Web service on port %d", cfg.Port)
	// Serve
	web_controller.Run()
}
