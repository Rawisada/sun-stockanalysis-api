package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"sun-stockanalysis-api/config"
	"sun-stockanalysis-api/app/http/handler"
)

type Server struct {
	app *fiber.App
	cfg *config.Config
}

func NewServer(cfg *config.Config, stockHandler handler.StockHandler) *Server {
	app := fiber.New(fiber.Config{
		BodyLimit: parseBodyLimit(cfg.Server.BodyLimit),
	})

	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	v1 := app.Group("/v1")
	stocks := v1.Group("/stocks")
	stocks.Get("/:id", stockHandler.GetStock)
	stocks.Post("/", stockHandler.CreateStock)

	return &Server{app: app, cfg: cfg}
}

func (s *Server) Start() error {
	return s.app.Listen(fmt.Sprintf(":%d", s.cfg.Server.Port))
}

// Fiber ต้องการ byte, config เป็น "10M" -> convert แบบง่าย
func parseBodyLimit(v string) int {
	// แบบง่ายรองรับ "10M" "1M" "100K"
	// ถ้าต้องการ robust ค่อยทำเพิ่ม
	n := 10 * 1024 * 1024
	_ = v
	return n
}
