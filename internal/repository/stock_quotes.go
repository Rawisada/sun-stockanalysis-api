package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
)

type StockQuoteRepository interface {
	Create(quote *models.StockQuote) error
	FindLatestBySymbol(symbol string) (*models.StockQuote, error)
	FindLatestBySymbolBetween(symbol string, start, end time.Time, limit int) ([]models.StockQuote, error)
	FindBySymbolBetween(symbol string, start, end time.Time) ([]models.StockQuote, error)
}

type StockQuoteRepositoryImpl struct {
	db *gorm.DB
}

func NewStockQuoteRepository(db *gorm.DB) StockQuoteRepository {
	return &StockQuoteRepositoryImpl{db: db}
}

func (r *StockQuoteRepositoryImpl) Create(quote *models.StockQuote) error {
	if quote == nil {
		return errors.New("stock quote is nil")
	}
	return r.db.Create(quote).Error
}

func (r *StockQuoteRepositoryImpl) FindLatestBySymbol(symbol string) (*models.StockQuote, error) {
	if symbol == "" {
		return nil, errors.New("symbol is empty")
	}
	var quote models.StockQuote
	if err := r.db.
		Where("symbol = ?", symbol).
		Order("created_at desc").
		First(&quote).Error; err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *StockQuoteRepositoryImpl) FindLatestBySymbolBetween(symbol string, start, end time.Time, limit int) ([]models.StockQuote, error) {
	if symbol == "" {
		return nil, errors.New("symbol is empty")
	}
	if limit <= 0 {
		return []models.StockQuote{}, nil
	}
	var quotes []models.StockQuote
	if err := r.db.
		Where("symbol = ? AND created_at >= ? AND created_at <= ?", symbol, start, end).
		Order("created_at desc").
		Limit(limit).
		Find(&quotes).Error; err != nil {
		return nil, err
	}
	return quotes, nil
}

func (r *StockQuoteRepositoryImpl) FindBySymbolBetween(symbol string, start, end time.Time) ([]models.StockQuote, error) {
	if symbol == "" {
		return nil, errors.New("symbol is empty")
	}
	var quotes []models.StockQuote
	if err := r.db.
		Where("symbol = ? AND created_at >= ? AND created_at <= ?", symbol, start, end).
		Order("created_at asc").
		Find(&quotes).Error; err != nil {
		return nil, err
	}
	return quotes, nil
}
