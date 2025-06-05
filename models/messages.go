package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SenderID   int                `json:"sender_id,omitempty" bson:"sender_id,omitempty"`
	ReceiverID int                `json:"receiver_id,omitempty" bson:"receiver_id,omitempty"`
	Message    string             `json:"message" bson:"message"`
	Timestamp  time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

