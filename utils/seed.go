package utils

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"example.com/fiberserver/config"
	"github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Seed() error {
	user_collection := config.MI.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, _ := user_collection.CountDocuments(ctx, bson.M{})
	if count == 0 {
		pass := "superpwd@" + strconv.Itoa(time.Now().Year())

		var users []interface{}

		for i := 0; i < 100; i++ {

			name := faker.FirstName() + faker.LastName()

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
			if err != nil {
				log.Fatal(err)
			}

			users = append(users, map[string]string{
				"name":     name,
				"email":    name + "@e-corp.com",
				"password": string(hashedPassword),
			})
		}

		_, err := user_collection.InsertMany(ctx, users)

		if err != nil {
			return fmt.Errorf("failed to insert documents: %w", err)
		}

		fmt.Println("Data seeded successfully")
	}

	return nil
}
