package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"sun-stockanalysis-api/internal/controllers"
)

func RegisterAuthRoutes(api huma.API, controllers *controllers.Controllers) {
	huma.Register(api, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/register",
		Summary: "Register",
		Tags:    v1Tags(),
	}, controllers.AuthController.Register)

	huma.Register(api, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/login",
		Summary: "Login",
		Tags:    v1Tags(),
	}, controllers.AuthController.Login)

	huma.Register(api, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/refresh",
		Summary: "Refresh access token",
		Tags:    v1Tags(),
	}, controllers.AuthController.Refresh)
}
