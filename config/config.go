package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB Config
var DB *gorm.DB

// Redis Config
var RedisClient *redis.Client

// Initialize the MySQL database connection
func InitDB() {
	var err error
	dsn := "root:Chipgau164@@tcp(localhost:3306)/promotional_campaign?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Successfully connected to MySQL Database.")
}

// Initialize Redis connection
func InitRedis() {
	// Create a context for Redis client
	ctx := context.Background()

	// Initialize Redis client
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
		DB:   0,                // Default Redis database
	})

	// Ping Redis to check connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Successfully connected to Redis.")
}
