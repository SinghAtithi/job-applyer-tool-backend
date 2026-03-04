package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/api/middleware"
	"example.com/api/routes"
	internalConfig "example.com/internal/config"
	"example.com/internal/database"
	"example.com/pkg/cache"
	"example.com/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Parse command-line flags
	devMode := flag.Bool("dev", false, "Run in development mode")
	debugMode := flag.Bool("debug", false, "Enable debug logging with stack traces")
	flag.Parse()

	// Load environment variables from .env file
	envFile := ".env"
	if *devMode {
		if _, err := os.Stat(".env.development"); err == nil {
			envFile = ".env.development"
		}
	}

	if err := godotenv.Load(envFile); err != nil {
		// Use logger even before full init — Default() is always available
		logger.Warn("failed to load %s: %v", envFile, err)
	}

	// If --dev flag is passed, set environment variables
	if *devMode {
		os.Setenv("CONFIG_MODE", "development")
		os.Setenv("DEV", "true")
	}

	// Load application configuration
	config, err := internalConfig.NewApplicationConfig()
	if err != nil {
		logger.Fatal("failed to load application config: %v", err)
	}

	if *devMode {
		config.SetDevMode(true)
	}
	if *debugMode {
		config.SetDebugMode(true)
	}

	// Initialize logger with config-based level and debug flag
	logger.Initialize(config.GetLogLevel(), config.IsDebug())

	if config.IsDev() {
		logger.Info("running in development mode")
	}
	if config.IsDebug() {
		logger.Info("debug mode enabled — stack traces will be included on errors")
	}

	// Initialize database
	db, err := database.NewDatabase()
	if err != nil {
		logger.Fatal("failed to connect to database: %v", err)
	}

	if err := database.EnsureTablesExist(db); err != nil {
		logger.Fatal("failed to create tables: %v", err)
	}

	// Initialize Redis cache (non-fatal if unavailable)
	cacheConfig := cache.DefaultCacheConfig()
	if err := cache.InitializeRedis(cacheConfig); err != nil {
		logger.Warn("failed to connect to Redis: %v (running without cache)", err)
	}

	// Create router
	router := gin.New()

	// Apply global middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", middleware.HealthCheck())

	// Application routes
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
		logger.Info("starting server on port %s", config.GetPort())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown: %v", err)
	}

	// Close DB connection
	if sqlDB, err := db.DB(); err == nil {
		if cerr := sqlDB.Close(); cerr != nil {
			logger.Error("failed to close database connection: %v", cerr)
		} else {
			logger.Info("database connection closed")
		}
	} else {
		logger.Error("unable to obtain underlying sql.DB to close: %v", err)
	}

	logger.Info("server exited properly")
}
