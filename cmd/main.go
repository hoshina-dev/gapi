package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hoshina-dev/gapi/internal/adapters/graph"
	"github.com/hoshina-dev/gapi/internal/adapters/http"
	"github.com/hoshina-dev/gapi/internal/adapters/infrastructure"
	"github.com/hoshina-dev/gapi/internal/adapters/repository"
	"github.com/hoshina-dev/gapi/internal/core/services"
)

func main() {
	cfg := infrastructure.LoadConfig()

	db := infrastructure.ConnectDB(cfg.DatabaseURL)
	redisClient := infrastructure.ConnectRedis(cfg)

	countryRepo := repository.NewAdminAreaRepository(db, redisClient)
	countryService := services.NewAdminAreaService(countryRepo)
	resolver := graph.NewResolver(countryService)

	app := http.SetupRouter(resolver, cfg)

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
