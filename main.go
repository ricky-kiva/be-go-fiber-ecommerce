package main

import (
	"be-go-fiber-ecommerce/db"
	"be-go-fiber-ecommerce/initializer"
	"be-go-fiber-ecommerce/route"

	"github.com/gofiber/fiber/v2"
)

func init() {
	initializer.LoadEnv()
}

func main() {
	app := fiber.New()
	db := db.InitDb()

	route.Setup(app, db)

	app.Listen(":8080")
}
