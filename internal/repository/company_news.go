package repository

import (
	"errors"

	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
)

type CompanyNewsRepository interface {
	CreateMany(items []models.CompanyNews) error
}

type CompanyNewsRepositoryImpl struct {
	db *gorm.DB
}

func NewCompanyNewsRepository(db *gorm.DB) CompanyNewsRepository {
	return &CompanyNewsRepositoryImpl{db: db}
}

func (r *CompanyNewsRepositoryImpl) CreateMany(items []models.CompanyNews) error {
	if len(items) == 0 {
		return errors.New("company_news items are empty")
	}
	return r.db.Create(&items).Error
}
