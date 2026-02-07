package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"sun-stockanalysis-api/internal/controllers"
)

func RegisterRelationNewsRoutes(api huma.API, controllers *controllers.Controllers, middleware func(huma.Context, func(huma.Context))) {
	protected := huma.NewGroup(api, "")
	protected.UseMiddleware(middleware)

	huma.Register(protected, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/relation-news",
		Summary:       "Create relation news",
		Tags:          v1Tags(),
		DefaultStatus: http.StatusCreated,
	}, controllers.RelationNewsController.Create)
}
