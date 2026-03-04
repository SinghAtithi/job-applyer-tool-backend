package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"example.com/agent"
	"example.com/internal/models"
	"example.com/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

// CoverLetterHelper provides helper methods for cover letter operations
type CoverLetterHelper struct {
	db *gorm.DB
}

// NewCoverLetterHelper creates a new helper instance
func NewCoverLetterHelper(db *gorm.DB) *CoverLetterHelper {
	return &CoverLetterHelper{db: db}
}

// GetJobDescriptionFromDB retrieves a cached job description by URL
func (cl *CoverLetterHelper) GetJobDescriptionFromDB(url string) (models.JobDescriptionTable, bool) {
	var jobDesc models.JobDescriptionTable

	result := cl.db.Model(&models.JobDescriptionTable{}).Where("url = ?", url).First(&jobDesc)
	if result.Error != nil {
		return models.JobDescriptionTable{}, false
	}

	return jobDesc, true
}

// CommitJobDescription saves a job description to the database.
// Returns an error instead of silently logging on failure.
func (cl *CoverLetterHelper) CommitJobDescription(content string, url string) error {
	data := models.JobDescriptionTable{
		URL:         url,
		ContentInfo: content,
	}

	result := cl.db.Create(&data)
	if result.Error != nil {
		logger.Error("failed to save job description: %v", result.Error)
		return result.Error
	}
	return nil
}

// GetUserDetails retrieves a user's resume by username
func (cl *CoverLetterHelper) GetUserDetails(userName string) (models.Resume, error) {
	var userData models.Resume

	result := cl.db.Model(&models.Resume{}).Where("user_name = ?", userName).First(&userData)
	if result.Error != nil {
		return models.Resume{}, errors.New("user not found")
	}
	return userData, nil
}

// GetCoverLetter retrieves an existing cover letter or generates a new one
func (cl *CoverLetterHelper) GetCoverLetter(jobDescription string, userName string, url string) (models.CoverLetterTable, error) {
	// Check for existing cover letter
	var coverLetterData models.CoverLetterTable
	result := cl.db.Model(&models.CoverLetterTable{}).
		Where("url = ?", url).
		Where("user_name = ?", userName).
		First(&coverLetterData)

	if result.Error == nil && coverLetterData.ContentInfo != "" {
		logger.Info("found cached cover letter for user=%s url=%s", userName, url)
		return coverLetterData, nil
	}

	// Fetch user details
	userDetails, err := cl.GetUserDetails(userName)
	if err != nil {
		return models.CoverLetterTable{}, fmt.Errorf("user not found: %w", err)
	}

	// Serialize user details for the AI agent
	userDetailsJSON, err := json.Marshal(userDetails)
	if err != nil {
		return models.CoverLetterTable{}, fmt.Errorf("failed to marshal user details: %w", err)
	}

	// Generate cover letter via AI
	coverLetterContent, err := agent.GetCoverLetterString(jobDescription, string(userDetailsJSON))
	if err != nil {
		return models.CoverLetterTable{}, fmt.Errorf("failed to generate cover letter: %w", err)
	}

	// Persist the generated cover letter
	coverLetterData = cl.CommitCoverLetter(models.CoverLetterTable{
		URL:         url,
		ContentInfo: coverLetterContent,
		UserName:    userName,
	})

	return coverLetterData, nil
}

// CommitCoverLetter saves a cover letter to the database
func (cl *CoverLetterHelper) CommitCoverLetter(coverLetter models.CoverLetterTable) models.CoverLetterTable {
	result := cl.db.Create(&coverLetter)
	if result.Error != nil {
		logger.Error("failed to save cover letter: %v", result.Error)
	}
	return coverLetter
}

// SetHeadersForPDF sets the response headers for PDF download
func (cl *CoverLetterHelper) SetHeadersForPDF(c *gin.Context, pdfBytes []byte) {
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=document.pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
}

// ReturnPDFResponse generates a PDF from the given content string
func (cl *CoverLetterHelper) ReturnPDFResponse(content string) ([]byte, error) {
	pdfBytes, err := stringToPDF(content)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	return pdfBytes, nil
}

// stringToPDF converts a plain-text string into a simple A4 PDF document
func stringToPDF(content string) ([]byte, error) {
	doc := gofpdf.New("P", "mm", "A4", "")
	doc.AddPage()
	doc.SetFont("Arial", "", 12)
	doc.MultiCell(190, 8, content, "", "", false)

	var buf bytes.Buffer
	if err := doc.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
