package infrastructure

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/klauspost/compress/s2"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
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
		Addr:         cfg.RedisURL,
		Password:     cfg.RedisPass,
		DB:           db,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,              // Increase connection pool
		MinIdleConns: 5,               // Keep minimum idle connections
		MaxIdleConns: 10,              // Maximum idle connections
		PoolTimeout:  4 * time.Second, // Pool timeout
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

	// Get compressed data from Redis
	compressed, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return false
	}

	// Decompress using S2
	data, err := s2.Decode(nil, compressed)
	if err != nil {
		return false
	}

	// Unmarshal using MessagePack
	if err := msgpack.Unmarshal(data, dest); err != nil {
		return false
	}

	return true
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}) {
	if c.client == nil {
		return
	}

	// Marshal using MessagePack
	data, err := msgpack.Marshal(value)
	if err != nil {
		log.Errorf("Failed to marshal data: %v", err)
		return
	}

	// Compress using S2
	compressed := s2.Encode(nil, data)

	// Store in Redis with TTL (24 hours)
	if err := c.client.Set(ctx, key, compressed, 24*time.Hour).Err(); err != nil {
		log.Errorf("Failed to set cache: %v", err)
	}
}

// Delete removes a key from cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Del(ctx, key).Err()
}

// DeletePattern removes all keys matching a pattern
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	if c.client == nil {
		return nil
	}

	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			log.Errorf("Failed to delete key %s: %v", iter.Val(), err)
		}
	}

	return iter.Err()
}

// Clear removes all keys from the cache
func (c *Cache) Clear(ctx context.Context) error {
	if c.client == nil {
		return nil
	}
	return c.client.FlushDB(ctx).Err()
}
