package repository

import (
	domain "sun-stockanalysis-api/internal/domain/stock"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type stockRepositoryImpl struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) domain.StockRepository {
	return &stockRepositoryImpl{db: db}
}

func (r *stockRepositoryImpl) FindByID(id uuid.UUID) (*domain.Stock, error) {
	var s domain.Stock
	if err := r.db.First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *stockRepositoryImpl) Create(s *domain.Stock) error {
	return r.db.Create(s).Error
}
