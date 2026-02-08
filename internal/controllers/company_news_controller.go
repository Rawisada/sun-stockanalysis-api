package controllers

import (
	"context"
	"net/http"
	"time"

	"sun-stockanalysis-api/internal/domains/company_news"
	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/pkg/apierror"
	"sun-stockanalysis-api/pkg/response"
)

type CompanyNewsController struct {
	service company_news.CompanyNewsService
}

func NewCompanyNewsController(service company_news.CompanyNewsService) *CompanyNewsController {
	return &CompanyNewsController{service: service}
}

type CompanyNewsListInput struct {
	Symbol string `query:"symbol" doc:"Base symbol" required:"true"`
	Start  string `query:"start" doc:"Start date (YYYY-MM-DD)" required:"true"`
	End    string `query:"end" doc:"End date (YYYY-MM-DD)" required:"true"`
}

type CompanyNewsListResponse struct {
	Status int `status:"default"`
	Body   response.ApiResponse[[]models.CompanyNews]
}

func (c *CompanyNewsController) ListBySymbolAndDate(ctx context.Context, input *CompanyNewsListInput) (*CompanyNewsListResponse, error) {
	if input == nil || input.Symbol == "" || input.Start == "" || input.End == "" {
		return nil, apierror.NewBadRequest("symbol, start, end are required")
	}
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	startDate, err := time.ParseInLocation("2006-01-02", input.Start, loc)
	if err != nil {
		return nil, apierror.NewBadRequest("invalid start date format (YYYY-MM-DD)")
	}
	endDate, err := time.ParseInLocation("2006-01-02", input.End, loc)
	if err != nil {
		return nil, apierror.NewBadRequest("invalid end date format (YYYY-MM-DD)")
	}
	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, loc)

	items, err := c.service.ListBySymbolAndDate(ctx, input.Symbol, start, end)
	if err != nil {
		return nil, apierror.NewInternalError(err.Error())
	}

	return &CompanyNewsListResponse{
		Status: http.StatusOK,
		Body:   response.Success(items),
	}, nil
}
