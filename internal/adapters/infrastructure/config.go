package infrastructure

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	CorsOrigins string
	Port        string
	RedisURL    string
	RedisPass   string
	RedisDB     int
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Warnf("Error loading .env file: %v", err)
	}

	// Default to Redis DB 0 if REDIS_DB is not set or invalid
	redisDBStr := os.Getenv("REDIS_DB")
	RedisDBInt := 0
	if redisDBStr != "" {
		parsed, err := strconv.Atoi(redisDBStr)
		if err != nil {
			log.Warnf("Error parsing REDIS_DB=%q, defaulting to Redis DB 0: %v", redisDBStr, err)
		} else {
			RedisDBInt = parsed
		}
	}

	return Config{
		DatabaseURL: os.Getenv("DATA_SOURCE_NAME"),
		CorsOrigins: os.Getenv("CORS_ORIGINS"),
		Port:        os.Getenv("PORT"),
		RedisURL:    os.Getenv("REDIS_URL"),
		RedisPass:   os.Getenv("REDIS_PASSWORD"),
		RedisDB:     RedisDBInt,
	}
}
