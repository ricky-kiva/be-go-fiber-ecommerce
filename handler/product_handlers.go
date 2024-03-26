package handler

import (
	"be-go-fiber-ecommerce/entity"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetAllProducts(c *fiber.Ctx) error {
	var products []entity.Product

	if result := h.DB.Find(&products); result.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(result.Error)
	}

	return c.JSON(products)
}

func (h *Handler) GetProductById(c *fiber.Ctx) error {
	var product entity.Product

	productId := c.Query("id")
	if productId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Product ID is required",
		})
	}

	if result := h.DB.First(&product, productId); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	return c.JSON(product)
}

func (h *Handler) GetProductsByCategoryId(c *fiber.Ctx) error {
	var products []entity.Product

	categoryId := c.Params("categoryId")
	if result := h.DB.Where("category_id = ?", categoryId).Find(&products); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(result.Error)
	}

	return c.JSON(products)
}

func (h *Handler) GetAllCategories(c *fiber.Ctx) error {
	var categories []entity.Category

	if result := h.DB.Find(&categories); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(result.Error)
	}

	return c.JSON(categories)
}
