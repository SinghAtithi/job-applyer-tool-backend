package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/api/middleware"
	"example.com/api/routes"
	"example.com/pkg/cache"
	internalConfig "example.com/internal/config"
	"example.com/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize database
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Ensure tables exist
	if err := database.EnsureTablesExist(db); err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	// Initialize Redis cache
	cacheConfig := cache.DefaultCacheConfig()
	if err := cache.InitializeRedis(cacheConfig); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Running without cache.", err)
	}

	// Create router with custom configuration
	router := gin.New()

	// Apply global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", middleware.HealthCheck())

	// Setup application routes
	config, cerr := internalConfig.NewApplicationConfig()
	if cerr != nil {
		log.Fatalf("Error loading application config: %v", cerr)
	}
	routes.SetupRoutes(router, db)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + config.GetPort(),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s", config.GetPort())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close DB connection (gorm.DB -> sql.DB)
	if sqlDB, derr := db.DB(); derr == nil {
		if cerr := sqlDB.Close(); cerr != nil {
			log.Printf("Error closing database connection: %v", cerr)
		} else {
			log.Println("Database connection closed")
		}
	} else {
		log.Printf("Unable to obtain underlying sql.DB to close: %v", derr)
	}

	fmt.Println("Server exited properly")
}
