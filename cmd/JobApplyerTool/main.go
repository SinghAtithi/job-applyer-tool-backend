package main

import (
	"example.com/api/routes"
	_ "example.com/internal/config"
	internalConfig "example.com/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	config := internalConfig.NewApplicationConfig()

	routes.SetupRoutes(router)

	router.Run(":" + config.GetPort())

}
