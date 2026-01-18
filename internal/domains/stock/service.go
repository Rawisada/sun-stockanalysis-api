package stock

import (
	"github.com/google/uuid"

	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/repository"
)

type StockService interface {
	GetStock(id uuid.UUID) (*models.Stock, error)
	CreateStock(symbol, name, sector string, price int) error
}

type StockServiceImpl struct {
	repo repository.StockRepository
}

func NewStockService(repo repository.StockRepository) StockService {
	return &StockServiceImpl{repo: repo}
}

func (s *StockServiceImpl) GetStock(id uuid.UUID) (*models.Stock, error) {
	return s.repo.FindByID(id)
}

func (s *StockServiceImpl) CreateStock(symbol, name, sector string, price int) error {
	return s.repo.Create(&models.Stock{
		// ไม่ต้อง set ID เพราะ DB default gen_random_uuid() จะสร้างให้
		Symbol: symbol,
		Name:   name,
		Sector: sector,
		Price:  price,
	})
}
