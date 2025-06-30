package websocketservice

// import (
// 	"context"

// 	wsmodels "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/models"
// 	"github.com/gorilla/websocket"
// )

// // import (
// // 	"context"
// // 	"encoding/json"
// // 	"errors"
// // 	"fmt"
// // 	"log"
// // 	"time"

// // 	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
// // 	wsmodels "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/models"
// // 	"go.mongodb.org/mongo-driver/bson/primitive"

// //	"github.com/gorilla/websocket"
// //
// // )
// func (w *webSocketService) AddClient(userID int, conn *websocket.Conn) error {}

// // 	w.lock.Lock()
// // 	defer w.lock.Unlock()

// // 	logPrefix := "[AddClient]"

// // 	if w.activeConnections >= MaxConnections {
// // 		log.Printf("%s Max connections reached (%d)", logPrefix, MaxConnections)
// // 		return errors.New("maximum connections reached")
// // 	}

// // 	if existingClient, exists := w.clients[userID]; exists {
// // 		log.Printf("%s Closing existing connection for user %d", logPrefix, userID)
// // 		w.closeClientUnsafe(userID, existingClient)
// // 	} else {

// // 		w.activeConnections++
// // 	}

// // 	client := &Client{
// // 		conn:      conn,
// // 		send:      make(chan []byte, ConnBufferSize),
// // 		closeChan: make(chan struct{}),
// // 	}

// // 	conn.SetReadLimit(MaxMessageSize)
// // 	conn.SetReadDeadline(time.Now().Add(PongWait))
// // 	conn.SetPongHandler(func(string) error {
// // 		conn.SetReadDeadline(time.Now().Add(PongWait))
// // 		return nil
// // 	})

// // 	w.clients[userID] = client

// // 	currentUsers := make([]int, 0, len(w.clients))
// // 	for id := range w.clients {
// // 		currentUsers = append(currentUsers, id)
// // 	}
// // 	log.Printf("%s Added client for user %d (total connections: %d, users: %v)",
// // 		logPrefix, userID, w.activeConnections, currentUsers)

// // 	go w.readPump(userID, client)
// // 	go w.writePump(userID, client)

// // 	w.notifyConnectionStatus(userID, true)

// // 	return nil
// // }

// func (w *webSocketService) closeClientUnsafe(userID int, client *Client) {}

// // 	select {
// // 	case <-client.closeChan:

// // 	default:
// // 		close(client.closeChan)
// // 	}

// // 	client.conn.Close()
// // 	delete(w.clients, userID)
// // 	w.activeConnections--
// // }

// func (w *webSocketService) readPump(userID int, client *Client) {}

// // 	logPrefix := "[ReadPump]"
// // 	log.Printf("%s Starting for user %d", logPrefix, userID)

// // 	defer func() {
// // 		if r := recover(); r != nil {
// // 			log.Printf("%s PANIC for user %d: %v", logPrefix, userID, r)
// // 		}
// // 		w.RemoveClient(userID)
// // 		log.Printf("%s Ended for user %d", logPrefix, userID)
// // 	}()

// // 	for {
// // 		_, messageBytes, err := client.conn.ReadMessage()
// // 		if err != nil {
// // 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// // 				log.Printf("%s Unexpected close error for user %d: %v", logPrefix, userID, err)
// // 			} else {
// // 				log.Printf("%s Read error for user %d: %v", logPrefix, userID, err)
// // 			}
// // 			return
// // 		}

// // 		client.conn.SetReadDeadline(time.Now().Add(PongWait))

// // 		if err := w.processIncomingMessage(userID, messageBytes); err != nil {
// // 			log.Printf("%s Error processing message from user %d: %v", logPrefix, userID, err)

// // 			errorMsg := map[string]interface{}{
// // 				"type":    "error",
// // 				"message": err.Error(),
// // 			}
// // 			jsonMsg := mustJSON(errorMsg)
// // 			select {
// // 			case client.send <- jsonMsg:
// // 			default:
// // 				log.Printf("%s Could not send error to user %d: channel full", logPrefix, userID)
// // 			}
// // 		}
// // 	}
// // }

// func (w *webSocketService) processIncomingMessage(senderID int, messageBytes []byte) error {}

// // 	logPrefix := "[ProcessMessage]"

// // 	var genericMsg map[string]interface{}
// // 	if err := json.Unmarshal(messageBytes, &genericMsg); err != nil {
// // 		log.Printf("%s JSON unmarshal error from user %d: %v", logPrefix, senderID, err)
// // 		return fmt.Errorf("invalid JSON format")
// // 	}

// // 	msgType, hasType := genericMsg["type"].(string)
// // 	if !hasType {
// // 		msgType = "message"
// // 	}

// // 	log.Printf("%s Received from user %d, type: %s", logPrefix, senderID, msgType)

// // 	switch msgType {
// // 	case "message", "":
// // 		var incomingMsg wsmodels.IncomingMessage
// // 		if err := json.Unmarshal(messageBytes, &incomingMsg); err != nil {
// // 			log.Printf("%s Failed to unmarshal message from user %d: %v", logPrefix, senderID, err)
// // 			return fmt.Errorf("invalid message format")
// // 		}
// // 		return w.BroadcastMessage(context.Background(), senderID, incomingMsg)

// // 	case "ping":
// // 		return w.handlePingMessage(senderID)

// // 	case "typing":
// // 		return w.handleTypingMessage(senderID, genericMsg)

// // 	default:
// // 		log.Printf("%s Unknown message type '%s' from user %d", logPrefix, msgType, senderID)
// // 		return fmt.Errorf("unknown message type: %s", msgType)
// // 	}
// // }

// func (w *webSocketService) handlePingMessage(userID int) error {}

// // 	w.lock.RLock()
// // 	client, exists := w.clients[userID]
// // 	w.lock.RUnlock()

// // 	if !exists {
// // 		return fmt.Errorf("client not found")
// // 	}

// // 	pongMsg := map[string]string{"type": "pong"}
// // 	jsonMsg := mustJSON(pongMsg)

// // 	select {
// // 	case client.send <- jsonMsg:
// // 		return nil
// // 	default:
// // 		return fmt.Errorf("client send channel full")
// // 	}
// // }

// func (w *webSocketService) handleTypingMessage(senderID int, genericMsg map[string]interface{}) error {
// }

// // 	receiverIDFloat, ok := genericMsg["receiver_id"].(float64)
// // 	if !ok {
// // 		return fmt.Errorf("receiver ID required for typing message")
// // 	}
// // 	receiverID := int(receiverIDFloat)

// // 	if receiverID == 0 {
// // 		return fmt.Errorf("receiver ID required for typing message")
// // 	}

// // 	w.lock.RLock()
// // 	receiverClient, exists := w.clients[receiverID]
// // 	w.lock.RUnlock()

// // 	if !exists {
// // 		return nil
// // 	}

// // 	typing := true
// // 	if typingStatus, ok := genericMsg["typing"].(bool); ok {
// // 		typing = typingStatus
// // 	} else if message, ok := genericMsg["message"].(string); ok {
// // 		typing = message == "start"

// // 	}

// // 	typingMsg := map[string]interface{}{
// // 		"type":      "typing",
// // 		"sender_id": senderID,
// // 		"typing":    typing,
// // 	}

// // 	jsonMsg := mustJSON(typingMsg)

// //		select {
// //		case receiverClient.send <- jsonMsg:
// //			return nil
// //		default:
// //			return fmt.Errorf("receiver send channel full")
// //		}
// //	}
// func (w *webSocketService) writePump(userID int, client *Client) {}

// // 	logPrefix := "[WritePump]"
// // 	log.Printf("%s Starting for user %d", logPrefix, userID)

// // 	ticker := time.NewTicker(PingPeriod)
// // 	defer func() {
// // 		ticker.Stop()
// // 		if r := recover(); r != nil {
// // 			log.Printf("%s PANIC for user %d: %v", logPrefix, userID, r)
// // 		}
// // 		log.Printf("%s Ended for user %d", logPrefix, userID)
// // 	}()

// // 	for {
// // 		select {
// // 		case message, ok := <-client.send:
// // 			client.conn.SetWriteDeadline(time.Now().Add(WriteWait))
// // 			if !ok {
// // 				log.Printf("%s send channel closed for user %d. Closing WS.", logPrefix, userID)
// // 				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
// // 				return
// // 			}

// // 			if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
// // 				log.Printf("%s Write error for user %d: %v", logPrefix, userID, err)
// // 				return
// // 			}

// // 		case <-ticker.C:
// // 			client.conn.SetWriteDeadline(time.Now().Add(WriteWait))
// // 			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// // 				log.Printf("%s Ping failed for user %d: %v", logPrefix, userID, err)
// // 				return
// // 			}

// // 		case <-client.closeChan:
// // 			log.Printf("%s closeChan closed for user %d — terminating writePump.", logPrefix, userID)
// // 			return
// // 		}
// // 	}
// // }

// func (w *webSocketService) IsUserOnline(userID int) bool {}

// // 	w.lock.RLock()
// // 	defer w.lock.RUnlock()
// // 	_, exists := w.clients[userID]
// // 	return exists
// // }
// // func (w *webSocketService) RemoveClient(userID int) {
// // 	w.lock.Lock()
// // 	defer w.lock.Unlock()

// // 	logPrefix := "[RemoveClient]"
// // 	client, exists := w.clients[userID]
// // 	if exists {

// // 		select {
// // 		case <-client.closeChan:
// // 		default:
// // 			close(client.closeChan)
// // 		}
// // 		client.conn.Close()

// // 		delete(w.clients, userID)
// // 		w.activeConnections--

// // 		log.Printf("%s Removed client for user %d (remaining: %d)",
// // 			logPrefix, userID, w.activeConnections)

// // 		go w.notifyConnectionStatus(userID, false)
// // 	} else {
// // 		log.Printf("%s No client found for user %d", logPrefix, userID)
// // 	}
// // }

// func (w *webSocketService) notifyConnectionStatus(userID int, online bool) {}

// // 	w.lock.RUnlock()

// // 	for id, client := range clients {
// // 		select {
// // 		case client.send <- mustJSON(status):
// // 			log.Printf("Notified user %d about user %d status: %s", id, userID, status.Status)
// // 		default:
// // 			log.Printf("Failed to notify user %d about user %d status - channel full", id, userID)
// // 		}
// // 	}
// // }

// func (w *webSocketService) GetClient(userID int) (*Client, bool) {}

// // 	w.lock.RLock()
// // 	defer w.lock.RUnlock()
// // 	client, exists := w.clients[userID]
// // 	return client, exists
// // }

// //	func mustJSON(v interface{}) []byte {
// //		b, err := json.Marshal(v)
// //		if err != nil {
// //			log.Printf("JSON marshaling error: %v", err)
// //			return []byte(`{"type":"error","message":"internal server error"}`)
// //		}
// //		return b
// //	}
// func (w *webSocketService) BroadcastMessage(ctx context.Context, senderID int, msg wsmodels.IncomingMessage) error {
// }

// // 	log.Printf("[Broadcast] Init | Sender: %d → Receiver: %d | Msg: %q (%d bytes)",
// // 		senderID, msg.ReceiverID, msg.Message, len(msg.Message))

// // 	if len(msg.Message) > MaxMessageSize {
// // 		log.Printf("[Broadcast] Validation failed: message too long")
// // 		return ErrMessageTooLong
// // 	}
// // 	if msg.ReceiverID == 0 {
// // 		log.Printf("[Broadcast] Validation failed: receiver id is required")
// // 		return errors.New("receiver id is required")
// // 	}

// // 	msgToStore := &models.Message{
// // 		ID:         primitive.NewObjectID(),
// // 		SenderID:   senderID,
// // 		ReceiverID: msg.ReceiverID,
// // 		Message:    msg.Message,
// // 		Timestamp:  time.Now().UTC(),
// // 	}

// // 	start := time.Now()
// // 	if err := w.message.InsertMessage(ctx, msgToStore); err != nil {
// // 		log.Printf("[Broadcast] DB insert failed | Error: %v | Duration: %v", err, time.Since(start))
// // 		return fmt.Errorf("failed to store message")
// // 	}
// // 	log.Printf("[Broadcast] DB insert success | Sender: %d | Receiver: %d | Duration: %v",
// // 		senderID, msg.ReceiverID, time.Since(start))

// // 	wsMsg := wsmodels.Message{
// // 		SenderID:   senderID,
// // 		ReceiverID: msg.ReceiverID,
// // 		Message:    msg.Message,
// // 		Timestamp:  time.Now().UTC(),
// // 	}
// // 	jsonMsg := mustJSON(wsMsg)

// // 	w.lock.RLock()
// // 	receiverClient, receiverExists := w.clients[msg.ReceiverID]
// // 	senderClient, senderExists := w.clients[senderID]
// // 	w.lock.RUnlock()

// // 	if receiverExists {
// // 		select {
// // 		case receiverClient.send <- jsonMsg:
// // 			log.Printf("[Broadcast] enqueued to receiver %d", msg.ReceiverID)
// // 		case <-time.After(time.Second * 1): // Add timeout to prevent blocking
// // 			log.Printf("[Broadcast] receiver queue timeout (%d)", msg.ReceiverID)
// // 			return ErrConnectionTimeout
// // 		}
// // 	} else {
// // 		log.Printf("[Broadcast] receiver offline (%d)", msg.ReceiverID)
// // 	}

// // 	if senderExists && senderID != msg.ReceiverID {
// // 		select {
// // 		case senderClient.send <- jsonMsg:
// // 			log.Printf("[Broadcast] echoed to sender %d", senderID)
// // 		case <-time.After(time.Millisecond * 500): // Shorter timeout for echo
// // 			log.Printf("[Broadcast] sender queue timeout (%d)", senderID)
// // 		}
// // 	}

// // 	return nil
// // }

// func (w *webSocketService) GetAllOnlineUsers() []int {}

// // 	w.lock.RLock()
// // 	defer w.lock.RUnlock()

// // 	users := make([]int, 0, len(w.clients))
// // 	for userID := range w.clients {
// // 		users = append(users, userID)
// // 	}
// // 	return users
// // }

// func (w *webSocketService) BroadcastToAll(message interface{}) {}

// // 	jsonMsg := mustJSON(message)

// // 	w.lock.RLock()
// // 	clients := make(map[int]*Client, len(w.clients))
// // 	for id, client := range w.clients {
// // 		clients[id] = client
// // 	}
// // 	w.lock.RUnlock()

// // 	for userID, client := range clients {
// // 		select {
// // 		case client.send <- jsonMsg:
// // 			log.Printf("Broadcasted message to user %d", userID)
// // 		default:
// // 			log.Printf("Failed to broadcast to user %d - channel full", userID)
// // 		}
// // 	}
// // }
