package route

import (
	"be-go-fiber-ecommerce/controller"
	"be-go-fiber-ecommerce/handler"
	"be-go-fiber-ecommerce/middleware"
	"be-go-fiber-ecommerce/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	app.Use(middleware.HandleError())

	h := handler.New(db)

	v1 := app.Group("/v1")
	midtrans := app.Group("/midtrans")

	validate := validator.New()
	midtransService := service.NewMidTransServiceImpl(validate)
	midtransController := controller.NewMidtransControllerImpl(midtransService)

	app.Get("/", h.AboutProject)

	v1.Post("/register", h.UserRegister)
	v1.Post("/login", h.UserLogin)

	v1.Get("/cart", middleware.AuthValidator, middleware.AuthUserIdExtraction, h.GetCart)
	v1.Post("/cart", middleware.AuthValidator, middleware.AuthUserIdExtraction, h.AddToCart)
	v1.Delete("/cart/items/:productID", middleware.AuthValidator, middleware.AuthUserIdExtraction, h.DeleteCartItem)

	v1.Get("/products", h.GetAllProducts)
	v1.Get("/products/info", h.GetProductById)
	v1.Get("/products/categories/:categoryId", h.GetProductsByCategoryId)
	v1.Get("/products/categories", h.GetAllCategories)

	midtrans.Post("/create", midtransController.Create)
}
