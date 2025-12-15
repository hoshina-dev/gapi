package infrastructure

import "os"

type Config struct {
	DatabaseURL string
	CorsOrigins string
	Port        string
}

func LoadConfig() Config {
	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		CorsOrigins: os.Getenv("CORS_ORIGINS"),
		Port:        os.Getenv("PORT"),
	}
}
