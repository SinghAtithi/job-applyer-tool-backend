package routes

import (
	"time"

	"example.com/api/handler"
	"example.com/api/middleware"
	"example.com/pkg/cache"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Default cache TTL of 25 minutes
const cacheTTL = 25 * time.Minute

// Default request timeout
const requestTimeout = 30 * time.Second

// Default body size limit (10MB)
const bodySizeLimit = 10 * 1024 * 1024

// SetupRoutes registers all application routes on the given router
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Apply global middleware
	setupGlobalMiddleware(router)

	// API routes group with versioning
	api := router.Group("/api/v1")
	{
		setupUserRoutes(api, db)
		setupResumeRoutes(api, db)
		setupCoverLetterRoutes(api, db)
	}

	// Health check route (outside API versioning)
	router.GET("/health", middleware.HealthCheck())
}

// setupGlobalMiddleware applies middleware to the router
func setupGlobalMiddleware(router *gin.Engine) {
	router.Use(middleware.RequestID())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.BodySizeLimit(bodySizeLimit))
	router.Use(middleware.ContentTypeValidator(
		"application/json",
		"application/pdf",
		"multipart/form-data",
	))
}

// setupUserRoutes registers user-related routes
func setupUserRoutes(api *gin.RouterGroup, db *gorm.DB) {
	userHandler := handler.NewUserHandler(db)
	userRoutes := api.Group("/users")
	{
		userRoutes.GET("", cache.CacheMiddleware(cacheTTL), userHandler.GetAllUsers)
		userRoutes.GET("/:id", cache.CacheMiddleware(cacheTTL), userHandler.GetUser)
		userRoutes.POST("", userHandler.CreateUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.PATCH("/:id", userHandler.PatchUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}
}

// setupResumeRoutes registers resume-related routes
func setupResumeRoutes(api *gin.RouterGroup, db *gorm.DB) {
	resumeHandler := handler.NewResumeHandler(db)
	resumeRoutes := api.Group("/resumes")
	{
		resumeRoutes.GET("", cache.CacheMiddleware(cacheTTL), resumeHandler.GetAllResumes)
		resumeRoutes.GET("/:id", cache.CacheMiddleware(cacheTTL), resumeHandler.GetResume)
		resumeRoutes.POST("", resumeHandler.CreateResume)
		resumeRoutes.PUT("/:id", resumeHandler.UpdateResume)
		resumeRoutes.PATCH("/:id", resumeHandler.PatchResume)
		resumeRoutes.DELETE("/:id", resumeHandler.DeleteResume)
		resumeRoutes.HEAD("/:id/availability", resumeHandler.IsResumeAvailable)
		resumeRoutes.GET("/:id/availability", resumeHandler.IsResumeAvailable)
	}
}

// setupCoverLetterRoutes registers cover-letter-related routes
func setupCoverLetterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	coverLetterHandler := handler.NewCoverLetterHandler(db)
	coverLetterRoutes := api.Group("/cover-letters")
	{
		coverLetterRoutes.GET("", cache.CacheMiddleware(cacheTTL), coverLetterHandler.GetAllCoverLetters)
		coverLetterRoutes.GET("/:id", cache.CacheMiddleware(cacheTTL), coverLetterHandler.GetCoverLetter)
		coverLetterRoutes.POST("", middleware.RequestTimeout(requestTimeout), coverLetterHandler.CreateCoverLetter)
		coverLetterRoutes.PUT("/:id", coverLetterHandler.UpdateCoverLetter)
		coverLetterRoutes.PATCH("/:id", coverLetterHandler.PatchCoverLetter)
		coverLetterRoutes.DELETE("/:id", coverLetterHandler.DeleteCoverLetter)
	}
}
