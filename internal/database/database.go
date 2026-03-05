package database

import (
	"fmt"
	"net/url"
	"os"

	"example.com/internal/models"
	"example.com/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// getDBConnectionDetails retrieves database configuration from environment variables
func getDBConnectionDetails() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Port:     getEnv("DB_PORT", "5432"),
		Username: getEnv("DB_USERNAME", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "resumedb"),
		SSLMode:  getEnv("DB_SSL_MODE", "require"),
	}
}

// getEnv retrieves environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewDatabase creates a new database connection
func NewDatabase() (*gorm.DB, error) {

	// 1. Prefer full connection string (cloud platforms)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {

		db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect using DATABASE_URL: %w", err)
		}

		logger.Info("Database connection established using DATABASE_URL")
		return db, nil
	}

	// 2. Fallback to manual config (local dev)
	dbConfig := getDBConnectionDetails()

	if dbConfig.SSLMode == "disable" {
		logger.Warn("DB SSL mode is 'disable' — TLS is disabled for DB connections. Ensure this is intentional.")
	}

	userinfo := url.UserPassword(dbConfig.Username, dbConfig.Password)

	u := &url.URL{
		Scheme: "postgres",
		User:   userinfo,
		Host:   fmt.Sprintf("%s:%s", dbConfig.Host, dbConfig.Port),
		Path:   "/" + dbConfig.DBName,
	}

	q := u.Query()
	q.Set("sslmode", dbConfig.SSLMode)
	u.RawQuery = q.Encode()

	dsn := u.String()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info("Database connection established using DB_* environment variables")
	return db, nil
}

// EnsureTablesExist creates tables if they don't exist using AutoMigrate
func EnsureTablesExist(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Resume{},
		&models.CoverLetterTable{},
		&models.JobDescriptionTable{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate tables: %w", err)
	}

	logger.Info("All tables ensured to exist")
	return nil
}
