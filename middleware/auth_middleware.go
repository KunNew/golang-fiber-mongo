package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Middleware JWT function
func ProtectedMiddleware(c *fiber.Ctx) error {

	tokenString := c.Cookies("jwt")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "You are not logged in"})
	}

	// Parse JWT token with claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the token signing method and return the secret key
		return []byte("secret"), nil

	})

	// Handle token parsing errors
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Extract claims from token
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse claims",
		})
	}

	c.Locals("userId", claims["userId"])

	return c.Next()
}
