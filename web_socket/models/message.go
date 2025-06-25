package wsmodels

import "time"

type MessageResponse struct {
    SenderID   int       `json:"sender_id" bson:"sender_id"`
    ReceiverID int       `json:"receiver_id,omitempty" bson:"receiver_id,omitempty"`
    Message    string    `json:"message" bson:"message"`
    Timestamp  time.Time `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
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