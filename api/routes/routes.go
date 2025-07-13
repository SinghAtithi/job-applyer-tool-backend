package routes

import (
	"example.com/api/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	userHandler := handler.NewUserHandler()
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

	resumeHandler := handler.NewResumeHandler()
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
}
