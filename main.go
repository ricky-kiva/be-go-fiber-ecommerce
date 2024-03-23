package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"Pink Floyd": "The Wall",
			"Queen":      "Bohemian Rhapsody",
		})
	})

	app.Listen(":8080")
}
