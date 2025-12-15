package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hoshina-dev/gapi/internal/adapters/graph"
	"github.com/hoshina-dev/gapi/internal/adapters/infrastructure"
	"github.com/hoshina-dev/gapi/internal/adapters/repository"
	"github.com/hoshina-dev/gapi/internal/core/services"
)

func main() {
	cfg := infrastructure.LoadConfig()

	db := infrastructure.ConnectDB(cfg.DatabaseURL)

	countryRepo := repository.NewCountryRepository(db)
	countryService := services.NewCountryService(countryRepo)
	resolver := graph.NewResolver(countryService)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CorsOrigins,
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Get("/health", healthCheck)
	app.Get("/", graph.PlaygroundHandler())
	app.All("/query", graph.GraphQLHandler(resolver))

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server running on :%s", cfg.Port)
	log.Printf("Connect to http://localhost:%s/ for GraphQL playground", cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok", "time": time.Now()})
}
