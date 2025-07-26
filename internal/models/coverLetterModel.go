package models

type CoverLetterTable struct {
	URL         string `json:"url" gorm:"primaryKey"`
	ContentInfo string `json:"contentInfo"`
	UserName    string `json:"user_name"`
}
