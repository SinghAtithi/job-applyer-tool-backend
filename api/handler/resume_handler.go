package handler

import (
	"example.com/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/ledongthuc/pdf"
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
	//There are two options, either pass the entire file as a multipart form or just the text content.
	resumeInJson, err := helper.HandleParseResume(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "CreateResume", "resume": resumeInJson})
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
