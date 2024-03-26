package handler

import (
	"be-go-fiber-ecommerce/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

func (h *Handler) AddToCart(c *fiber.Ctx) error {
	userId := c.Locals("userID").(uint)

	var cartInput struct {
		ProductID uint
		Quantity  int
	}

	if err := c.BodyParser(&cartInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error on binding",
		})
	}

	var product models.Product
	if err := h.DB.First(&product, cartInput.ProductID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Product not found",
		})
	}

	var cart models.Cart
	if err := h.DB.FirstOrCreate(&cart, models.Cart{UserID: userId}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not find or create cart",
		})
	}

	var cartItem models.CartItem
	err := h.DB.Where("cart_id = ? AND product_id = ?", cart.ID, cartInput.ProductID).First(&cartItem).Error
	if err == nil {
		cartItem.Quantity += cartInput.Quantity

		if cartItem.Quantity <= 0 {
			if err := h.DB.Delete(&cartItem).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": "Could not delete item from cart",
				})
			}
		} else {
			if err := h.DB.Save(&cartItem).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": "Could not update item in cart",
				})
			}
		}
	} else if err == gorm.ErrRecordNotFound {
		if cartInput.Quantity > 0 {
			cartItem = models.CartItem{
				CartID:    cart.ID,
				ProductID: cartInput.ProductID,
				Quantity:  cartInput.Quantity,
			}

			if err := h.DB.Create(&cartItem).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": "Could not add item to cart",
				})
			}
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Item quantity must be greater than 0",
			})
		}
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error checking cart item",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Cart updated successfully",
	})
}

func (h *Handler) GetCart(c *fiber.Ctx) error {
	userId := c.Locals("userID").(uint)

	var cart models.Cart
	if err := h.DB.Where("user_id = ?", userId).First(&cart).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cart was not initialized. Please add items",
		})
	}

	var cartItems []models.CartItem
	if err := h.DB.Where("cart_id = ?", cart.ID).Preload("Product").Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not load cart items",
		})
	}

	if len(cartItems) == 0 {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Your cart is empty",
		})
	}

	return c.JSON(cartItems)
}

func (h *Handler) DeleteCartItem(c *fiber.Ctx) error {
	userId := c.Locals("userID").(uint)
	productId, err := c.ParamsInt("productID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product ID",
		})
	}

	var cart models.Cart
	if err := h.DB.Where("user_id = ?", userId).First(&cart).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Cart was not initialized. Please add items",
		})
	}

	if err := h.DB.Where("cart_id = ? AND product_id = ?", cart.ID, productId).Delete(&models.CartItem{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not delete cart item",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
