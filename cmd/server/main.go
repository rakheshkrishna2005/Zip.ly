package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rakheshkrishna2005/url-shortener/internal/api"
	"github.com/rakheshkrishna2005/url-shortener/internal/api/handlers"
	"github.com/rakheshkrishna2005/url-shortener/internal/config"
	"github.com/rakheshkrishna2005/url-shortener/internal/repository/postgres"
	"github.com/rakheshkrishna2005/url-shortener/internal/service"
)

func main() {
	// Load configuration
	cfg := config.New()
	
	// Connect to database
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	// Set up connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// Create repositories
	urlRepo := postgres.NewURLRepository(db)
	
	// Create services
	urlService := service.NewURLService(urlRepo, cfg)
	
	// Create handlers
	urlHandler := handlers.NewURLHandler(urlService)
	healthHandler := handlers.NewHealthHandler()
	
	// Set up router
	router := api.NewRouter(urlHandler, healthHandler)
	
	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
	
	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Server shutting down...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server exited properly")
}