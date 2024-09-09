package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Email    string             `json:"email" gorm:"unique" bson:"email,omitempty" validate:"required"`
	Password string             `json:"password,omitempty" validate:"required"`
	ImageUrl string             `json:"imageUrl,omitempty" bson:"imageUrl,omitempty"`
}

type SignInInput struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}
