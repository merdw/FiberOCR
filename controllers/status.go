package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/otiai10/gosseract/v2"
)

const version = "0.5.0"

// Status ...
func Status(c *fiber.Ctx) error {
	langs, err := gosseract.GetAvailableLanguages()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	client := gosseract.NewClient()
	defer func(client *gosseract.Client) {
		err = client.Close()
		if err != nil {

		}
	}(client)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Hello!",
		"version": version,
		"tesseract": fiber.Map{
			"version":   client.Version(),
			"languages": langs,
		},
	})
}
