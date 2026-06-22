package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/starquill/URL-Shortening-Service/internal/cache"
	"github.com/starquill/URL-Shortening-Service/internal/config"
	"github.com/starquill/URL-Shortening-Service/internal/database"
	"github.com/starquill/URL-Shortening-Service/internal/generator"
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

	// Connect to Redis
	ttl, err := time.ParseDuration(cfg.Redis.TTL)
	if err != nil {
		log.Fatalf("invalid redis ttl: %v", err)
	}

	redisCache, err := cache.NewRedisCache(cfg.Redis.URL, cfg.Redis.Password, 0, ttl)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer redisCache.Close()
	log.Println("redis cache initialized")

	// Initialize short code generator
	codeGenerator := generator.NewShortCodeGenerator(urlStore, 7)
	log.Printf("short code generator initialized: %v", codeGenerator != nil)

	// Initialize URL handler
	urlHandler := handler.NewURLHandler(urlStore, redisCache, codeGenerator, cfg.Server.BaseURL)

	// Setup router
	r := chi.NewRouter()

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   cfg.Cors.AllowedOrigins,
		AllowedMethods:   cfg.Cors.AllowedMethods,
		AllowedHeaders:   cfg.Cors.AllowedHeaders,
		AllowCredentials: true,
	})
	r.Use(corsHandler.Handler)

	// Serve static files from frontend build
	fileServer := http.FileServer(http.Dir("./frontend/build"))
	r.Handle("/static/*", fileServer)
	r.Handle("/favicon.ico", fileServer)
	r.Handle("/manifest.json", fileServer)

	// Health check
	r.Get("/health", handler.Health)

	// Redirect endpoint (must come before /api routes)
	r.Get("/{code}", urlHandler.RedirectURL)

	// API endpoints
	r.Post("/shorten", urlHandler.CreateShortURL)
	r.Route("/api/urls/{code}", func(r chi.Router) {
		r.Get("/", urlHandler.GetURL)
		r.Put("/", urlHandler.UpdateURL)
		r.Delete("/", urlHandler.DeleteURL)
	})

	address := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("server listening on %s", address)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
