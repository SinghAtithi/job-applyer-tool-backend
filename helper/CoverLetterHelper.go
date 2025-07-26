package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"example.com/agent"
	"example.com/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	grom "gorm.io/gorm"
	"log"
)

type CoverLetterHelper struct {
	db *grom.DB
}

func NewCoverLetterHelper(db *grom.DB) *CoverLetterHelper {
	return &CoverLetterHelper{
		db: db,
	}
}

func (cl *CoverLetterHelper) GetJobDescriptionFromDB(url string) (models.JobDescriptionTable, bool) {
	var coverLetters models.JobDescriptionTable

	query := cl.db.Model(&models.JobDescriptionTable{})

	query = query.Where("url = ?", url)

	result := query.First(&coverLetters)

	if result.Error != nil {
		return models.JobDescriptionTable{}, false
	}

	return coverLetters, true
}

func (cl *CoverLetterHelper) CommitJobDescription(content string, url string) {

	coverLetterData := models.JobDescriptionTable{
		URL:         url,
		ContentInfo: content,
	}

	result := cl.db.Create(&coverLetterData)
	if result.Error != nil {
		log.Printf("ERRRRRRRRRRRROOORRRR" + result.Error.Error())
	}
}

func (cl *CoverLetterHelper) GetUserDetails(userName string) (models.Resume, error) {
	var userData models.Resume
	query := cl.db.Model(&models.Resume{})
	result := query.Where("user_name = ?", userName).First(&userData)

	if result.Error != nil {
		return models.Resume{}, errors.New("user not found")
	}
	return userData, nil
}

func (cl *CoverLetterHelper) GetCoverLetter(jobDescription string, userName string, url string) (models.CoverLetterTable, error) {
	var coverLetterData models.CoverLetterTable
	query := cl.db.Model(&models.CoverLetterTable{})
	result := query.Where("url = ?", url).Where("user_name = ?", userName).First(&coverLetterData)
	if result.Error == nil && coverLetterData.ContentInfo != "" {
		return coverLetterData, nil
	}

	userDetails, err := cl.GetUserDetails(userName)
	if err != nil {
		return models.CoverLetterTable{}, errors.New("user not found")
	}

	userDetailsJSON, err := json.Marshal(userDetails)
	if err != nil {
		return models.CoverLetterTable{}, errors.New("user details marshal failed")
	}
	userDetailsString := string(userDetailsJSON)

	coverLetterDataString, err := agent.GetCoverLetterString(jobDescription, userDetailsString)

	coverLetterData = cl.CommitCoverLetter(models.CoverLetterTable{
		URL:         url,
		ContentInfo: coverLetterDataString,
		UserName:    userName,
	})

	return coverLetterData, nil
}

func (cl *CoverLetterHelper) CommitCoverLetter(coverLetter models.CoverLetterTable) models.CoverLetterTable {
	result := cl.db.Create(&coverLetter)
	if result.Error != nil {
		log.Printf("ERRRRRRRRRRRROOORRRR" + result.Error.Error())
	}
	return coverLetter
}

func (cl *CoverLetterHelper) SetHeadersForPDf(context *gin.Context, pdfBytes []byte) {

	// Set response headers for PDF download
	context.Header("Content-Type", "application/pdf")
	context.Header("Content-Disposition", "attachment; filename=document.pdf")
	context.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

}

func (cl *CoverLetterHelper) ReturnPDFResponse(content string) ([]byte, error) {
	// Generate PDF from string
	pdfBytes, err := stringToPDF(content)
	if err != nil {
		return pdfBytes, errors.New("failed to generate PDF")
	}
	return pdfBytes, nil
}

func stringToPDF(content string) ([]byte, error) {
	// Create new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "", 12)

	// Add content to PDF with word wrapping
	pdf.MultiCell(190, 8, content, "", "", false)

	// Output PDF to bytes buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
