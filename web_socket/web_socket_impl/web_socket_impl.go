package websocketimpl

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type IncomingMessage struct {
	ReceiverID int    `json:"receiver_id"`
	Message    string `json:"message"`
}

func HandleConnection(userID int, conn *websocket.Conn, wsService WebSocketService) {

	wsService.AddClient(userID, conn)

	defer func() {
		wsService.RemoveClient(userID)
		conn.Close()
		log.Printf("Connection closed for user: %d", userID)
	}()

	welcomeMsg := map[string]interface{}{
		"type":    "connection",
		"message": "Connected successfully",
		"user_id": userID,
	}
	if err := conn.WriteJSON(welcomeMsg); err != nil {
		log.Printf("Error sending welcome message to user %d: %v", userID, err)
		return
	}

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for user %d: %v", userID, err)
			}
			break
		}

		log.Printf("Received message from user %d: %s", userID, string(msgBytes))

		var incomingMsg IncomingMessage
		if err := json.Unmarshal(msgBytes, &incomingMsg); err != nil {
			log.Printf("Error parsing message from user %d: %v", userID, err)

			errorMsg := map[string]interface{}{
				"type":    "error",
				"message": "Invalid message format",
			}
			conn.WriteJSON(errorMsg)
			continue
		}

		if incomingMsg.ReceiverID == 0 || incomingMsg.Message == "" {
			log.Printf("Invalid message from user %d: missing receiver_id or message", userID)

			errorMsg := map[string]interface{}{
				"type":    "error",
				"message": "receiver_id and message are required",
			}
			conn.WriteJSON(errorMsg)
			continue
		}

		err = wsService.BroadcastMessage(userID, incomingMsg.ReceiverID, incomingMsg.Message)
		if err != nil {
			log.Printf("Error broadcasting message from user %d to user %d: %v",
				userID, incomingMsg.ReceiverID, err)

			errorMsg := map[string]interface{}{
				"type":    "error",
				"message": "Failed to deliver message: " + err.Error(),
			}
			conn.WriteJSON(errorMsg)
		} else {

			confirmMsg := map[string]interface{}{
				"type":    "confirmation",
				"message": "Message delivered successfully",
			}
			conn.WriteJSON(confirmMsg)
		}
	}
}
