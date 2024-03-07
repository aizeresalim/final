// image_controller.go

package controller

import (
	"github.com/gofiber/fiber/v2"
	"math/rand"
)

var letters = []rune("abcdefghijklmnopqrsuvwxyz")

func randLetter(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func UploadImage(c *fiber.Ctx) error {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Parsing form failed",
			"message": err.Error(),
		})
	}

	// Get files from the form
	files := form.File["image"]
	fileNames := make([]string, len(files))

	// Iterate over each file
	for i, file := range files {
		// Generate a unique file name
		fileName := randLetter(5) + "-" + file.Filename
		// Save the file
		if err := c.SaveFile(file, "./uploads/"+fileName); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Saving file failed",
				"message": err.Error(),
			})
		}
		// Store the file name
		fileNames[i] = fileName
	}

	// Return file names of uploaded files
	return c.JSON(fiber.Map{
		"fileNames": fileNames,
	})
}
