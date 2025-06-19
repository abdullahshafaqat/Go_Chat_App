package websocketimpl

import (
	"errors"
	"log"

	"github.com/gorilla/websocket"
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
    w.lock.RLock()
    defer w.lock.RUnlock()

    conn, exists := w.clients[receiverID]
    if !exists {
        return ErrClientNotFound
    }

    msg := map[string]interface{}{
        "sender_id": senderID,
        "message":   message,
    }

    return conn.WriteJSON(msg)
}

// Custom errors
var (
    ErrClientNotFound = errors.New("receiver client not found")
)