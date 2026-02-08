package repository

import (
	"errors"

	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
)

type AlertEventRepository interface {
	Create(event *models.AlertEvent) error
}

type AlertEventRepositoryImpl struct {
	db *gorm.DB
}

func NewAlertEventRepository(db *gorm.DB) AlertEventRepository {
	return &AlertEventRepositoryImpl{db: db}
}

func (r *AlertEventRepositoryImpl) Create(event *models.AlertEvent) error {
	if event == nil {
		return errors.New("alert event is nil")
	}
	return r.db.Create(event).Error
}
