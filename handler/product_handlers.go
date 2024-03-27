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

func (h *Handler) UpdateProduct(c *fiber.Ctx) error {
	productId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product ID",
		})
	}

	updateData := make(map[string]interface{})
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error parsing request data",
		})
	}

	var product entity.Product
	if err := h.DB.First(&product, productId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Product not found",
		})
	}

	if result := h.DB.Model(&product).Updates(updateData); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not update product",
		})
	}

	return c.JSON(product)
}
