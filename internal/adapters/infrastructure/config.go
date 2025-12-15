package infrastructure

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	CorsOrigins string
	Port        string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Warnf("Error loading .env file: %v", err)
	}

	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		CorsOrigins: os.Getenv("CORS_ORIGINS"),
		Port:        os.Getenv("PORT"),
	}
}
