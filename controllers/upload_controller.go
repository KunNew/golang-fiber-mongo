package controllers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		log.Println("Error in uploading Image : ", err)
	}

	splitFileName := strings.Split(file.Filename, ".")
	extension := splitFileName[len(splitFileName)-1]
	newFileName := fmt.Sprintf("%s.%s", time.Now().Format("2006-01-02-15-04-05"), extension)

	fileHeader, _ := file.Open()
	defer fileHeader.Close()

	// make folder in root dir
	folderUpload := filepath.Join(".", "uploads")
	if err := os.MkdirAll(folderUpload, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = c.SaveFile(file, "./uploads/"+newFileName)

	if err != nil {
		log.Println("Error in saving Image :", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	imageUrl := fmt.Sprintf("http://localhost:3000/uploads/%s", newFileName)

	data := map[string]interface{}{

		"imageName": newFileName,
		"imageUrl":  imageUrl,
		"header":    file.Header,
		"size":      file.Size,
	}

	return c.JSON(fiber.Map{"status": 201, "message": "Image uploaded successfully", "data": data})
}
