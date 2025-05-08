package main

import (
	"flag"
	"log"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp"
)

func main() {
	// Flags
	cfgPath := flag.String("conf", "config.yml", "Path to config file for this service")

	// Load configurations
	cfg, err := configs.GetWebConfig(*cfgPath)
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

	// Serve
	web_controller.Run()
}
