package controllers

import (
	"github.com/gofiber/fiber/v2"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

func FileUpload(ctx *fiber.Ctx) error {
	// Get uploaded file
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(err)
	}

	files := form.File["file"]
	if len(files) == 0 {
		return ctx.Status(http.StatusBadRequest).JSON("No file uploaded")
	}

	file := files[0]

	// Create physical file
	tempfile, err := os.CreateTemp("", "fiberocr"+"-")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(err)
	}
	defer func() {
		err = tempfile.Close()
		if err != nil {
			return
		}
		err = os.Remove(tempfile.Name())
		if err != nil {
			return
		}
	}()

	// Make uploaded physical
	upload, err := file.Open()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(err)
	}
	defer func(upload multipart.File) {
		err = upload.Close()
		if err != nil {

		}
	}(upload)

	if _, err := io.Copy(tempfile, upload); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(err)
	}

	client := gosseract.NewClient()
	defer func(client *gosseract.Client) {
		err = client.Close()
		if err != nil {
			return
		}
	}(client)

	err = client.SetImage(tempfile.Name())
	if err != nil {
		return err
	}
	client.Languages = []string{"eng"}

	languages := ctx.FormValue("languages")
	if languages != "" {
		client.Languages = strings.Split(languages, ",")
	}

	whitelist := ctx.FormValue("whitelist")
	if whitelist != "" {
		err = client.SetWhitelist(whitelist)
		if err != nil {
			return err
		}
	}

	var out string
	switch ctx.FormValue("format") {
	case "hocr":
		out, err = client.HOCRText()
	default:
		out, err = client.Text()
	}
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(err)
	}

	return ctx.Status(http.StatusOK).JSON(map[string]interface{}{
		"result":  strings.Trim(out, ctx.FormValue("trim")),
		"version": version,
	})
}
