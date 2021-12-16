package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Link        string             `json:"link,omitempty" bson:"link,omitempty"`
	Time        string             `json:"time,omitempty" bson:"time,omitempty"`
	Lecturer    string             `json:"lecturer,omitempty" bson:"lecturer,omitempty"`
	Participant []EventActivity    `json:"participant,omitempty" bson:"participant,omitempty"`
}
