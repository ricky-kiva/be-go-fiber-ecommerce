package service

import (
	"be-go-fiber-ecommerce/models/web"

	"github.com/gofiber/fiber/v2"
)

type MidtransService interface {
	CartTransaction(c *fiber.Ctx) (web.MidtransResponse, error)
}
