package websocketimpl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func (w *webSocketService) AddClient(userID int, conn *websocket.Conn) {
	w.lock.Lock()
	defer w.lock.Unlock()
	if existingConn, exists := w.clients[userID]; exists {
		existingConn.Close()
	}
	w.clients[userID] = conn
	log.Printf("Client connected - UserID: %d, Active connections: %d", userID, len(w.clients))
}

func (w *webSocketService) RemoveClient(userID int) {
	w.lock.Lock()
	defer w.lock.Unlock()
	if conn, exists := w.clients[userID]; exists {
		conn.Close()
		delete(w.clients, userID)
		log.Printf("Client disconnected - UserID: %d", userID)
	}
}

func (w *webSocketService) GetClient(userID int) (*websocket.Conn, bool) {
	w.lock.RLock()
	defer w.lock.RUnlock()
	conn, exists := w.clients[userID]
	return conn, exists
}
func (w *webSocketService) BroadcastMessage(senderID, receiverID int, message string) error {
	if message == "" {
		return errors.New("message cannot be empty")
	}


	msgToStore := &models.Message{
		ID:         primitive.NewObjectID(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    message,
		Timestamp:  time.Now().UTC(),
	}


	ginCtx := &gin.Context{
		Request: &http.Request{
			Method: "POST",
			URL:    &url.URL{Path: "/ws-message"},
		},
	}
	ginCtx.Request = ginCtx.Request.WithContext(context.Background())


	if err := w.message.InsertMessage(ginCtx, msgToStore); err != nil {
		log.Printf("Failed to store message in MongoDB: %v", err)
		return fmt.Errorf("failed to store message: %w", err)
	}


	w.lock.RLock()
	defer w.lock.RUnlock()

	conn, exists := w.clients[receiverID]
	if !exists {
		return ErrClientNotFound
	}
	wsMsg := map[string]interface{}{
		"sender_id":   senderID,
		"receiver_id": receiverID,
		"message":     message,
		"timestamp":   time.Now().UTC(),
	}

	if err := conn.WriteJSON(wsMsg); err != nil {
		if websocket.IsCloseError(err, websocket.CloseNormalClosure,
			websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			return ErrConnectionClosed
		}
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
