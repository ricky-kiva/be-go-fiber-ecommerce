package auth

import (
	"be-go-fiber-ecommerce/models"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *fiber.Ctx, db *gorm.DB) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if !isEmailValid(data["email"]) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid email format",
		})
	}

	if len(data["password"]) < 3 || strings.Contains(data["password"], " ") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password must be at least 3 characters long & contain no spaces",
		})
	}

	hashedPassword, err := hashPassword(data["password"])
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not hash password",
		})
	}

	user := models.User{
		Email:    data["email"],
		Password: hashedPassword,
	}

	result := db.Create(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func isEmailValid(s string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(s)
}
