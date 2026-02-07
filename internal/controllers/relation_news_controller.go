package controllers

import (
	"context"
	"net/http"

	"sun-stockanalysis-api/internal/domains/relation_news"
	"sun-stockanalysis-api/pkg/apierror"
	"sun-stockanalysis-api/pkg/response"
)

type RelationNewsController struct {
	service relation_news.RelationNewsService
}

func NewRelationNewsController(service relation_news.RelationNewsService) *RelationNewsController {
	return &RelationNewsController{service: service}
}

type CreateRelationNewsResponse struct {
	Status int `status:"default"`
	Body   response.ApiResponse[any]
}

func (c *RelationNewsController) Create(ctx context.Context, input *relation_news.CreateRelationNewsInput) (*CreateRelationNewsResponse, error) {
	_ = ctx

	if input == nil {
		return nil, apierror.NewBadRequest("request body required")
	}
	if input.Body.Symbol == "" {
		return nil, apierror.NewBadRequest("symbol required")
	}

	if err := c.service.CreateRelations(input.Body.Symbol, input.Body.RelationSymbols); err != nil {
		return nil, apierror.NewInternalError(err.Error())
	}

	return &CreateRelationNewsResponse{
		Status: http.StatusCreated,
		Body:   response.Success[any]("relation news created successfully"),
	}, nil
}
