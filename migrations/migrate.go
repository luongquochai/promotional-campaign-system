package main

import (
	"log"

	"github.com/luongquochai/promotional-campaign-system/config"
	database "github.com/luongquochai/promotional-campaign-system/database/mysql"
	"github.com/luongquochai/promotional-campaign-system/models"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config/file/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	database.InitDB(cfg)
	database.DB.AutoMigrate(&models.User{}, &models.Campaign{}, &models.Voucher{}, &models.Purchase{}) // Add other models as needed
	println("Database migration completed!")
}
