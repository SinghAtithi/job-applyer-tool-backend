package handler

import (
	"example.com/agent"
	"example.com/helper"
	"example.com/internal/models"
	"github.com/gin-gonic/gin"
	_ "gorm.io/gorm"
	grom "gorm.io/gorm"
	"net/http"
)

type CoverLetterHandler struct {
	db *grom.DB
}

func (h CoverLetterHandler) CreateCoverLetter(context *gin.Context) {

	var requestBody models.JobDescriptionTable

	if err := context.BindJSON(&requestBody); err != nil || requestBody.URL == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coverLetterHelper := helper.NewCoverLetterHelper(h.db)

	jobDescriptionData, isJobDescriptionPresent := coverLetterHelper.GetJobDescriptionFromDB(requestBody.URL)

	var jobDescriptionContent string
	var err error
	if isJobDescriptionPresent {
		jobDescriptionContent = jobDescriptionData.ContentInfo
	} else {
		// Make the agent call and get the cover letter
		jobDescriptionContent, err = agent.GetJobDescription(requestBody.TextContent)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		coverLetterHelper.CommitJobDescription(jobDescriptionContent, requestBody.URL)
	}
	//userDetails, err = coverLetterHelper.GetUserDetails(requestBody.UserName)
	coverLetter, err := coverLetterHelper.GetCoverLetter(jobDescriptionContent, requestBody.UserName, requestBody.URL)

	if err != nil || coverLetter.ContentInfo == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Unable to createCoverLetter"})
		return
	}

	pdfBytes, err := coverLetterHelper.ReturnPDFResponse(coverLetter.ContentInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	coverLetterHelper.SetHeadersForPDf(context, pdfBytes)

	// Return PDF as response
	context.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h CoverLetterHandler) GetAllCoverLetters(context *gin.Context) {

}

func (h CoverLetterHandler) GetCoverLetter(context *gin.Context) {
	context.JSON(http.StatusCreated, models.JobDescriptionTable{URL: "example.com", ContentInfo: "Example Content"})
}

func (h CoverLetterHandler) UpdateCoverLetter(context *gin.Context) {

}

func (h CoverLetterHandler) PatchCoverLetter(context *gin.Context) {

}

func (h CoverLetterHandler) DeleteCoverLetter(context *gin.Context) {

}

func NewCoverLetter(db *grom.DB) *CoverLetterHandler {
	return &CoverLetterHandler{db: db}
}
