package handler

import (
	"net/http"

	"example.com/agent"
	"example.com/helper"
	"example.com/internal/dto"
	"example.com/internal/models"
	"example.com/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CoverLetterHandler handles cover-letter–related HTTP requests
type CoverLetterHandler struct {
	db *gorm.DB
}

// NewCoverLetterHandler creates a new cover letter handler
func NewCoverLetterHandler(db *gorm.DB) *CoverLetterHandler {
	return &CoverLetterHandler{db: db}
}

// CreateCoverLetter generates a cover letter from a job description URL
func (h *CoverLetterHandler) CreateCoverLetter(c *gin.Context) {
	var requestBody models.JobDescriptionTable

	if err := c.BindJSON(&requestBody); err != nil || requestBody.URL == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"VALIDATION_ERROR", "Invalid request body: URL is required", "",
		))
		return
	}

	coverLetterHelper := helper.NewCoverLetterHelper(h.db)

	jobDescriptionData, isPresent := coverLetterHelper.GetJobDescriptionFromDB(requestBody.URL)

	var jobDescriptionContent string
	var err error
	if isPresent {
		jobDescriptionContent = jobDescriptionData.ContentInfo
	} else {
		jobDescriptionContent, err = agent.GetJobDescription(requestBody.TextContent)
		if err != nil {
			logger.Error("failed to get job description from agent: %v", err)
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(
				"BAD_REQUEST", "Failed to extract job description", err.Error(),
			))
			return
		}
		if commitErr := coverLetterHelper.CommitJobDescription(jobDescriptionContent, requestBody.URL); commitErr != nil {
			logger.Error("failed to save job description: %v", commitErr)
			// Non-fatal: continue even if caching fails
		}
	}

	coverLetter, err := coverLetterHelper.GetCoverLetter(jobDescriptionContent, requestBody.UserName, requestBody.URL)
	if err != nil || coverLetter.ContentInfo == "" {
		logger.Error("failed to generate cover letter: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"INTERNAL_ERROR", "Unable to create cover letter", "",
		))
		return
	}

	pdfBytes, err := coverLetterHelper.ReturnPDFResponse(coverLetter.ContentInfo)
	if err != nil {
		logger.Error("failed to generate PDF: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"INTERNAL_ERROR", "Failed to generate PDF", err.Error(),
		))
		return
	}

	coverLetterHelper.SetHeadersForPDF(c, pdfBytes)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GetAllCoverLetters retrieves all cover letters
// TODO: Implement retrieval logic
func (h *CoverLetterHandler) GetAllCoverLetters(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.ErrorResponse(
		"NOT_IMPLEMENTED", "GetAllCoverLetters is not yet implemented", "",
	))
}

// GetCoverLetter retrieves a single cover letter
// TODO: Implement retrieval logic
func (h *CoverLetterHandler) GetCoverLetter(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.ErrorResponse(
		"NOT_IMPLEMENTED", "GetCoverLetter is not yet implemented", "",
	))
}

// UpdateCoverLetter updates a cover letter
// TODO: Implement update logic
func (h *CoverLetterHandler) UpdateCoverLetter(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.ErrorResponse(
		"NOT_IMPLEMENTED", "UpdateCoverLetter is not yet implemented", "",
	))
}

// PatchCoverLetter partially updates a cover letter
// TODO: Implement patch logic
func (h *CoverLetterHandler) PatchCoverLetter(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.ErrorResponse(
		"NOT_IMPLEMENTED", "PatchCoverLetter is not yet implemented", "",
	))
}

// DeleteCoverLetter deletes a cover letter
// TODO: Implement deletion logic
func (h *CoverLetterHandler) DeleteCoverLetter(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.ErrorResponse(
		"NOT_IMPLEMENTED", "DeleteCoverLetter is not yet implemented", "",
	))
}
