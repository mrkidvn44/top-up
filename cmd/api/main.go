package main

import (
	"log"
	"top-up-api/config"
	"top-up-api/internal/app"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Enter the token with the `Bearer` prefix, e.g., `Bearer <token>`
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	app.Run(cfg)
}
