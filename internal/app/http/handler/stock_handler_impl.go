package handler

import (
	stock "project-go-basic/internal/usecase/stock"
)

type stockHandlerImpl struct {
	stockService stock.StockService
}

func NewStockHandlerImpl(stockService stock.StockService) StockHandler {
	return  &stockHandlerImpl{stockService}
}