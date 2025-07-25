package handler

import (
	"example.com/helper"
	"example.com/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/ledongthuc/pdf"
)

type ResumeHandler struct {
	db *gorm.DB
}

func NewResumeHandler(db *gorm.DB) *ResumeHandler {
	return &ResumeHandler{db: db}
}

func (h *ResumeHandler) GetAllResumes(c *gin.Context) {
	resumeHandlerHelper := helper.NewResumeHandlerHelper(h.db)
	var resumes []models.Resume
	resumes, err := resumeHandlerHelper.GetAllResume(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch resumes",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resumes)

}

func (h *ResumeHandler) GetResume(c *gin.Context) {
	resumeHandlerHelper := helper.NewResumeHandlerHelper(h.db)
	var resume models.Resume
	resume, err := resumeHandlerHelper.GetResume(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch resumes",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resume)
}

func (h *ResumeHandler) CreateResume(c *gin.Context) {
	//There are two options, either pass the entire file as a multipart form or just the text content.
	resumeHandlerHelper := helper.NewResumeHandlerHelper(h.db)
	resumeInJson, err := resumeHandlerHelper.HandleParseResume(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	responseData, err := resumeHandlerHelper.CommitResumeToDB(resumeInJson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resumeUUID, err := uuid.NewUUID()
	if err != nil || resumeUUID == uuid.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate UUID"})
		return
	}
	resumeInJson.UserID = resumeUUID

	//helper.commitResumeToDB(resumeInJson)

	c.JSON(http.StatusCreated, gin.H{"message": "CreateResume", "resume": responseData})
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
