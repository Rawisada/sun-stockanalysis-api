package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"sun-stockanalysis-api/internal/controllers"
)

func RegisterStockRoutes(api huma.API, controllers *controllers.Controllers, middleware func(huma.Context, func(huma.Context))) {
	protected := huma.NewGroup(api, "")
	protected.UseMiddleware(middleware)

	huma.Register(protected, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/stocks/{id}",
		Summary: "Get stock by ID",
		Tags:    v1Tags(),
	}, controllers.StockController.GetStock)

	huma.Register(protected, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/stocks",
		Summary:       "Create stock",
		Tags:          v1Tags(),
		DefaultStatus: http.StatusCreated,
	}, controllers.StockController.CreateStock)
}
