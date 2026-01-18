package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
)

type StockRepository interface {
	FindByID(id uuid.UUID) (*models.Stock, error)
	Create(stock *models.Stock) error
}

type StockRepositoryImpl struct {
	db *gorm.DB // GORM DB instance
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &StockRepositoryImpl{db: db}
}

func (r *StockRepositoryImpl) FindByID(id uuid.UUID) (*models.Stock, error) {
	var s models.Stock
	if err := r.db.First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StockRepositoryImpl) Create(s *models.Stock) error {
	return r.db.Create(s).Error
}
