package stock

import (
	domain "project-go-basic/internal/domain/stock"
)

type stockServiceImpl struct {
	stockRepo domain.StockRepository
}

func NewStockServiceImpl(stockRepo domain.StockRepository) StockService {
	return &stockServiceImpl{stockRepo}
}

