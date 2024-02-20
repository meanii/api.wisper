package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/meanii/api.wisper/clients"
)

type User struct {
	Id       primitive.ObjectID   `json:"_id,omitempty"      bson:"_id,omitempty"`
	Username string               `json:"username,omitempty"                      validate:"required"`
	Password string               `json:"password,omitempty"                      validate:"required"`
	Friends  []primitive.ObjectID `json:"friends,omitempty"`
}

func (u *User) Collection() *mongo.Collection {
	return clients.GetDatabase().Collection("users")
}
