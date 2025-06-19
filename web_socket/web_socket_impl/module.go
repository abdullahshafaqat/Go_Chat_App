package websocketimpl

import (
	"sync"

	"github.com/abdullahshafaqat/Go_Chat_App.git/db/mongodb"
	"github.com/gorilla/websocket"
)

type WebSocketService interface {
	AddClient(userID int, conn *websocket.Conn)
	RemoveClient(userID int)
	GetClient(userID int) (*websocket.Conn, bool)
	BroadcastMessage(senderID, receiverID int, message string) error
}

type webSocketService struct {
	clients map[int]*websocket.Conn
	lock    sync.RWMutex
	message mongodb.Database
}

func NewWebSocketService(message mongodb.Database) WebSocketService {
	return &webSocketService{
		clients: make(map[int]*websocket.Conn),
		message: message,
	}
}
