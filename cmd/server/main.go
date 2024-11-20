package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	// Create an HTTP server
	server := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	// Run the server in a goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	// Create a channel to listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) // Listen for SIGINT (Ctrl+C) or other termination signals

	<-quit // Block until a signal is received
	log.Println("Shutting down server...")

	// Create a context with a timeout for cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Cleanup resources (e.g., DB, Redis)
	database.CloseDB() // Add this method to close database connections
	redis.CloseRedis() // Add this method to close Redis connections

	log.Println("Server exiting")
}
