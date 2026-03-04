package database

import (
	"example.com/internal/models"
	"gorm.io/gorm"
)

type ResumeRepository interface {
	GetByID(id uint) (*models.Resume, error)
	Create(resume *models.Resume) error
	Update(resume *models.Resume) error
	Delete(id uint) error
	IsResumeAvailable(id uint) (bool, error)
}

type resumeRepository struct {
	db *gorm.DB
}

// Fixed: Return interface type, not concrete type
func NewResumeRepository(db *gorm.DB) ResumeRepository {
	return &resumeRepository{db: db} // Fixed: Return pointer
}

func (r *resumeRepository) GetByID(id uint) (*models.Resume, error) {
	var resume models.Resume
	err := r.db.First(&resume, id).Error
	return &resume, err
}

func (r *resumeRepository) Create(resume *models.Resume) error {
	return r.db.Create(resume).Error
}

// Fixed: Added missing Update method
func (r *resumeRepository) Update(resume *models.Resume) error {
	return r.db.Save(resume).Error
}

// Fixed: Added missing Delete method
func (r *resumeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Resume{}, id).Error
}

func (r *resumeRepository) IsResumeAvailable(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Resume{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
