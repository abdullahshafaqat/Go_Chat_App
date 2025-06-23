package websocketimpl

import (
	"context"
	"errors"
	"sync"

	"github.com/abdullahshafaqat/Go_Chat_App.git/db/mongodb"
	wsmodels "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/models"
	"github.com/gorilla/websocket"
)

var (
	ErrClientNotFound    = errors.New("receiver not available")
	ErrConnectionClosed  = errors.New("connection closed")
	ErrInvalidMessage    = errors.New("invalid message format")
	ErrMessageTooLong    = errors.New("message too long")
	ErrConnectionTimeout = errors.New("connection timeout")
)

type Client struct {
	conn      *websocket.Conn
	send      chan []byte
	closeChan chan struct{}
}

type WebSocketService interface {
	AddClient(userID int, conn *websocket.Conn) error
	RemoveClient(userID int)
	GetClient(userID int) (*Client, bool)
	BroadcastMessage(ctx context.Context, senderID int, msg wsmodels.IncomingMessage) error
	GetOnlineUsers() []int
}

type webSocketService struct {
	clients map[int]*Client
	lock    sync.RWMutex
	message mongodb.Database
}

func NewWebSocketService(message mongodb.Database) WebSocketService {
	return &webSocketService{
		clients: make(map[int]*Client),
		message: message,
	}
}
