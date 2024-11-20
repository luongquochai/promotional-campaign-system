package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/config"
	"github.com/luongquochai/promotional-campaign-system/routes"
)

func main() {
	// Initialize DB and Redis
	config.InitDB()
	config.InitRedis()

	router := gin.Default()

	// Set up routes
	routes.SetupRoutes(router)

	// Run server
	router.Run(":8080")
}
