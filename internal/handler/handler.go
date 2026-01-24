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

func RegisterRoutes(rootApi huma.API, controllers *controllers.Controllers, authSecret, authIssuer string) {
	rootApi.UseMiddleware(requestIDMiddleware)

	registerHealthHandlers(rootApi, controllers)
	v1Api := huma.NewGroup(rootApi, apiBasePath)

	registerV1(v1Api, controllers, authSecret, authIssuer)
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

func registerV1(api huma.API, controllers *controllers.Controllers, authSecret, authIssuer string) {
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

	protected := huma.NewGroup(api, "")
	protected.UseMiddleware(authMiddleware(authSecret, authIssuer))

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

func requestIDMiddleware(ctx huma.Context, next func(huma.Context)) {
	if ctx.Header("X-Request-Id") == "" {
		ctx.SetHeader("X-Request-Id", uuid.NewString())
	}
	next(ctx)
}
