package handler

import (
	_ "gorm.io/gorm"
	grom "gorm.io/gorm"
)

type CoverLetterHandler struct {
	db *grom.DB
}

func NewCoverLetter(db *grom.DB) *CoverLetterHandler {
	return &CoverLetterHandler{
		db: db,
	}
}
