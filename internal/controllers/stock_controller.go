package controllers

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"sun-stockanalysis-api/internal/domains/stock"
	"sun-stockanalysis-api/internal/models"
)

type StockController struct {
	stockService stock.StockService
}

func NewStockController(stockService stock.StockService) *StockController {
	return &StockController{stockService: stockService}
}

type GetStockInput struct {
	ID string `path:"id" doc:"Stock ID (UUID)"`
}

type StockResponse struct {
	Status int `status:"default"`
	Body   DataResponse[*models.Stock]
}

func (c *StockController) GetStock(ctx context.Context, input *GetStockInput) (*StockResponse, error) {
	_ = ctx

	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, NewBadRequest("invalid stock id")
	}

	s, err := c.stockService.GetStock(id)
	if err != nil {
		return nil, NewNotFound(err.Error())
	}

	return &StockResponse{
		Status: http.StatusOK,
		Body:   SuccessResponse(s),
	}, nil
}

type CreateStockInput struct {
	Body struct {
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
		Sector string `json:"sector"`
		Price  int    `json:"price"`
	}
}

type StatusResponse struct {
	Status int `status:"default"`
	Body   DataResponse[any]
}

func (c *StockController) CreateStock(ctx context.Context, input *CreateStockInput) (*StatusResponse, error) {
	_ = ctx

	if input.Body.Symbol == "" || input.Body.Name == "" {
		return nil, NewBadRequest("symbol and name required")
	}

	if err := c.stockService.CreateStock(
		input.Body.Symbol,
		input.Body.Name,
		input.Body.Sector,
		input.Body.Price,
	); err != nil {
		return nil, NewInternalError(err.Error())
	}

	return &StatusResponse{
		Status: http.StatusCreated,
		Body:   SuccessResponse[any]("Stock created successfully"),
	}, nil
}
