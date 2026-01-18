package controllers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
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
	Body *models.Stock
}

func (c *StockController) GetStock(ctx context.Context, input *GetStockInput) (*StockResponse, error) {
	_ = ctx

	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid uuid")
	}

	s, err := c.stockService.GetStock(id)
	if err != nil {
		return nil, huma.Error404NotFound(err.Error())
	}

	return &StockResponse{Body: s}, nil
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
	Body struct {
		Status string `json:"status"`
	}
}

func (c *StockController) CreateStock(ctx context.Context, input *CreateStockInput) (*StatusResponse, error) {
	_ = ctx

	if input.Body.Symbol == "" || input.Body.Name == "" {
		return nil, huma.Error400BadRequest("symbol and name required")
	}

	if err := c.stockService.CreateStock(
		input.Body.Symbol,
		input.Body.Name,
		input.Body.Sector,
		input.Body.Price,
	); err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	var res StatusResponse
	res.Body.Status = "created"
	return &res, nil
}
