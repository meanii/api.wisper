package models

import (
	"github.com/meanii/api.wisper/clients"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string               `json:"username,omitempty" validate:"required"`
	Password string               `json:"password,omitempty" validate:"required"`
	Friends  []primitive.ObjectID `json:"friends,omitempty"`
}

var UserModel = clients.GetDatabase().Collection("users")
