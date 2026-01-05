package stock

import "project-go-basic/internal/domain/stock"

type stockServiceImpl struct {
	repo stock.Repository
}

func NewStockService(repo stock.Repository) StockService {
	return &stockServiceImpl{
		repo: repo,
	}
}

func (s *stockServiceImpl) GetStock(id uint) (*stock.Stock, error) {
	return s.repo.FindByID(id)
}

func (s *stockServiceImpl) CreateStock(symbol, name, market string) error {
	return s.repo.Create(&stock.Stock{
		Symbol: symbol,
		Name:   name,
		Market: market,
	})
}
