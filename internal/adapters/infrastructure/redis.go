package infrastructure

import (
	"context"
	"encoding/json"
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

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Get(ctx context.Context, key string, dest interface{}) bool {
	if c.client == nil {
		return false
	}
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	if json.Unmarshal([]byte(data), dest) != nil {
		return false
	}
	return true
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}) {
	if c.client == nil {
		return
	}
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	c.client.Set(ctx, key, data, 0)
}
