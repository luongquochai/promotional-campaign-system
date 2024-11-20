package main

import (
	"github.com/luongquochai/promotional-campaign-system/config"
	"github.com/luongquochai/promotional-campaign-system/models"
)

func main() {
	config.InitDB()
	config.DB.AutoMigrate(&models.User{}, &models.Campaign{}, &models.Voucher{}, &models.Purchase{}) // Add other models as needed
	println("Database migration completed!")
}
