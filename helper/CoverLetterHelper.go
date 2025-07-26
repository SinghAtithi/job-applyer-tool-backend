package helper

import (
	"example.com/internal/models"
	"fmt"
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

func (cl *CoverLetterHelper) GetCoverLetterFromDB(url string) (models.CoverLetterTable, bool) {
	var coverLetters models.CoverLetterTable

	query := cl.db.Model(&models.CoverLetterTable{})

	query = query.Where("url = ?", url)
	fmt.Printf("Searching for URL: %s\n", url)

	var count int64
	cl.db.Model(&models.CoverLetterTable{}).Count(&count)
	fmt.Printf("Total records in table: %d\n", count)

	result := query.First(&coverLetters)

	cl.db = cl.db.Debug()
	if result.Error != nil {
		return models.CoverLetterTable{}, false
	}

	return coverLetters, true
}

func (cl *CoverLetterHelper) CommitJobDescription(content string, url string) {

	coverLetterData := models.CoverLetterTable{
		URL:         url,
		ContentInfo: content,
	}

	result := cl.db.Create(&coverLetterData)
	if result.Error != nil {
		log.Printf("ERRRRRRRRRRRROOORRRR" + result.Error.Error())
	}
}
