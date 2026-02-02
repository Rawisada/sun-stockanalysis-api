package stock

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"sun-stockanalysis-api/internal/models"
	repositorymock "sun-stockanalysis-api/internal/mocks/repository"
)

type StockServiceSuite struct {
	suite.Suite
	repo    *repositorymock.MockStockRepository
	service StockService
}

func (s *StockServiceSuite) SetupTest() {
	s.repo = repositorymock.NewMockStockRepository(s.T())
	s.service = NewStockService(s.repo)
}

func (s *StockServiceSuite) TestGetStock_ReturnsStock() {
	id := uuid.New()
	expected := &models.Stock{
		ID:     id,
		Symbol: "AAPL",
		Name:   "Apple Inc.",
	}

	s.repo.EXPECT().FindByID(id).Return(expected, nil)

	result, err := s.service.GetStock(id)

	s.NoError(err)
	s.Equal(expected, result)
}

func (s *StockServiceSuite) TestGetStock_ReturnsError() {
	id := uuid.New()
	wantErr := errors.New("not found")

	s.repo.EXPECT().FindByID(id).Return((*models.Stock)(nil), wantErr)

	result, err := s.service.GetStock(id)

	s.Nil(result)
	s.EqualError(err, wantErr.Error())
}

func (s *StockServiceSuite) TestCreateStock_Persists() {
	input := CreateStockInput{}
	input.Body.Symbol = "TSLA"
	input.Body.Name = "Tesla, Inc."
	input.Body.Sector = "Automotive"
	input.Body.Exchange = "NASDAQ"
	input.Body.AssetType = "Stock"
	input.Body.Currency = "USD"

	s.repo.EXPECT().Create(mock.MatchedBy(func(stock *models.Stock) bool {
		s.NotNil(stock)
		return stock.Symbol == input.Body.Symbol &&
			stock.Name == input.Body.Name &&
			stock.Sector == input.Body.Sector &&
			stock.Exchange == input.Body.Exchange &&
			stock.AssetType == input.Body.AssetType &&
			stock.Currency == input.Body.Currency
	})).Return(nil)

	err := s.service.CreateStock(input)

	s.NoError(err)
}

func (s *StockServiceSuite) TestCreateStock_ReturnsError() {
	input := CreateStockInput{}
	input.Body.Symbol = "NVDA"
	input.Body.Name = "NVIDIA Corporation"
	input.Body.Sector = "Technology"
	input.Body.Exchange = "NASDAQ"
	input.Body.AssetType = "Stock"
	input.Body.Currency = "USD"

	wantErr := errors.New("create failed")

	s.repo.EXPECT().Create(mock.Anything).Return(wantErr)

	err := s.service.CreateStock(input)

	s.EqualError(err, wantErr.Error())
}

func TestStockServiceSuite(t *testing.T) {
	suite.Run(t, new(StockServiceSuite))
}
