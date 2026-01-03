package infrastructure

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg Config) *redis.Client {
	if cfg.RedisURL == "" {
		log.Info("Redis is disabled due to incomplete configuration.")
		return nil
	}

	db, err := strconv.Atoi(cfg.RedisDB)
	if err != nil {
		log.Errorf("Invalid Redis DB: %v, disabling Redis", err)
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:        cfg.RedisURL,
		Password:    cfg.RedisPass,
		DB:          db,
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
