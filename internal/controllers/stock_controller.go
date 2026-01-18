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
		return &StockResponse{
			Status: http.StatusBadRequest,
			Body:   NewDataResponse[*models.Stock](InvalidStatus("invalid uuid"), nil),
		}, nil
	}

	s, err := c.stockService.GetStock(id)
	if err != nil {
		return &StockResponse{
			Status: http.StatusNotFound,
			Body:   NewDataResponse[*models.Stock](NewStatus("404", err.Error(), nil), nil),
		}, nil
	}

	return &StockResponse{
		Status: http.StatusOK,
		Body:   NewDataResponse(SuccessStatus(), s),
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
		return &StatusResponse{
			Status: http.StatusBadRequest,
			Body:   NewDataResponse[any](InvalidStatus("symbol and name required"), nil),
		}, nil
	}

	if err := c.stockService.CreateStock(
		input.Body.Symbol,
		input.Body.Name,
		input.Body.Sector,
		input.Body.Price,
	); err != nil {
		return &StatusResponse{
			Status: http.StatusInternalServerError,
			Body:   NewDataResponse[any](NewStatus("500", err.Error(), nil), nil),
		}, nil
	}

	return &StatusResponse{
		Status: http.StatusCreated,
		Body:   NewDataResponse[any](SuccessStatus(), nil),
	}, nil
}
