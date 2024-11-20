package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/config"
	database "github.com/luongquochai/promotional-campaign-system/database/mysql"
	"github.com/luongquochai/promotional-campaign-system/database/redis"
	"github.com/luongquochai/promotional-campaign-system/routes"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config/file/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize DB and Redis
	database.InitDB(cfg)
	redis.InitRedis(cfg)

	router := gin.Default()

	// Set up routes
	routes.SetupRoutes(router)

	// Run server
	router.Run(":8080")
}
