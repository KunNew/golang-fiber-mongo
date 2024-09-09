package main

import (
	"log"
	"os"

	"example.com/fiberserver/config"
	"example.com/fiberserver/routes"
	"example.com/fiberserver/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	app := fiber.New()

	// Serve static files from the uploads directory
	app.Static("/uploads", "./uploads")

	app.Use(cors.New())

	//run database
	config.ConnectDB()

	//seed database
	utils.Seed()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// init routes
	routes.SetupRoutes(app)

	// custom 404 handler at router tail
	// app.Use(func(ctx *fiber.Ctx) error {
	// 	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
	// 		"message": "Route Not Found",
	// 	})
	// })

	app.Listen(":3000")
}
