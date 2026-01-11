package handler

import "github.com/gofiber/fiber/v2"

type StockHandler interface {
	GetStock(c *fiber.Ctx) error
	CreateStock(c *fiber.Ctx) error
}