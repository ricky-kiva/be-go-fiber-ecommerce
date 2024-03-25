package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthValidator(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	const bearerPrefix = "Bearer "

	if !strings.HasPrefix(auth, bearerPrefix) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing or malformed JWT",
		})
	}

	tokenString := strings.TrimPrefix(auth, bearerPrefix)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, valid := t.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("invalid token: %s", t.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid or expired login session",
		})
	}

	return c.Next()
}
