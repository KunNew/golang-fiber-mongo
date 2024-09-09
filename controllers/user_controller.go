package controllers

import (
	"context"
	"log"
	"strconv"
	"time"

	"example.com/fiberserver/config"
	"example.com/fiberserver/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUsers(c *fiber.Ctx) error {
	user_collection := config.MI.DB.Collection("users")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var users []models.User

	filter := bson.M{}
	projection := bson.D{{Key: "password", Value: 0}}

	findOptions := options.Find().SetProjection(projection)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limitVal, _ := strconv.Atoi(c.Query("limit", "10"))

	var limit int64 = int64(limitVal)

	total, _ := user_collection.CountDocuments(ctx, filter)

	findOptions.SetSkip((int64(page) - 1) * limit)

	findOptions.SetLimit(limit)

	cursor, err := user_collection.Find(ctx, filter, findOptions)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User Not found",
			"error":   err,
		})
	}
	for cursor.Next(ctx) {
		var user models.User

		err := cursor.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":  users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func GetCurrentUser(c *fiber.Ctx) error {

	userId := c.Locals("userId").(string)

	return c.JSON(fiber.Map{
		"userId": userId,
	})
}
