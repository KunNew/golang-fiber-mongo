package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"example.com/fiberserver/config"
	"example.com/fiberserver/models"
	"example.com/fiberserver/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	user_collection := config.MI.DB.Collection("users")

	// validate := validator.New()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user models.User

	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		// log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	fileUpload, errFile := utils.Upload(c)

	if errFile == nil {
		user.ImageUrl = fileUpload["imageUrl"].(string)
	}

	// fmt.Print("Good Morning" + errImg.Error())
	fmt.Printf("%+v\n", fileUpload["imageUrl"])

	// fmt.Print(uploadFile)

	// if validationErr := validate.Struct(&user); validationErr != nil {
	// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 	"error": validationErr.Error(),
	// })
	// }

	filter := bson.M{"email": user.Email}

	err := user_collection.FindOne(ctx, filter).Decode(&user)

	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
	}

	user.Password = string(hashedPassword)

	user.ID = primitive.NewObjectID()

	user.Password = string(hashedPassword)

	result, err := user_collection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "User failed to insert",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    result,
		"success": true,
	})
}

func Login(c *fiber.Ctx) error {

	user_collection := config.MI.DB.Collection("users")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var payload models.SignInInput

	if err := c.BodyParser(&payload); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	filter := bson.M{"email": payload.Email}
	// Find the user by credentials

	var user models.User

	err := user_collection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid email or Password", "error": err})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid email or Password"})
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"userId": user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"exp":    time.Now().Add(time.Minute * 1).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Set JWT token in cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(time.Minute * 1), // Expires in 24 hours
		HTTPOnly: true,
		Secure:   true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{"token": t})
}

func Logout(c *fiber.Ctx) error {

	// Clear JWT token by setting an empty value and expired time in the cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Minute), // Expired 1 hour ago
		HTTPOnly: true,
		Secure:   true,
	}
	c.Cookie(&cookie)

	// Return success response indicating logout was successful
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Logout successful",
	})
}
