package routes

import (
	"time"

	"example.com/api/handler"
	"example.com/pkg/cache"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Default cache TTL of 25 minutes
const cacheTTL = 25 * time.Minute

// SetupRoutes registers all application routes on the given router
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// User routes
	userHandler := handler.NewUserHandler(db)
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", cache.CacheMiddleware(cacheTTL), userHandler.GetAllUsers)
		userRoutes.GET("/:id", cache.CacheMiddleware(cacheTTL), userHandler.GetUser)
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.PATCH("/:id", userHandler.PatchUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	// Resume routes
	resumeHandler := handler.NewResumeHandler(db)
	resumeRoutes := router.Group("/resumes")
	{
		resumeRoutes.GET("/", cache.CacheMiddleware(cacheTTL), resumeHandler.GetAllResumes)
		resumeRoutes.GET("/:id", cache.CacheMiddleware(cacheTTL), resumeHandler.GetResume)
		resumeRoutes.POST("/", resumeHandler.CreateResume)
		resumeRoutes.PUT("/:id", resumeHandler.UpdateResume)
		resumeRoutes.PATCH("/:id", resumeHandler.PatchResume)
		resumeRoutes.DELETE("/:id", resumeHandler.DeleteResume)
		resumeRoutes.GET("/isResumeAvailable/:id", resumeHandler.IsResumeAvailable)
	}

	// Cover letter routes
	coverLetterHandler := handler.NewCoverLetterHandler(db)
	coverLetterRoutes := router.Group("/coverLetter")
	{
		coverLetterRoutes.GET("/", cache.CacheMiddleware(cacheTTL), coverLetterHandler.GetAllCoverLetters)
		coverLetterRoutes.GET("/:id", cache.CacheMiddleware(cacheTTL), coverLetterHandler.GetCoverLetter)
		coverLetterRoutes.POST("/", coverLetterHandler.CreateCoverLetter)
		coverLetterRoutes.PUT("/:id", coverLetterHandler.UpdateCoverLetter)
		coverLetterRoutes.PATCH("/:id", coverLetterHandler.PatchCoverLetter)
		coverLetterRoutes.DELETE("/:id", coverLetterHandler.DeleteCoverLetter)
	}
}
