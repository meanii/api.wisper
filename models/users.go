package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID   `json:"id,omitempty"`
	Username string               `json:"username,omitempty" validate:"required"`
	Password string               `json:"password,omitempty" validate:"required"`
	Friends  []primitive.ObjectID `json:"friends,omitempty"`
}
