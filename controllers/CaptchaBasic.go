package controllers

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/otiai10/gosseract/v2"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

const (
	th1 float64 = 140

	th2 float64 = 140.0
	sig float64 = 1.5
)

func savePreprocessedImage(img image.Image, filename string) {
	out, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating preprocessed image:", err)
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)
	err = png.Encode(out, img)
	if err != nil {
		return
	}
	fle, err := os.Stat(out.Name())
	fmt.Println("filenin yol:", fle.Name())
	abs_fname, err := filepath.Abs(fle.Name())
	fmt.Println("gercekyok:", abs_fname)
}

func Scale(file string) error {

	img, err := imaging.Open(file)
	if err != nil {
		fmt.Println("Error opening image:", err)
		return err
	}
	// Convert the image to grayscale
	grayImg := imaging.Grayscale(img)

	// Apply the first threshold
	firstThreshold := imaging.AdjustContrast(grayImg, th1)

	// Apply Gaussian blur
	blurred := imaging.Blur(firstThreshold, sig)

	// Apply the second threshold
	finalThreshold := imaging.AdjustContrast(blurred, th2)

	// Edge enhance and sharpen the image
	//final := imaging.Overlay(finalThreshold, img, image.Pt(0, 0), 1.0)

	final := finalThreshold
	//filed, err := os.Create("output.png")
	//if err != nil {
	//	panic(err)
	//}
	//defer filed.Close()

	// PNG formatında kaydet
	//err = png.Encode(filed, final)
	//if err != nil {
	//	panic(err)
	//}

	// Save the preprocessed image (optional)
	//savePreprocessedImage(final, "preprocessed.png")
	savePreprocessedImage(final, file)
	return nil

}

func CaptchaBasic(c *fiber.Ctx) error {

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

	err = Scale(tempfile.Name())
	if err != nil {
		fmt.Println("Scale hata", err)
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

	if body.Whitelist != "" {
		err = client.SetWhitelist(body.Whitelist)
		if err != nil {
			return err
		}
	}

	result, err := processText(client, &body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Println("result:", result)
	x := strings.ReplaceAll(result, " ", "")
	fmt.Println("x result:", result)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result":  x,
		"version": version,
	})
}
