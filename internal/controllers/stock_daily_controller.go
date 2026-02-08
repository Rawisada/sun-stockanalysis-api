package controllers

import (
	"context"
	"net/http"

	"sun-stockanalysis-api/internal/domains/stock_daily"
	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/pkg/apierror"
	"sun-stockanalysis-api/pkg/response"
)

type StockDailyController struct {
	service stock_daily.StockDailyService
}

func NewStockDailyController(service stock_daily.StockDailyService) *StockDailyController {
	return &StockDailyController{service: service}
}

type StockDailyListInput struct {
	Symbol string `query:"symbol" doc:"Filter by symbol" required:"true"`
}

type StockDailyListResponse struct {
	Status int `status:"default"`
	Body   response.ApiResponse[[]models.StockDaily]
}

func (c *StockDailyController) ListBySymbol(ctx context.Context, input *StockDailyListInput) (*StockDailyListResponse, error) {
	if input == nil || input.Symbol == "" {
		return nil, apierror.NewBadRequest("symbol required")
	}

	metrics, err := c.service.ListBySymbol(ctx, input.Symbol)
	if err != nil {
		return nil, apierror.NewInternalError(err.Error())
	}

	return &StockDailyListResponse{
		Status: http.StatusOK,
		Body:   response.Success(metrics),
	}, nil
}
