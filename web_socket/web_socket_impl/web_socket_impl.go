package websocketservice

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
		log.Printf("[HandleConnection] Failed to add client %d: %v\n", userID, err)
		conn.Close()
		return
	}
	defer func() {
		wsService.RemoveClient(userID)
		conn.Close()
		log.Printf("[HandleConnection] Connection closed for user %d\n", userID)
	}()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[HandleConnection] User %d disconnected unexpectedly: %v\n", userID, err)
			}
			break
		}

		var incoming wsmodels.IncomingMessage
		if err := json.Unmarshal(msgBytes, &incoming); err != nil {
			log.Printf("[HandleConnection] Invalid message format from user %d\n", userID)
			sendError(conn, "Invalid message format")
			continue
		}

		log.Printf("[HandleConnection] User %d sending to %d: %s\n",
			userID, incoming.ReceiverID, incoming.Message)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = wsService.BroadcastMessage(ctx, userID, incoming)
		cancel()

		if err != nil {
			log.Printf("[HandleConnection] Failed to deliver message: %v\n", err)
			sendError(conn, "Failed to deliver message")
		} else {
			sendConfirmation(conn, "Message delivered")
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

func sendJSON(conn *websocket.Conn, v interface{}) {
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second)) // WriteWait
	err := conn.WriteJSON(v)
	if err != nil {
		log.Printf("[sendJSON] failed to write: %v\n", err)
	}
}
