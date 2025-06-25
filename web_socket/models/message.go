package wsmodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	SenderID   int                `bson:"sender_id"`
	ReceiverID int                `bson:"receiver_id"`
	Message    string             `bson:"message"`
	Timestamp  time.Time          `bson:"timestamp"`
}

type IncomingMessage struct {
	ReceiverID int    `json:"receiver_id"`
	Message    string `json:"message"`
}

type ConnectionStatus struct {
	UserID  int    `json:"user_id"`
	Status  string `json:"status"` // "online", "offline"
	Message string `json:"message,omitempty"`
}
