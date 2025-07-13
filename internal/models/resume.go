package models

import "time"

type Resume struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	UserID      uint      `json:"user_id" grom:"foreignKey:UserID;constraint:OnUpdate:CASCADE;"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FileType    string    `json:"file_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SoftwareDeveloperEntryLevel struct {
	Resume
	ProgrammingLanguages []string `json:"programming_languages"`
	Frameworks           []string `json:"frameworks"`
	Tools                []string `json:"tools"`
}
