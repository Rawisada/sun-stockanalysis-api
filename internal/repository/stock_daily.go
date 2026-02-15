package repository

import (
	"errors"

	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
)

type StockDailyRepository interface {
	Create(metric *models.StockDaily) error
	FindLatestBySymbol(symbol string) (*models.StockDaily, error)
	FindPreviousBySymbol(symbol string) ([]models.StockDaily, error)
	FindBySymbol(symbol string) ([]models.StockDaily, error)
}

type StockDailyRepositoryImpl struct {
	db *gorm.DB
}

func NewStockDailyRepository(db *gorm.DB) StockDailyRepository {
	return &StockDailyRepositoryImpl{db: db}
}

func (r *StockDailyRepositoryImpl) Create(metric *models.StockDaily) error {
	if metric == nil {
		return errors.New("stock daily is nil")
	}
	return r.db.Create(metric).Error
}

func (r *StockDailyRepositoryImpl) FindLatestBySymbol(symbol string) (*models.StockDaily, error) {
	if symbol == "" {
		return nil, errors.New("symbol is empty")
	}
	var metric models.StockDaily
	if err := r.db.
		Where("symbol = ?", symbol).
		Order("created_at desc").
		First(&metric).Error; err != nil {
		return nil, err
	}
	return &metric, nil
}

func (r *StockDailyRepositoryImpl) FindBySymbol(symbol string) ([]models.StockDaily, error) {
	if symbol == "" {
		return nil, errors.New("symbol is empty")
	}
	var metrics []models.StockDaily
	if err := r.db.
		Where("symbol = ?", symbol).
		Order("created_at desc").
		Find(&metrics).Error; err != nil {
		return nil, err
	}
	return metrics, nil
}

func (r *StockDailyRepositoryImpl) FindPreviousBySymbol(symbol string) ([]models.StockDaily, error) {
	if symbol == "" {
		return nil, errors.New("symbol is empty")
	}
	var metrics []models.StockDaily
	if err := r.db.
		Where("symbol = ?", symbol).
		Order("created_at desc").
		Offset(1).
		Limit(1).
		Find(&metrics).Error; err != nil {
		return nil, err
	}
	return metrics, nil
}
