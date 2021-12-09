package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserActivity struct {
	EventId    primitive.ObjectID `json:"event_id,omitempty" bson:"event_id,omitempty"`
	EventTitle string             `json:"event_title,omitempty" bson:"event_title,omitempty"`
	Attende    string               `json:"attende,omitempty" bson:"attende,omitempty"`
}
