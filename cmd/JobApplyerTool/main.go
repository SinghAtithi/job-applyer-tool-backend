package main

import (
	"example.com/api/routes"
	internalConfig "example.com/internal/config"
	"example.com/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // Add the actual dotenv package
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.NewDatabase()

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	router := gin.Default()

	config := internalConfig.NewApplicationConfig()

	routes.SetupRoutes(router, db)

	router.Run(":" + config.GetPort())
}
