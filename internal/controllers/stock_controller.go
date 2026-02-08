package controllers

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"sun-stockanalysis-api/internal/domains/stock"
	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/pkg/apierror"
	"sun-stockanalysis-api/pkg/response"
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
	Body   response.ApiResponse[*models.Stock]
}

func (c *StockController) GetStock(ctx context.Context, input *GetStockInput) (*StockResponse, error) {
	_ = ctx

	id, err := uuid.Parse(input.ID)

	if err != nil {
		return nil, apierror.NewBadRequest("invalid stock id")
	}

	s, err := c.stockService.GetStock(id)

	if err != nil {
		return nil, apierror.NewNotFound("stock not found")
	}

	return &StockResponse{
		Status: http.StatusOK,
		Body:   response.Success(s),
	}, nil
}

type StockListResponse struct {
	Status int `status:"default"`
	Body   response.ApiResponse[[]models.Stock]
}

func (c *StockController) ListStocks(ctx context.Context, input *EmptyRequest) (*StockListResponse, error) {
	_ = ctx
	_ = input

	stocks, err := c.stockService.ListAll()
	if err != nil {
		return nil, apierror.NewInternalError(err.Error())
	}

	return &StockListResponse{
		Status: http.StatusOK,
		Body:   response.Success(stocks),
	}, nil
}

type StockCreateResponse struct {
	Status int `status:"default"`
	Body   response.ApiResponse[any]
}

func (c *StockController) CreateStock(ctx context.Context, input *stock.CreateStockInput) (*StockCreateResponse, error) {
	_ = ctx

	if err := validateCreateStockInput(input); err != nil {
		return nil, err
	}

	if err := c.stockService.CreateStock(*input); err != nil {
		return nil, apierror.NewInternalError(err.Error())
	}

	return &StockCreateResponse{
		Status: http.StatusCreated,
		Body:   response.Success[any]("stock created successfully"),
	}, nil
}

func validateCreateStockInput(input *stock.CreateStockInput) error {
	if input.Body.Symbol == "" {
		return apierror.NewBadRequest("symbol required")
	}

	return nil
}
