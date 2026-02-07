package controllers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"sun-stockanalysis-api/internal/domains/stock"
	stockmock "sun-stockanalysis-api/internal/mocks/domains/stock"
	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/pkg/apierror"
)

type StockControllerSuite struct {
	suite.Suite
	stockService *stockmock.MockStockService
	controller   *StockController
}

func (s *StockControllerSuite) SetupTest() {
	s.stockService = stockmock.NewMockStockService(s.T())
	s.controller = NewStockController(s.stockService)
}

func (s *StockControllerSuite) TestGetStock_InvalidID() {
	input := &GetStockInput{ID: "not-uuid"}

	resp, err := s.controller.GetStock(context.Background(), input)

	s.Nil(resp)
	s.Error(err)
	s.Equal(apierror.ErrCodeBadRequest, err.(*apierror.APIError).Code)
}

func (s *StockControllerSuite) TestGetStock_NotFound() {
	id := uuid.New()
	input := &GetStockInput{ID: id.String()}

	s.stockService.EXPECT().GetStock(id).Return((*models.Stock)(nil), errors.New("not found"))

	resp, err := s.controller.GetStock(context.Background(), input)

	s.Nil(resp)
	s.Error(err)
	s.Equal(apierror.ErrCodeNotFound, err.(*apierror.APIError).Code)
}

func (s *StockControllerSuite) TestGetStock_Success() {
	id := uuid.New()
	input := &GetStockInput{ID: id.String()}
	expected := &models.Stock{
		ID:     id,
		Symbol: "AAPL",
		Name:   "Apple Inc.",
	}

	s.stockService.EXPECT().GetStock(id).Return(expected, nil)

	resp, err := s.controller.GetStock(context.Background(), input)

	s.NoError(err)
	s.NotNil(resp)
	s.Equal(http.StatusOK, resp.Status)
	s.Equal(expected, resp.Body.Data)
}

func (s *StockControllerSuite) TestCreateStock_ValidationError() {
	input := &stock.CreateStockInput{}
	input.Body.Symbol = ""

	resp, err := s.controller.CreateStock(context.Background(), input)

	s.Nil(resp)
	s.Error(err)
	s.Equal(apierror.ErrCodeBadRequest, err.(*apierror.APIError).Code)
}

func (s *StockControllerSuite) TestCreateStock_InternalError() {
	input := &stock.CreateStockInput{}
	input.Body.Symbol = "AAPL"

	s.stockService.EXPECT().CreateStock(*input).Return(errors.New("db down"))

	resp, err := s.controller.CreateStock(context.Background(), input)

	s.Nil(resp)
	s.Error(err)
	s.Equal(apierror.ErrCodeInternalError, err.(*apierror.APIError).Code)
}

func (s *StockControllerSuite) TestCreateStock_Success() {
	input := &stock.CreateStockInput{}
	input.Body.Symbol = "AAPL"

	s.stockService.EXPECT().CreateStock(*input).Return(nil)

	resp, err := s.controller.CreateStock(context.Background(), input)

	s.NoError(err)
	s.NotNil(resp)
	s.Equal(http.StatusCreated, resp.Status)
	s.Equal("stock created successfully", resp.Body.Data)
}

func TestStockControllerSuite(t *testing.T) {
	suite.Run(t, new(StockControllerSuite))
}
