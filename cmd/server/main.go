package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/config"
	database "github.com/luongquochai/promotional-campaign-system/database/mysql"
	"github.com/luongquochai/promotional-campaign-system/database/redis"
	_ "github.com/luongquochai/promotional-campaign-system/docs" // This imports the generated Swagger docs
	"github.com/luongquochai/promotional-campaign-system/routes"
	ginSwagger "github.com/swaggo/gin-swagger"                // Swagger UI handler
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles" // Correct import path for swaggerFiles
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

	// Initialize Gin router
	router := gin.Default()

	// Set up routes
	routes.SetupRoutes(router)

	// Serve Swagger UI at /docs
	// router.GET("/docs/*any", swagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run the server
	router.Run(":8080")
}
