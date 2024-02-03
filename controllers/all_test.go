package controllers

import (
	"github.com/gofiber/fiber/v2"
	"testing"
)

func testserver() *fiber.App {
	app := fiber.New()

	app.Post("/base64", Base64)
	app.Post("/file", FileUpload)
	app.Get("/status", Status)
	app.Get("/", Index)

	return app
}

func TestBase64(t *testing.T) {
	// TODO
}

func TestFileUpload(t *testing.T) {
	// TODO
}

func TestStatus(t *testing.T) {
	// TODO
}

func TestIndex(t *testing.T) {
	// TODO
}
