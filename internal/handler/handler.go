package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"sun-stockanalysis-api/internal/controllers"
)

const (
	apiBasePath = "/v1"
	stocksTag   = "stocks"
)

func v1Tags(additionalTags ...string) []string {
	tags := []string{"PosClient", "v1"}
	tags = append(tags, additionalTags...)
	return tags
}



type healthResponse struct {
	Body struct {
		Status string `json:"status"`
	}
}

func RegisterRoutes(app *fiber.App, stockController *controllers.StockController) huma.API {
	app.Use(recover.New())
	app.Use(logger.New())

	cfg := huma.DefaultConfig("sun-stockanalysis-api", "1.0.0")
	cfg.DocsPath = ""
	api := humafiber.New(app, cfg)

	app.Get("/docs", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(swaggerHTML)
	})

	huma.Register(api, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/health",
		Summary: "Health check",
		Tags:    []string{"system"},
	}, func(ctx context.Context, input *struct{}) (*healthResponse, error) {
		_ = ctx
		_ = input

		var res healthResponse
		res.Body.Status = "ok"
		return &res, nil
	})

	huma.Register(api, huma.Operation{
		Method:  http.MethodGet,
		Path:    apiBasePath + "/stocks/{id}",
		Summary: "Get stock by ID",
		Tags:    []string{stocksTag},
	}, stockController.GetStock)

	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          apiBasePath + "/stocks",
		Summary:       "Create stock",
		Tags:          []string{stocksTag},
		DefaultStatus: http.StatusCreated,
	}, stockController.CreateStock)

	return api
}

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = function () {
        window.ui = SwaggerUIBundle({
          url: "/openapi.json",
          dom_id: "#swagger-ui",
        });
      };
    </script>
  </body>
</html>
`
