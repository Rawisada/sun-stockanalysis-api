package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
)

type CompanyNewsRepository interface {
	CreateMany(items []models.CompanyNews) error
	FindBySymbolsAndDate(symbols []string, start, end time.Time) ([]models.CompanyNews, error)
	DeleteBefore(t time.Time) error
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

func (r *CompanyNewsRepositoryImpl) FindBySymbolsAndDate(symbols []string, start, end time.Time) ([]models.CompanyNews, error) {
	if len(symbols) == 0 {
		return []models.CompanyNews{}, nil
	}
	var items []models.CompanyNews
	if err := r.db.
		Where("symbol IN ? AND created_at >= ? AND created_at <= ?", symbols, start, end).
		Order("created_at desc").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CompanyNewsRepositoryImpl) DeleteBefore(t time.Time) error {
	if t.IsZero() {
		return nil
	}
	return r.db.
		Where("created_at < ?", t).
		Delete(&models.CompanyNews{}).Error
}
