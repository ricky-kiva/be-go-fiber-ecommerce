package controller

import (
	"be-go-fiber-ecommerce/helper"
	"be-go-fiber-ecommerce/models/web"
	"be-go-fiber-ecommerce/service"

	"github.com/gofiber/fiber/v2"
)

type MidtransControllerImpl struct {
	MidtransService service.MidtransService
}

func NewMidtransControllerImpl(midtransService service.MidtransService) *MidtransControllerImpl {
	return &MidtransControllerImpl{
		MidtransService: midtransService,
	}
}

func (controller *MidtransControllerImpl) Create(c *fiber.Ctx) error {
	var request web.MidtransRequest
	if err := c.BodyParser(&request); err != nil {
		helper.PanicIfError(err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Binding body error",
		})
	}

	midtransResponse := controller.MidtransService.Create(c, request)
	webRespose := web.WebResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   midtransResponse,
	}

	return c.Status(fiber.StatusOK).JSON(webRespose)
}
