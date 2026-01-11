package handler

import (
	stock "sun-stockanalysis-api/internal/usecase/stock"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type stockHandlerImpl struct {
	stockService stock.StockService
}

func NewStockHandlerImpl(stockService stock.StockService) StockHandler {
	return &stockHandlerImpl{stockService: stockService}
}

func (h *stockHandlerImpl) GetStock(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid uuid"})
	}

	s, err := h.stockService.GetStock(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(s)
}

type createStockReq struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Sector string `json:"sector"`
	Price  int    `json:"price"`
}

func (h *stockHandlerImpl) CreateStock(c *fiber.Ctx) error {
	var req createStockReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}
	if req.Symbol == "" || req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "symbol and name required"})
	}

	if err := h.stockService.CreateStock(req.Symbol, req.Name, req.Sector, req.Price); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"status": "created"})
}
