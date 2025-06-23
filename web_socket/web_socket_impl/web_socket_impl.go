package websocketimpl

import (
	"context"
	"encoding/json"
	"log"
	"time"

	wsmodels "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/models"
	"github.com/gorilla/websocket"
)
func HandleConnection(userID int, conn *websocket.Conn, wsService WebSocketService) {

	if err := wsService.AddClient(userID, conn); err != nil {
		log.Printf("Failed to add client %d: %v", userID, err)
		conn.Close()
		return
	}

	defer func() {
		wsService.RemoveClient(userID)
		conn.Close()
		log.Printf("Connection closed for user %d", userID)
	}()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("User %d unexpected disconnect: %v", userID, err)
			}
			break
		}

		var incomingMsg wsmodels.IncomingMessage
		if err := json.Unmarshal(msgBytes, &incomingMsg); err != nil {
			sendError(conn, "Invalid message format")
			continue
		}

		
		var lastErr error
		for i := 0; i < MaxRetryCount; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = wsService.BroadcastMessage(ctx, userID, incomingMsg)
			cancel()

			if err == nil {
				sendConfirmation(conn, "Message delivered")
				break
			}

			lastErr = err
			log.Printf("Broadcast attempt %d failed for user %d: %v", i+1, userID, err)
			time.Sleep(time.Second * time.Duration(i+1))
		}

		if lastErr != nil {
			log.Printf("Final broadcast failed for user %d: %v", userID, lastErr)
			sendError(conn, "Failed to deliver message")
		}
	}
}

func sendError(conn *websocket.Conn, message string) {
	sendJSON(conn, map[string]interface{}{
		"type":    "error",
		"message": message,
	})
}

func sendConfirmation(conn *websocket.Conn, message string) {
	sendJSON(conn, map[string]interface{}{
		"type":    "confirmation",
		"message": message,
	})
}

func sendJSON(conn *websocket.Conn, v interface{}) error {
	conn.SetWriteDeadline(time.Now().Add(WriteWait))
	return conn.WriteJSON(v)
}