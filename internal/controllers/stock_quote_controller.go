package controllers

import (
	"context"
	"net/http"

	"sun-stockanalysis-api/internal/domains/stock_quotes"
	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/pkg/apierror"
	"sun-stockanalysis-api/pkg/response"
)

type StockQuoteController struct {
	service stock_quotes.StockQuoteService
}

func NewStockQuoteController(service stock_quotes.StockQuoteService) *StockQuoteController {
	return &StockQuoteController{service: service}
}

type StockQuoteListResponse struct {
	Status int `status:"default"`
	Body   response.ApiResponse[[]models.StockQuote]
}

type StockQuoteListInput struct {
	Symbol string `query:"symbol" doc:"Filter by symbol"`
}

func (c *StockQuoteController) ListAll(ctx context.Context, input *StockQuoteListInput) (*StockQuoteListResponse, error) {
	symbol := ""
	if input != nil {
		symbol = input.Symbol
	}
	quotes, err := c.service.List(ctx, symbol)
	if err != nil {
		return nil, apierror.NewInternalError(err.Error())
	}

	return &StockQuoteListResponse{
		Status: http.StatusOK,
		Body:   response.Success(quotes),
	}, nil
}
