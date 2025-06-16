package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	SenderID   int                `json:"sender_id" bson:"sender_id"`
	ReceiverID int                `json:"receiver_id,omitempty" bson:"receiver_id,omitempty"`
	Message    string             `json:"message" bson:"message"` // This is the field we need to focus on
	Timestamp  time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

type UpdateMessageRequest struct {
	ID      string `json:"id" binding:"required"`
	Message string `json:"message" binding:"required"`
}
type SendMessageRequest struct {
	Content    string `json:"content" binding:"required"`
	ReceiverID int    `json:"receiver_id,omitempty"`
}
