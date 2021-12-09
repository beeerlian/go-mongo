package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventActivity struct {
	UserId  primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Email   string             `json:"email,omitempty" bson:"email,omitempty"`
	Attende string             `json:"attende,omitempty" bson:"attende,omitempty"`
}
