package routes

import (
	"example.com/fiberserver/controllers"
	"example.com/fiberserver/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api", logger.New())

	// auth

	auth := api.Group("/auth")

	auth.Post("/login", controllers.Login)

	auth.Post("/signup", controllers.SignUp)

	auth.Post("/logout", controllers.Logout)

	// users
	users := api.Group("/users", middleware.ProtectedMiddleware)

	users.Get("/", controllers.GetUsers)

	// upload

	api.Post("/upload", controllers.Upload)

}
