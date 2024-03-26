package controller

import "github.com/gofiber/fiber/v2"

type MidtransController interface {
	Create(c *fiber.Ctx)
}
