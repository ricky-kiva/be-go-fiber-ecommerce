package route

import (
	"be-go-fiber-ecommerce/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"Project": "Fish E-Commerce",
			"Dev":     "Rickyslash",
		})
	})

	v1 := app.Group("/v1")

	v1.Get("/products", func(c *fiber.Ctx) error {
		var products []models.Product

		if result := db.Find(&products); result.Error != nil {
			c.Status(fiber.StatusInternalServerError).JSON(result.Error)
		}

		return c.JSON(products)
	})

	v1.Get("/products/categories/:categoryId", func(c *fiber.Ctx) error {
		var products []models.Product

		categoryId := c.Params("categoryId")
		if result := db.Where("category_id = ?", categoryId).Find(&products); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(result.Error)
		}

		return c.JSON(products)
	})

	v1.Get("/products/:productId", func(c *fiber.Ctx) error {
		var product models.Product

		productId := c.Params("productId")
		if result := db.First(&product, productId); result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(result.Error)
		}

		return c.JSON(product)
	})

	v1.Get("/products/categories", func(c *fiber.Ctx) error {
		var categories []models.Category

		if result := db.Find(&categories); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(result.Error)
		}

		return c.JSON(categories)
	})
}
