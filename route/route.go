package route

import (
	"be-go-fiber-ecommerce/auth"
	"be-go-fiber-ecommerce/handler"
	"be-go-fiber-ecommerce/middleware"
	"be-go-fiber-ecommerce/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	h := handler.New(db)

	v1 := app.Group("/v1")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"Project": "Fish E-Commerce",
			"Dev":     "Rickyslash",
		})
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		return auth.RegisterHandler(c, db)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return auth.LoginHandler(c, db)
	})

	app.Get("/cart", middleware.AuthValidator, middleware.AuthUserIdExtraction, h.GetCart)
	app.Post("/cart", middleware.AuthValidator, middleware.AuthUserIdExtraction, h.AddToCart)

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

	v1.Get("/products/info", func(c *fiber.Ctx) error {
		var product models.Product

		productId := c.Query("id")
		if productId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Product ID is required",
			})
		}

		if result := db.First(&product, productId); result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Product not found",
			})
		}

		return c.JSON(product)
	})

	v1.Get("/products/categories", middleware.AuthValidator, func(c *fiber.Ctx) error {
		var categories []models.Category

		if result := db.Find(&categories); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(result.Error)
		}

		return c.JSON(categories)
	})
}
