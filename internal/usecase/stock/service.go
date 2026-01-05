package stock

import stock "project-go-basic/internal/domain/stock" 

type StockService interface {
	GetStock(id uint) (*stock.Stock, error)
	CreateStock(symbol, name, market string) error
}