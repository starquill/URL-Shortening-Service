package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/starquill/URL-Shortening-Service/internal/config"
	"github.com/starquill/URL-Shortening-Service/internal/database"
	"github.com/starquill/URL-Shortening-Service/internal/handler"
	"github.com/starquill/URL-Shortening-Service/internal/store"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Run database migrations
	log.Println("running database migrations...")
	if err := database.RunMigrations(cfg.Database.URL); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Connect to database
	pool, err := database.NewPool(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize store
	urlStore := store.NewURLStore(pool)
	log.Printf("url store initialized: %v", urlStore != nil)

	r := chi.NewRouter()
	r.Get("/health", handler.Health)

	address := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("server listening on %s", address)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
