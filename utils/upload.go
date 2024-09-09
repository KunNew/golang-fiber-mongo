package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Upload(c *fiber.Ctx) (fiber.Map, error) {
	file, err := c.FormFile("image")
	if err != nil {

		log.Println("Error in uploading Image : ", err)
		// return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })
		return nil, err
	}

	splitFileName := strings.Split(file.Filename, ".")
	extension := splitFileName[len(splitFileName)-1]
	newFileName := fmt.Sprintf("%s.%s", time.Now().Format("2006-01-02-15-04-05"), extension)

	fileHeader, _ := file.Open()
	defer fileHeader.Close()

	// make folder in root dir
	folderUpload := filepath.Join(".", "uploads")
	if err := os.MkdirAll(folderUpload, os.ModePerm); err != nil {
		// return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"message": err.Error(),
		// })
		return nil, err
	}

	err = c.SaveFile(file, "./uploads/"+newFileName)

	if err != nil {
		log.Println("Error in saving Image :", err)
		// return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"message": err.Error(),
		// })
		return nil, err
	}

	imageUrl := fmt.Sprintf("http://localhost:3000/uploads/%s", newFileName)

	data := fiber.Map{
		"imageName": newFileName,
		"imageUrl":  imageUrl,
		"header":    file.Header,
		"size":      file.Size,
	}
	return data, nil
}
