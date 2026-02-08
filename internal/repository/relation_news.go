package repository

import (
	"errors"

	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
)

type RelationNewsRepository interface {
	CreateMany(items []models.RelationNews) error
	ListDistinctRelationSymbols() ([]string, error)
	ListRelationSymbolsBySymbol(symbol string) ([]string, error)
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

func (r *RelationNewsRepositoryImpl) ListRelationSymbolsBySymbol(symbol string) ([]string, error) {
	if symbol == "" {
		return nil, errors.New("symbol is empty")
	}
	var symbols []string
	if err := r.db.Model(&models.RelationNews{}).
		Select("distinct relation_symbol").
		Where("symbol = ? AND relation_symbol <> ''", symbol).
		Pluck("relation_symbol", &symbols).Error; err != nil {
		return nil, err
	}
	return symbols, nil
}
