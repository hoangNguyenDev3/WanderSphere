package main

import (
	"flag"
	"log"

	"github.com/hoangNguyenDev3/WanderSphere/configs"
	"github.com/hoangNguyenDev3/WanderSphere/internal/app/webapp"
)

func main() {
	cfgPath := flag.String("conf", "configs/config.yaml", "path to the config file")

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
