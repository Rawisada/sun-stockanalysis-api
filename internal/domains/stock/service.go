package stock

import (
	"github.com/google/uuid"

	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/repository"
)

type StockService interface {
	GetStock(id uuid.UUID) (*models.Stock, error)
	CreateStock(input CreateStockInput) error
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

func (s *StockServiceImpl) CreateStock(input CreateStockInput) error {


	return s.repo.Create(&models.Stock{
		Symbol:      input.Body.Symbol,
		Name:        input.Body.Name,
		Sector:      input.Body.Sector,
		Price:       input.Body.Price,
		Exchange:    input.Body.Exchange,
		AssetType:   input.Body.AssetType,
		Currency:    input.Body.Currency,
	})
}
