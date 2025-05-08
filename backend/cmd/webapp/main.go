package main

import (
	"flag"
	"log"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/app/webapp"
)

func main() {
	cfgPath := flag.String("conf", "config.yaml", "path to the config file")

	cfg, err := configs.ParseConfig(*cfgPath)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	web_controller, err := webapp.NewWebController(&cfg.WebConfig)
	if err != nil {
		log.Fatalf("Failed to create web controller: %v", err)
		return
	}

	web_controller.Run()
}
