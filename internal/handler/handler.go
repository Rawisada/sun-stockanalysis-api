package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"sun-stockanalysis-api/internal/controllers"
)

const apiBasePath = "/v1"

func v1Tags() []string {
	return []string{"v1"}
}

func RegisterRoutes(rootApi huma.API, controllers *controllers.Controllers) {
	rootApi.UseMiddleware(requestIDMiddleware)

	registerHealthHandlers(rootApi, controllers)
	v1Api := huma.NewGroup(rootApi, apiBasePath)

	registerV1(v1Api, controllers)
}

func registerHealthHandlers(api huma.API, controllers *controllers.Controllers) {
	tags := []string{"Health Check"}

	huma.Register(api, huma.Operation{
		Path:   "/healthz",
		Method: http.MethodGet,
		Tags:   tags,
	}, controllers.HealthController.Healthz)

	huma.Register(api, huma.Operation{
		Path:   "/readyz",
		Method: http.MethodGet,
		Tags:   tags,
	}, controllers.HealthController.Readyz)
}

func registerV1(api huma.API, controllers *controllers.Controllers) {
	huma.Register(api, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/stocks/{id}",
		Summary: "Get stock by ID",
		Tags:    v1Tags(),
	}, controllers.StockController.GetStock)

	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/stocks",
		Summary:       "Create stock",
		Tags:          v1Tags(),
		DefaultStatus: http.StatusCreated,
	}, controllers.StockController.CreateStock)
}

func requestIDMiddleware(ctx huma.Context, next func(huma.Context)) {
	if ctx.Header("X-Request-Id") == "" {
		ctx.SetHeader("X-Request-Id", uuid.NewString())
	}
	next(ctx)
}
