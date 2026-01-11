package stock

import (
	domain "sun-stockanalysis-api/internal/domain/stock"
	"github.com/google/uuid"
)

type stockServiceImpl struct {
	repo domain.StockRepository
}

func NewStockService(repo domain.StockRepository) StockService {
	return &stockServiceImpl{repo: repo}
}

func (s *stockServiceImpl) GetStock(id uuid.UUID) (*domain.Stock, error) {
	return s.repo.FindByID(id)
}

func (s *stockServiceImpl) CreateStock(symbol, name, sector string, price int) error {
	return s.repo.Create(&domain.Stock{
		// ไม่ต้อง set ID เพราะ DB default gen_random_uuid() จะสร้างให้
		Symbol: symbol,
		Name:   name,
		Sector: sector,
		Price:  price,
	})
}
