package service

import (
	"be-go-fiber-ecommerce/models/web"

	"github.com/gofiber/fiber/v2"
)

type MidtransService interface {
	Create(c *fiber.Ctx, request web.MidtransRequest) web.MidtransResponse
}
