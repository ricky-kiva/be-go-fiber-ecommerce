package middleware

import (
	"be-go-fiber-ecommerce/models"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthValidator(c *fiber.Ctx) error {
	authValue := c.Get("Authorization")
	const bearerPrefix = "Bearer "

	if !strings.HasPrefix(authValue, bearerPrefix) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing or malformed JWT",
		})
	}

	tokenString := strings.TrimPrefix(authValue, bearerPrefix)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, valid := t.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("invalid token: %s", t.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid login session",
		})
	}

	return c.Next()
}

func AuthUserIdExtraction(c *fiber.Ctx) error {
	authValue := c.Get("Authorization")
	const bearerPrefix = "Bearer "

	tokenString := strings.TrimPrefix(authValue, bearerPrefix)

	token, err := jwt.ParseWithClaims(tokenString, &models.UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("error parsing token: %v", err),
		})
	}

	if claims, ok := token.Claims.(*models.UserClaims); ok && token.Valid {
		c.Locals("userID", claims.UserID)
		return c.Next()
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid login session",
		})
	}
}
