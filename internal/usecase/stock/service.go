package stock

import domain "sun-stockanalysis-api/internal/domain/stock"

import "github.com/google/uuid"

type StockService interface {
	GetStock(id uuid.UUID) (*domain.Stock, error)
	CreateStock(symbol, name, sector string, price int) error
}