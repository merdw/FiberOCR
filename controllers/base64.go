package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/otiai10/gosseract/v2"
)

type Body struct {
	Base64    string `json:"base64"`
	Trim      string `json:"trim"`
	Languages string `json:"languages"`
	Whitelist string `json:"whitelist"`
}

// Base64 ...
func Base64(c *fiber.Ctx) error {
	var body Body

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	tempfile, err := createTempFile()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer func() {
		err = cleanupTempFile(tempfile)
		if err != nil {
			// Handle error accordingly
		}
	}()

	if len(body.Base64) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "base64 string required"})
	}

	decoded, err := decodeBase64(body.Base64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = tempfile.Write(decoded)
	if err != nil {
		return err
	}

	client := gosseract.NewClient()
	defer func(client *gosseract.Client) {
		err = client.Close()
		if err != nil {
			// Handle error accordingly
		}
	}(client)

	err = configureTesseract(client, tempfile, &body)
	if err != nil {
		return err
	}

	result, err := processText(client, &body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result":  result,
		"version": version,
	})
}
