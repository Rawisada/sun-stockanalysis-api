package controllers

import (
	"context"

	"sun-stockanalysis-api/internal/repository"
)

type HealthController struct {
	healthRepository repository.HealthRepository
	serverVersion    string
}

func NewHealthController(healthRepository repository.HealthRepository, serverVersion string) *HealthController {
	return &HealthController{
		healthRepository: healthRepository,
		serverVersion:    serverVersion,
	}
}

type HealthRequest struct {
}

type HealthResponseBody struct {
	Version string `json:"version"`
}

func (hc *HealthController) Healthz(ctx context.Context, req *EmptyRequest) (*Response[HealthResponseBody], error) {
	return &Response[HealthResponseBody]{
		Body: HealthResponseBody{
			Version: hc.serverVersion,
		},
	}, nil
}

func (hc *HealthController) Readyz(ctx context.Context, req *EmptyRequest) (*EmptyResponse, error) {
	err := hc.healthRepository.CheckDBStatus(ctx)

	if err != nil {
		return nil, err
	}

	return &EmptyResponse{
		Status: 200,
	}, nil
}
