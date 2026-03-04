package handler

import (
	"net/http"

	"example.com/helper"
	"example.com/internal/dto"
	"example.com/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ResumeHandler handles resume-related HTTP requests
type ResumeHandler struct {
	db *gorm.DB
}

// NewResumeHandler creates a new resume handler
func NewResumeHandler(db *gorm.DB) *ResumeHandler {
	return &ResumeHandler{db: db}
}

// GetAllResumes retrieves all resumes
func (h *ResumeHandler) GetAllResumes(c *gin.Context) {
	resumeHelper := helper.NewResumeHandlerHelper(h.db)

	resumes, err := resumeHelper.GetAllResume(c)
	if err != nil {
		logger.Error("failed to fetch resumes: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR", "Failed to fetch resumes", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(resumes, "Resumes retrieved successfully"))
}

// GetResume retrieves a single resume
func (h *ResumeHandler) GetResume(c *gin.Context) {
	resumeHelper := helper.NewResumeHandlerHelper(h.db)

	resume, err := resumeHelper.GetResume(c)
	if err != nil {
		logger.Error("failed to fetch resume: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR", "Failed to fetch resume", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(resume, "Resume retrieved successfully"))
}

// CreateResume processes an uploaded resume and stores it
func (h *ResumeHandler) CreateResume(c *gin.Context) {
	resumeHelper := helper.NewResumeHandlerHelper(h.db)

	resumeInJson, err := resumeHelper.HandleParseResume(c)
	if err != nil {
		logger.Error("failed to parse resume: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"BAD_REQUEST", "Failed to parse resume", err.Error(),
		))
		return
	}

	// Read username from request body (form field)
	if username := c.PostForm("username"); username != "" {
		resumeInJson.UserName = username
	}

	responseData, err := resumeHelper.CommitResumeToDB(resumeInJson)
	if err != nil {
		logger.Error("failed to commit resume to database: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR", "Failed to save resume", err.Error(),
		))
		return
	}

	resumeUUID, err := uuid.NewUUID()
	if err != nil || resumeUUID == uuid.Nil {
		logger.Error("failed to generate UUID for resume: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"INTERNAL_ERROR", "Failed to generate resume identifier", "",
		))
		return
	}
	resumeInJson.UserID = resumeUUID

	c.JSON(http.StatusCreated, dto.SuccessResponse(responseData, "Resume created successfully"))
}

// UpdateResume updates an existing resume
// TODO: Implement full resume update logic
func (h *ResumeHandler) UpdateResume(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.ErrorResponse(
		"NOT_IMPLEMENTED", "UpdateResume is not yet implemented", "",
	))
}

// PatchResume partially updates a resume
// TODO: Implement partial resume update logic
func (h *ResumeHandler) PatchResume(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.ErrorResponse(
		"NOT_IMPLEMENTED", "PatchResume is not yet implemented", "",
	))
}

// DeleteResume deletes a resume
// TODO: Implement resume deletion logic
func (h *ResumeHandler) DeleteResume(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.ErrorResponse(
		"NOT_IMPLEMENTED", "DeleteResume is not yet implemented", "",
	))
}

// IsResumeAvailable checks if a resume exists for the given identifier
func (h *ResumeHandler) IsResumeAvailable(c *gin.Context) {
	resumeHelper := helper.NewResumeHandlerHelper(h.db)
	id := c.Param("id")

	exists, err := resumeHelper.IsResumeAvailable(id)
	if err != nil {
		logger.Error("failed to check resume availability (id=%s): %v", id, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR", "Failed to check resume availability", err.Error(),
		))
		return
	}

	if exists {
		c.JSON(http.StatusOK, dto.SuccessResponse(nil, "Resume is available"))
	} else {
		c.JSON(http.StatusNotFound, dto.ErrorResponse(
			"NOT_FOUND", "Resume not found", "",
		))
	}
}
