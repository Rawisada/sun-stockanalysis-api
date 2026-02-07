package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"sun-stockanalysis-api/internal/controllers"
)

func RegisterHealthRoutes(api huma.API, controllers *controllers.Controllers) {
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
