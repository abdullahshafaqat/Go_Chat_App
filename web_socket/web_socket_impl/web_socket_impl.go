package websocketservice

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func HandleConnection(userID int, conn *websocket.Conn, wsService WebSocketService) {
	if err := wsService.AddClient(userID, conn); err != nil {
		log.Printf("[HandleConnection] Failed to add client %d: %v\n", userID, err)
		conn.Close()
		return
	}
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	
	select {
	case <-ticker.C:
		log.Printf("[HandleConnection] Max connection time reached for user %d", userID)
	case <-func() <-chan struct{} {
		client, exists := wsService.GetClient(userID)
		if !exists || client == nil {

			ch := make(chan struct{})
			close(ch)
			return ch
		}
		return client.closeChan
	}():
		log.Printf("[HandleConnection] Close signal received for user %d", userID)
	}
}