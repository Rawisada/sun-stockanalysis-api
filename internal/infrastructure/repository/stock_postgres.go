package repository

import (
    stock "project-go-basic/internal/domain/stock"
)

type stockRepositoryImpl struct {}

func NewStockRepository() stock.StockRepository {
	return &stockRepositoryImpl{}
}