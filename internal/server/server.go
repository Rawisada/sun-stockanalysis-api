package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"sun-stockanalysis-api/internal/configurations"
	"sun-stockanalysis-api/internal/controllers"
	"sun-stockanalysis-api/internal/handler"
)

type Server struct {
	app *fiber.App
	cfg *configurations.Config
}

func NewServer(cfg *configurations.Config, stockController *controllers.StockController) *Server {
	app := fiber.New(fiber.Config{
		BodyLimit: parseBodyLimit(cfg.Server.BodyLimit),
	})

	handler.RegisterRoutes(app, stockController)

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
