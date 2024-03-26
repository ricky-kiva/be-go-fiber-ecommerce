package controller

import (
	"be-go-fiber-ecommerce/service"

	"github.com/gofiber/fiber/v2"
)

type MidtransController struct {
	MidtransService service.MidtransService
}

func NewMidtransController(midtransService service.MidtransService) *MidtransController {
	return &MidtransController{
		MidtransService: midtransService,
	}
}

func (controller *MidtransController) Pay(c *fiber.Ctx) error {
	midtransResponse, err := controller.MidtransService.CartTransaction(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   midtransResponse,
	})
}
