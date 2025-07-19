package handler

import (
	"example.com/api/handler/helper"
	"log"
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
	resume, err := helper.HandleCreateResume(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Error while making POST operation in Resume" + err.Error())
		return
	}
	c.JSON(http.StatusCreated, resume)
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
