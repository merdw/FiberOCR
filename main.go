package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/merdw/fiberocr/controllers"
	"log"
	"os"
)

func main() {
	engine := html.New("./app/views", ".html")
	app := fiber.New(
		fiber.Config{
			Views:       engine,
			ViewsLayout: "",
			AppName:     "",
		})
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Static("/assets", "./app/assets")
	//app.Static("/", "./app/views/")

	app.Get("/status", controllers.Status)
	app.Post("/base64", controllers.Base64)
	app.Post("/captchabasic", controllers.CaptchaBasic)
	app.Post("/file", controllers.FileUpload)
	app.Get("/", controllers.Index)

	//app.Get("/", func(c *fiber.Ctx) error {
	//	return c.Render("index", fiber.Map{
	//		"AppName": "fiberocr",
	//	})
	//})
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.Render("test", fiber.Map{
			"AppName": "fiberocrdeneme",
		})
	})
	//app.Get("/test2", func(c *fiber.Ctx) error {
	//	return c.Render("test2", fiber.Map{
	//		"AppName": "fiberocrdeneme",
	//	})
	//})

	//logger = log.New(os.Stdout, fmt.Sprintf("[%s] ", "fiberocr"), 0)

	port := os.Getenv("PORT")

	//port = "8080"

	if port == "" {

		log.Fatalln("Requiresfd denv `PORT` is not specified.")
	}
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Dizin bilgisi alınamadı:", err)
		return
	}

	// Dizini ekrana yazdır
	fmt.Println("Çalıştığınız dizin:", currentDir)
	fmt.Printf("Server is running on :%s...\n")

	err = app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
