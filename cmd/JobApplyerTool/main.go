package main

import (
	"example.com/api/routes"
	internalConfig "example.com/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // Add the actual dotenv package
	"log"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	config := internalConfig.NewApplicationConfig()

	routes.SetupRoutes(router)

	router.Run(":" + config.GetPort())
}
