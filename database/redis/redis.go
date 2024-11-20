package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/luongquochai/promotional-campaign-system/config"
)

var (
	Redis *redis.Client
	ctx   = context.Background()
)

// Initialize Redis Client
func InitRedis(config *config.Config) {
	Redis = redis.NewClient(&redis.Options{
		Addr: config.Redis_Addr, // Redis server address
		// Password: os.Getenv("REDIS_PASSWORD"), // No password set
		DB: config.Redis_DB, // Use default DB
	})

	// Ping Redis to test the connection
	_, err := Redis.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connected successfully")
}

func CloseRedis() {
	if Redis != nil {
		if err := Redis.Close(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		} else {
			log.Println("Redis connection closed.")
		}
	}
}

// Set key-value pair in Redis
func Set(key string, value string, expiration time.Duration) error {
	return Redis.Set(ctx, key, value, expiration).Err()
}

// Get value by key from Redis
func Get(key string) (string, error) {
	return Redis.Get(ctx, key).Result()
}

// Delete key from Redis
func Delete(key string) error {
	return Redis.Del(ctx, key).Err()
}
