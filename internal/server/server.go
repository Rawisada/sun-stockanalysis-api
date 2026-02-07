package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"

	"sun-stockanalysis-api/internal/configurations"
	"sun-stockanalysis-api/internal/controllers"
	"sun-stockanalysis-api/internal/handler"
	"sun-stockanalysis-api/pkg/apierror"
	"sun-stockanalysis-api/pkg/logger"
)

type Server struct {
	cfg *configurations.Config
	app *fiber.App
	log *logger.Logger
}

func NewServer(cfg *configurations.Config, controllers *controllers.Controllers, log *logger.Logger) *Server {
	app := fiber.New(fiber.Config{
		AppName: "sun-stockanalysis-api",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status, resp := apierror.ToResponse(err)
			return c.Status(status).JSON(resp)
		},
	})

	contextPath := normalizeContextPath(cfg.Server.ContextPath)
	apiGroup := app.Group(contextPath)

	apiGroup.Use(func(c *fiber.Ctx) error {
		path := c.Path()
		if path == contextPath+"/docs" ||
			path == contextPath+"/openapi.json" ||
			path == contextPath+"/openapi.yaml" ||
			path == contextPath+"/schemas" {
			return c.Next()
		}

		correlationID := c.Get("X-Correlation-Id")
		if correlationID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "X-Correlation-Id header is required")
		}

		c.Set("X-Correlation-Id", correlationID)
		ctx := context.WithValue(c.UserContext(), logger.CorrelationIDKey, correlationID)
		c.SetUserContext(ctx)
		c.Context().SetUserValue(logger.CorrelationIDKey, correlationID)

		start := time.Now()
		err := c.Next()
		if log != nil {
			log.With("correlation_id", correlationID).Infof(
				"%s %s %d %s",
				c.Method(),
				c.OriginalURL(),
				c.Response().StatusCode(),
				time.Since(start),
			)
		}
		return err
	})

	huma.NewError = apierror.NewHumaError
	huma.NewErrorWithContext = func(_ huma.Context, status int, msg string, errs ...error) huma.StatusError {
		return apierror.NewHumaError(status, msg, errs...)
	}

	apiConfig := huma.DefaultConfig("sun-stockanalysis-api", "1.0.0")
	apiConfig.Servers = []*huma.Server{{URL: contextPath}}
	if apiConfig.Components.SecuritySchemes == nil {
		apiConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{}
	}
	apiConfig.Components.SecuritySchemes["CorrelationId"] = &huma.SecurityScheme{
		Type: "apiKey",
		In:   "header",
		Name: "X-Correlation-Id",
	}
	apiConfig.Components.SecuritySchemes["BearerAuth"] = &huma.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	}

	humaAPI := humafiber.NewWithGroup(app, apiGroup, apiConfig)
	handler.RegisterRoutes(humaAPI, controllers, cfg.State.Secret, cfg.State.Issuer)
	addCorrelationIDToOpenAPI(humaAPI)
	addBearerAuthToOpenAPI(humaAPI)

	app.Get(contextPath+"/docs", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(swaggerUIHTML(contextPath))
	})
	app.Get(contextPath+"/swagger-ui.html", func(c *fiber.Ctx) error {
		return c.Redirect(contextPath+"/swagger-ui/index.html", fiber.StatusFound)
	})
	app.Get(contextPath+"/swagger-ui/index.html", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(swaggerUIHTML(contextPath))
	})

	return &Server{
		cfg: cfg,
		app: app,
		log: log,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)
	return s.app.Listen(addr)
}

func (s *Server) Stop() error {
	timeout := s.cfg.Server.TimeOut
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.app.ShutdownWithContext(ctx)
}

func addCorrelationIDToOpenAPI(api huma.API) {
	openapi := api.OpenAPI()
	correlationParam := &huma.Param{
		Name:        "X-Correlation-Id",
		In:          "header",
		Description: "Request correlation ID for tracking",
		Required:    true,
		Schema: &huma.Schema{
			Type: "string",
		},
	}
	for path := range openapi.Paths {
		pathItem := openapi.Paths[path]
		operations := []*huma.Operation{
			pathItem.Get,
			pathItem.Post,
			pathItem.Put,
			pathItem.Patch,
			pathItem.Delete,
		}
		for _, op := range operations {
			if op != nil {
				op.Parameters = append(op.Parameters, correlationParam)
			}
		}
	}
}

func addBearerAuthToOpenAPI(api huma.API) {
	openapi := api.OpenAPI()
	security := map[string][]string{"BearerAuth": {}}
	for path := range openapi.Paths {
		pathItem := openapi.Paths[path]
		operations := []*huma.Operation{
			pathItem.Get,
			pathItem.Post,
			pathItem.Put,
			pathItem.Patch,
			pathItem.Delete,
		}
		for _, op := range operations {
			if op != nil {
				op.Security = append(op.Security, security)
			}
		}
	}
}

func normalizeContextPath(path string) string {
	p := strings.TrimSpace(path)
	if p == "" {
		return ""
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return strings.TrimSuffix(p, "/")
}

func swaggerUIHTML(contextPath string) string {
	specURL := contextPath + "/openapi.json"
	return `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
    <style>
      :root {
        --springdoc-blue: #1b78c8;
        --springdoc-blue-dark: #165f9f;
        --springdoc-bg: #f7f9fc;
        --springdoc-text: #1f2a44;
      }
      body {
        margin: 0;
        background: var(--springdoc-bg);
        color: var(--springdoc-text);
        font-family: "Segoe UI", "Helvetica Neue", Arial, sans-serif;
      }
      .swagger-ui .topbar {
        background-color: var(--springdoc-blue);
        border-bottom: 1px solid rgba(0,0,0,0.08);
        padding: 8px 0;
      }
      .swagger-ui .topbar .download-url-wrapper .download-url-button {
        background: var(--springdoc-blue-dark);
        border: none;
      }
      .swagger-ui .topbar .download-url-wrapper input[type="text"] {
        border: 1px solid rgba(0,0,0,0.1);
      }
      .swagger-ui .info {
        margin: 30px 0 10px;
      }
      .swagger-ui .info .title {
        color: var(--springdoc-text);
      }
      .swagger-ui .btn.authorize {
        background: var(--springdoc-blue);
        border-color: var(--springdoc-blue);
      }
      .swagger-ui .opblock.opblock-get {
        border-color: rgba(27,120,200,0.3);
        background: rgba(27,120,200,0.05);
      }
      .swagger-ui .opblock.opblock-post {
        border-color: rgba(73,204,144,0.3);
        background: rgba(73,204,144,0.05);
      }
      .swagger-ui .opblock.opblock-put {
        border-color: rgba(252,161,48,0.3);
        background: rgba(252,161,48,0.05);
      }
      .swagger-ui .opblock.opblock-delete {
        border-color: rgba(249,62,62,0.3);
        background: rgba(249,62,62,0.05);
      }
    </style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = function () {
        window.ui = SwaggerUIBundle({
          url: "` + specURL + `",
          dom_id: "#swagger-ui",
        });
      };
    </script>
  </body>
</html>
`
}
