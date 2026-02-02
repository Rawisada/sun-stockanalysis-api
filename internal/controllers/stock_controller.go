package controllers

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	common "sun-stockanalysis-api/internal/common"
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
	Body   common.DataResponse[*models.Stock]
}

func (c *StockController) GetStock(ctx context.Context, input *GetStockInput) (*StockResponse, error) {
	_ = ctx

	id, err := uuid.Parse(input.ID)

	if err != nil {
		return nil, common.NewBadRequest("invalid stock id")
	}

	s, err := c.stockService.GetStock(id)

	if err != nil {
		return nil, common.NewNotFound("stock not found")
	}

	return &StockResponse{
		Status: http.StatusOK,
		Body:   common.SuccessResponse(s),
	}, nil
}


func (c *StockController) CreateStock(ctx context.Context, input *stock.CreateStockInput) (*common.StatusResponse, error) {
	_ = ctx

	if err := validateCreateStockInput(input); err != nil {
		return nil, err
	}

	if err := c.stockService.CreateStock(*input); err != nil {
		return nil, common.NewInternalError(err.Error())
	}

	return &common.StatusResponse{
		Status: http.StatusCreated,
		Body:   common.SuccessResponse[any]("stock created successfully"),
	}, nil
}

func validateCreateStockInput(input *stock.CreateStockInput) error {
	if input.Body.Symbol == "" {
		return common.NewBadRequest("symbol required")
	}

	return nil
}
