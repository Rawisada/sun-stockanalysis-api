package server

import (
	"fmt"
	"log"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"go.opentelemetry.io/otel"

	common "sun-stockanalysis-api/internal/common"
	"sun-stockanalysis-api/internal/controllers"
	"sun-stockanalysis-api/internal/handler"
)

type Server struct {
	App    *fiber.App
	Api    huma.API
	Config *ServerConfig
}

type ServerTitle string
type ServerVersion string

type ServerConfig struct {
	Title                   ServerTitle
	Version                 ServerVersion
	Port                    string
	EnableTrustedProxyCheck bool
	TrustedProxies          []string
	MaxPayloadSizeKB        int
	TimeoutSeconds          int
	AuthSecret              string
	AuthIssuer              string
}

const (
	apiKeySchemeName = "BearerAuth"
)

func newHumaConfig(title string, version string) huma.Config {
	huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
		return common.NewErrorResponse(status, msg, errs...)
	}

	cfg := huma.DefaultConfig(title, version)
	cfg.DocsPath = ""
	cfg.CreateHooks = nil
	cfg.Transformers = nil
	cfg.OnAddOperation = nil
	cfg.OpenAPI.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		apiKeySchemeName: {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
	cfg.OpenAPI.Security = []map[string][]string{
		{apiKeySchemeName: {}},
	}
	return cfg
}

func NewServer(serverConfig ServerConfig, controllers *controllers.Controllers) *Server {
	humaConfig := newHumaConfig(string(serverConfig.Title), string(serverConfig.Version))
	app := fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		EnableTrustedProxyCheck: serverConfig.EnableTrustedProxyCheck,
		TrustedProxies:          serverConfig.TrustedProxies,
		ProxyHeader:             fiber.HeaderXForwardedFor,
		BodyLimit:               serverConfig.MaxPayloadSizeKB * 1024,
	})

	api := humafiber.New(app, humaConfig)
	handler.RegisterRoutes(api, controllers, serverConfig.AuthSecret, serverConfig.AuthIssuer)

	app.Use(openTelemetryMiddleware())
	app.Use(logger.New())
	app.Use(appLoggerMiddleware())

	if serverConfig.TimeoutSeconds > 0 {
		app.Use(timeout.NewWithContext(func(c *fiber.Ctx) error {
			return c.Next()
		}, time.Duration(serverConfig.TimeoutSeconds)*time.Second))
	}

	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		log.Printf("server starting on :%s", listenData.Port)
		return nil
	})

	app.Get("/docs", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(scalarHTML)
	})

	return &Server{
		App:    app,
		Api:    api,
		Config: &serverConfig,
	}
}

func (s *Server) Start() error {
	return s.App.Listen(":" + s.Config.Port)
}

func appLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			log.Printf("request error: %s %s: %v", c.Method(), c.Path(), err)
		}
		return err
	}
}

func openTelemetryMiddleware() fiber.Handler {
	tracer := otel.Tracer("sun-stockanalysis-api")
	return func(c *fiber.Ctx) error {
		ctx, span := tracer.Start(c.UserContext(), fmt.Sprintf("%s %s", c.Method(), c.Path()))
		defer span.End()
		c.SetUserContext(ctx)
		return c.Next()
	}
}

const scalarHTML = `<!DOCTYPE html>
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

// const scalarHTML = `<!doctype html>
// <html>
//   <head>
//     <title>Scalar API Reference</title>
//     <meta charset="utf-8" />
//     <meta
//       name="viewport"
//       content="width=device-width, initial-scale=1" />
//   </head>

//   <body>
//     <div id="app"></div>

//     <!-- Load the Script -->
//     <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>

//     <!-- Initialize the Scalar API Reference -->
//     <script>
//       Scalar.createApiReference('#app', {
//         // The URL of the OpenAPI/Swagger document
//         url: '/openapi.json',
//       })
//     </script>
//   </body>
// </html>	`
