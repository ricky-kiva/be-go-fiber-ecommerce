package main

import (
	"be-go-fiber-ecommerce/db"
	// "be-go-fiber-ecommerce/initializer"
	"be-go-fiber-ecommerce/route"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// func init() {
// 	initializer.LoadEnv()
// }

func main() {
	app := fiber.New()
	db := db.InitDb()

	route.Setup(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := app.Listen(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
