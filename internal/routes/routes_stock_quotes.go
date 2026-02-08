package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"sun-stockanalysis-api/internal/controllers"
)

func RegisterStockQuoteRoutes(api huma.API, controllers *controllers.Controllers, middleware func(huma.Context, func(huma.Context))) {
	protected := huma.NewGroup(api, "")
	protected.UseMiddleware(middleware)

	huma.Register(protected, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/stock-quotes",
		Summary: "List all stock quotes",
		Tags:    v1Tags(),
	}, controllers.StockQuoteController.ListAll)
}
