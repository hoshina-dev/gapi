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
	RedisDB     string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Warnf("Error loading .env file: %v", err)
	}

	redisURL := os.Getenv("REDIS_URL")
	redisPass := os.Getenv("REDIS_PASSWORD")
	redisDBStr := os.Getenv("REDIS_DB")

	// Validate Redis config: all parameters must be set
	if redisURL == "" || redisDBStr == "" {
		log.Warn("Redis configuration incomplete. All REDIS_URL, REDIS_PASSWORD, and REDIS_DB must be set. Redis will not be used.")
		redisURL = ""
		redisPass = ""
		redisDBStr = ""
	} else {
		// Check if REDIS_DB is a valid integer
		_, err := strconv.Atoi(redisDBStr)
		if err != nil {
			log.Warnf("Invalid REDIS_DB=%q, Redis will not be used: %v", redisDBStr, err)
			redisURL = ""
			redisPass = ""
			redisDBStr = ""
		}
	}

	return Config{
		DatabaseURL: os.Getenv("DATA_SOURCE_NAME"),
		CorsOrigins: os.Getenv("CORS_ORIGINS"),
		Port:        os.Getenv("PORT"),
		RedisURL:    redisURL,
		RedisPass:   redisPass,
		RedisDB:     redisDBStr,
	}
}
