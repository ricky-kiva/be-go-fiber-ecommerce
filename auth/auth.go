package auth

import (
	"be-go-fiber-ecommerce/models"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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

	result := db.Create(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func LoginHandler(c *fiber.Ctx, db *gorm.DB) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	result := db.Where("email = ?", data["email"]).First(&user)
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
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
		Issuer:    "rickyslash.my.id",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	newClaims := UserClaims{
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
