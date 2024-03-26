package middleware

import (
	"be-go-fiber-ecommerce/helper"
	"be-go-fiber-ecommerce/models"
	"be-go-fiber-ecommerce/models/web"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
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

func HandleError() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			if errorsValidation(c, err) {
				return nil
			}

			internalServerError(c, err)
		}

		return nil
	}
}

func errorsValidation(c *fiber.Ctx, err error) bool {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]web.ErrorResponse, len(ve))
		for _, fe := range ve {
			out = append(out, web.ErrorResponse{
				Field:   fe.Field(),
				Message: helper.MessageForTag(fe.Tag()),
			})
		}

		webResponse := web.WebResponse{
			Code:   fiber.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   out,
		}

		c.Status(fiber.StatusBadRequest).JSON(webResponse)
		return true
	}
	return false
}

func internalServerError(c *fiber.Ctx, err error) {
	webResponse := web.WebResponse{
		Code:   fiber.StatusInternalServerError,
		Status: "INTERNAL SERVER ERROR",
		Data:   err.Error(),
	}

	c.Status(fiber.StatusInternalServerError).JSON(webResponse)
}
