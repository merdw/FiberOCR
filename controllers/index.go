package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// Index ...
func Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"AppName": "FiberOCR",
	})
}
