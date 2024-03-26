package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Handler {
	return Handler{DB: db}
}

func (h *Handler) AboutProject(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"Project": "Fish E-Commerce",
		"Dev":     "Rickyslash",
	})
}
