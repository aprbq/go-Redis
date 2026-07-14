package handlers

import "github.com/gofiber/fiber/v3"

type CatalogHandler interface {
	GetProducts(c fiber.Ctx) error
}
