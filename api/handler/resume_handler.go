package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResumeHandler struct{}

func NewResumeHandler() *ResumeHandler {
	return &ResumeHandler{}
}

func (h *ResumeHandler) GetAllResumes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetAllResumes"})
}

func (h *ResumeHandler) GetResume(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetResume"})
}

func (h *ResumeHandler) CreateResume(c *gin.Context) {
	// Parse the multipart form
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing form"})
		return
	}

	// Get the file from the form
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving file"})
		return
	}
	defer file.Close()

	// Get the other form values
	userID := c.Request.FormValue("userID")
	title := c.Request.FormValue("title")
	description := c.Request.FormValue("description")

	// You can now process the file and the other form values
	// For example, save the file to a specific location
	// and create a new resume record in the database.

	c.JSON(http.StatusCreated, gin.H{
		"message":     "CreateResume",
		"filename":    handler.Filename,
		"userID":      userID,
		"title":       title,
		"description": description,
	})
}

func (h *ResumeHandler) UpdateResume(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "UpdateResume"})
}

func (h *ResumeHandler) PatchResume(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PatchResume"})
}

func (h *ResumeHandler) DeleteResume(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DeleteResume"})
}
