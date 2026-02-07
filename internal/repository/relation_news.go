package repository

import (
	"errors"

	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
)

type RelationNewsRepository interface {
	CreateMany(items []models.RelationNews) error
	ListDistinctRelationSymbols() ([]string, error)
}

type RelationNewsRepositoryImpl struct {
	db *gorm.DB
}

func NewRelationNewsRepository(db *gorm.DB) RelationNewsRepository {
	return &RelationNewsRepositoryImpl{db: db}
}

func (r *RelationNewsRepositoryImpl) CreateMany(items []models.RelationNews) error {
	if len(items) == 0 {
		return errors.New("relation_news items are empty")
	}
	return r.db.Create(&items).Error
}

func (r *RelationNewsRepositoryImpl) ListDistinctRelationSymbols() ([]string, error) {
	var symbols []string
	if err := r.db.Model(&models.RelationNews{}).
		Select("distinct relation_symbol").
		Where("relation_symbol <> ''").
		Pluck("relation_symbol", &symbols).Error; err != nil {
		return nil, err
	}
	return symbols, nil
}
