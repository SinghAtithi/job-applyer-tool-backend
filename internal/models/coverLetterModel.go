package models

type JobDescriptionTable struct {
	URL         string `json:"url"`
	ContentInfo string `json:"contentInfo"`
	TextContent string `json:"textContent"`
	UserName    string `json:"userName"`
}

type CoverLetterTable struct {
	URL         string `json:"url" gorm:"primaryKey"`
	ContentInfo string `json:"contentInfo"`
}
