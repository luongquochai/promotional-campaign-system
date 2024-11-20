package database

import (
	"log"

	"github.com/luongquochai/promotional-campaign-system/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB Config
var DB *gorm.DB

// Initialize the MySQL database connection
func InitDB(config *config.Config) {
	var err error
	dsn := config.Dsn
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Successfully connected to MySQL Database.")
}
