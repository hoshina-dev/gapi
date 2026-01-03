package infrastructure

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:        cfg.RedisURL,
		Password:    cfg.RedisPass,
		DB:          cfg.RedisDB,
		MaxRetries:  3,
		DialTimeout: 5 * time.Second,
	})

	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Errorf("Failed to connect to Redis: %v", err)
		// Return nil if connection fails, allowing the app to run without Redis
		return nil
	}

	return client
}
