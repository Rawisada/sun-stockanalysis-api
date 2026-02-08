package handler

import (
	"github.com/danielgtaylor/huma/v2"

	"sun-stockanalysis-api/internal/controllers"
	"sun-stockanalysis-api/internal/routes"
)

const apiBasePath = "/v1"

func RegisterRoutes(rootApi huma.API, controllers *controllers.Controllers, authSecret, authIssuer string) {
	// rootApi.UseMiddleware(requestIDMiddleware)

	routes.RegisterHealthRoutes(rootApi, controllers)
	v1Api := huma.NewGroup(rootApi, apiBasePath)

	routes.RegisterAuthRoutes(v1Api, controllers)
	routes.RegisterStockRoutes(v1Api, controllers, authMiddleware(authSecret, authIssuer))
	routes.RegisterStockQuoteRoutes(v1Api, controllers, authMiddleware(authSecret, authIssuer))
	routes.RegisterStockDailyRoutes(v1Api, controllers, authMiddleware(authSecret, authIssuer))
	routes.RegisterCompanyNewsRoutes(v1Api, controllers, authMiddleware(authSecret, authIssuer))
	routes.RegisterRelationNewsRoutes(v1Api, controllers, authMiddleware(authSecret, authIssuer))
}

// func requestIDMiddleware(ctx huma.Context, next func(huma.Context)) {
// 	if ctx.Header("X-Request-Id") == "" {
// 		ctx.SetHeader("X-Request-Id", uuid.NewString())
// 	}
// 	next(ctx)
// }
