package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"sun-stockanalysis-api/internal/controllers"
)

func RegisterCompanyNewsRoutes(api huma.API, controllers *controllers.Controllers, middleware func(huma.Context, func(huma.Context))) {
	protected := huma.NewGroup(api, "")
	protected.UseMiddleware(middleware)

	huma.Register(protected, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/company-news",
		Summary: "List company news by relation symbols and date range",
		Tags:    v1Tags(),
	}, controllers.CompanyNewsController.ListBySymbolAndDate)
}
