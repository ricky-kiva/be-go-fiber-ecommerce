package route

import (
	"be-go-fiber-ecommerce/controller"
	"be-go-fiber-ecommerce/handler"
	"be-go-fiber-ecommerce/middleware"
	"be-go-fiber-ecommerce/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	h := handler.New(db)

	v1 := app.Group("/v1")

	midtransService := service.NewMidTransServiceImpl(db)
	midtransController := controller.NewMidtransController(midtransService)

	app.Get("/", h.AboutProject)

	v1.Post("/register", h.UserRegister)
	v1.Post("/login", h.UserLogin)

	v1.Get("/cart", middleware.AuthValidator, middleware.AuthUserIdExtraction, h.GetCart)
	v1.Post("/cart", middleware.AuthValidator, middleware.AuthUserIdExtraction, h.AddToCart)
	v1.Delete("/cart/items/:productID", middleware.AuthValidator, middleware.AuthUserIdExtraction, h.DeleteCartItem)
	v1.Get("/cart/pay", middleware.AuthValidator, middleware.AuthUserIdExtraction, midtransController.Pay)

	v1.Get("/products", h.GetAllProducts)
	v1.Get("/products/info", h.GetProductById)
	v1.Get("/products/categories/:categoryId", h.GetProductsByCategoryId)
	v1.Get("/products/categories", h.GetAllCategories)
}
