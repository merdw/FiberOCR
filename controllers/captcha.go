package controllers

import (
	"encoding/base64"
	"github.com/otiai10/gosseract/v2"
	"os"
	"regexp"
	"strings"
)

//func Solve()  {
//	client := gosseract.NewClient()
//	defer func(client *gosseract.Client) {
//		err := client.Close()
//		if err != nil {
//			return
//		}
//	}(client)
//
//	client.Languages = []string{"eng"}
//	if body.Languages != "" {
//		client.Languages = strings.Split(body.Languages, ",")
//	}
//	err := client.SetImage(tempfile.Name())
//	if err != nil {
//		return err
//	}
//	if body.Whitelist != "" {
//		err = client.SetWhitelist(body.Whitelist)
//		if err != nil {
//			return err
//		}
//	}
//}

func decodeBase64(base64String string) ([]byte, error) {
	base64String = regexp.MustCompile("data:image\\/png;base64,").ReplaceAllString(base64String, "")
	base64String = strings.TrimSpace(base64String)
	decoded, err := base64.StdEncoding.DecodeString(base64String)
	return decoded, err
}

func createTempFile() (*os.File, error) {
	return os.CreateTemp("", "fiberocr"+"-")
}

func cleanupTempFile(tempfile *os.File) error {
	err := tempfile.Close()
	if err != nil {
		return err
	}
	return os.Remove(tempfile.Name())
}

func configureTesseract(client *gosseract.Client, tempfile *os.File, body *Body) error {
	client.Languages = []string{"eng"}
	if body.Languages != "" {
		client.Languages = strings.Split(body.Languages, ",")
	}

	err := client.SetImage(tempfile.Name())
	if err != nil {
		return err
	}

	if body.Whitelist != "" {
		err = client.SetWhitelist(body.Whitelist)
		if err != nil {
			return err
		}
	}

	return nil
}

func processText(client *gosseract.Client, body *Body) (string, error) {
	text, err := client.Text()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(text), nil
}
