package routes

import (
	"example.com/api/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {

	userHandler := handler.NewUserHandler(db)
	// User routes
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", userHandler.GetAllUsers)
		userRoutes.GET("/:id", userHandler.GetUser)
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.PATCH("/:id", userHandler.PatchUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	resumeHandler := handler.NewResumeHandler(db)
	// Resume routes
	resumeRoutes := router.Group("/resumes")
	{
		resumeRoutes.GET("/", resumeHandler.GetAllResumes)
		resumeRoutes.GET("/:id", resumeHandler.GetResume)
		resumeRoutes.POST("/", resumeHandler.CreateResume)
		resumeRoutes.PUT("/:id", resumeHandler.UpdateResume)
		resumeRoutes.PATCH("/:id", resumeHandler.PatchResume)
		resumeRoutes.DELETE("/:id", resumeHandler.DeleteResume)
	}

	coverLetterHandler := handler.NewCoverLetter(db)

	coverLetterRoutes := router.Group("/coverLetter")
	{
		coverLetterRoutes.GET("/", coverLetterHandler.GetAllCoverLetters)
		coverLetterRoutes.GET("/:id", coverLetterHandler.GetCoverLetter)
		coverLetterRoutes.POST("/", coverLetterHandler.CreateCoverLetter)
		coverLetterRoutes.PUT("/:id", coverLetterHandler.UpdateCoverLetter)
		coverLetterRoutes.PATCH("/:id", coverLetterHandler.PatchCoverLetter)
		coverLetterRoutes.DELETE("/:id", coverLetterHandler.DeleteCoverLetter)
	}

}
