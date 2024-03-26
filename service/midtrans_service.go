package service

import (
	"be-go-fiber-ecommerce/model"

	"github.com/gofiber/fiber/v2"
)

type MidtransService interface {
	CartTransaction(c *fiber.Ctx) (model.MidtransResponse, error)
}
