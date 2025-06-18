package models

import "time"

type MessageResponse struct {
	SenderID   int       `json:"sender_id" bson:"sender_id"`
	ReceiverID int       `json:"receiver_id,omitempty" bson:"receiver_id,omitempty"`
	Message    string    `json:"message" bson:"message"` // This is the field we need to focus on
	Timestamp  time.Time `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}
