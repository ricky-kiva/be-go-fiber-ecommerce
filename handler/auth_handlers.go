package handler

import (
	"be-go-fiber-ecommerce/models"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (h *Handler) UserRegister(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if !isEmailValid(data["email"]) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid email format",
		})
	}

	if len(data["password"]) < 3 || strings.Contains(data["password"], " ") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Password must be at least 3 characters long & contain no spaces",
		})
	}

	hashedPassword, err := hashPassword(data["password"])
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not hash password",
		})
	}

	user := models.User{
		Email:    data["email"],
		Password: hashedPassword,
	}

	result := h.DB.Create(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *Handler) UserLogin(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	result := h.DB.Where("email = ?", data["email"]).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Email is not registered",
		})
	}

	if !checkPasswordHash(data["password"], user.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Incorrect password",
		})
	}

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		Issuer:    "rickyslash.my.id",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	newClaims := models.UserClaims{
		RegisteredClaims: claims,
		UserID:           user.ID,
	}

	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	secret := os.Getenv("JWT_SECRET")

	token, err := sign.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create token",
		})
	}

	return c.JSON(fiber.Map{"token": token})
}
