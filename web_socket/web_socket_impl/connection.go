package websocketservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	wsmodels "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/models"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/websocket"
)

func (w *webSocketService) AddClient(userID int, conn *websocket.Conn) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	logPrefix := "[AddClient]"

	defer func() {
		if r := recover(); r != nil {
			log.Printf("%s PANIC while adding client for user %d: %v", logPrefix, userID, r)
		}
	}()

	if w.activeConnections >= MaxConnections {
		log.Printf("%s Max connections reached (%d)", logPrefix, MaxConnections)
		return errors.New("maximum connections reached")
	}

	if existingClient, exists := w.clients[userID]; exists {
		log.Printf("%s Closing existing client for user %d", logPrefix, userID)
		select {
		case <-existingClient.closeChan:
			log.Printf("%s closeChan already closed for user %d", logPrefix, userID)
		default:
			close(existingClient.closeChan)
		}
		existingClient.conn.Close()
		delete(w.clients, userID)
		w.activeConnections--
	}

	client := &Client{
		conn:      conn,
		send:      make(chan []byte, ConnBufferSize),
		closeChan: make(chan struct{}),
	}

	conn.SetReadLimit(MaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(PongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	w.clients[userID] = client
	log.Println("clients: ", client)
	log.Println("87", w.clients)
	w.activeConnections++
	log.Printf("%s Added client for user %d (total: %d)", logPrefix, userID, w.activeConnections)

	go w.writePump(userID, client)
	w.notifyConnectionStatus(userID, true)

	return nil
}

func (w *webSocketService) writePump(userID int, client *Client) {
	logPrefix := "[WritePump]"
	log.Printf("%s Starting for user %d", logPrefix, userID)

	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		if r := recover(); r != nil {
			log.Printf("%s PANIC for user %d: %v", logPrefix, userID, r)
		}
		w.RemoveClient(userID)
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				log.Printf("%s send channel closed for user %d. Closing WS.", logPrefix, userID)
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w := client.conn
			if err := w.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("%s Write error for user %d: %v", logPrefix, userID, err)
				return
			}

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("%s Ping failed for user %d: %v", logPrefix, userID, err)
				return
			}
		case <-client.closeChan:
			log.Printf("%s closeChan closed for user %d — terminating writePump.", logPrefix, userID)
			return
		}
	}
}

func (w *webSocketService) IsUserOnline(userID int) bool {
	w.lock.RLock()
	defer w.lock.RUnlock()
	_, exists := w.clients[userID]
	return exists
}
func (w *webSocketService) RemoveClient(userID int) {
	w.lock.Lock()
	defer w.lock.Unlock()

	logPrefix := "[RemoveClient]"
	client, exists := w.clients[userID]
	if exists {
		select {
		case <-client.closeChan:
			log.Printf("%s closeChan already closed for user %d", logPrefix, userID)
		default:
			close(client.closeChan)
		}

		client.conn.Close()
		delete(w.clients, userID)
		w.activeConnections--
		log.Printf("%s Removed client for user %d (remaining: %d)", logPrefix, userID, w.activeConnections)
	} else {
		log.Printf("%s No client found for user %d", logPrefix, userID)
	}
}

func (w *webSocketService) notifyConnectionStatus(userID int, online bool) {
	status := wsmodels.ConnectionStatus{
		UserID: userID,
		Status: "online",
	}
	if !online {
		status.Status = "offline"
		status.Message = "User disconnected"
	}

	w.lock.RLock()
	defer w.lock.RUnlock()

	for id, client := range w.clients {
		if id == userID {
			continue
		}
		select {
		case client.send <- mustJSON(status):
		default:
			log.Printf("Failed to notify user %d - channel full", id)
		}
	}
}

func (w *webSocketService) GetClient(userID int) (*Client, bool) {
	w.lock.RLock()
	defer w.lock.RUnlock()
	client, exists := w.clients[userID]
	return client, exists
}

func mustJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
func (w *webSocketService) BroadcastMessage(ctx context.Context, senderID int, msg wsmodels.IncomingMessage) error {
	log.Printf("[Broadcast] Init | Sender: %d → Receiver: %d | Msg: %q (%d bytes)",
		senderID, msg.ReceiverID, msg.Message, len(msg.Message))

	if len(msg.Message) > MaxMessageSize {
		log.Printf("[Broadcast] Validation failed: message too long")
		return ErrMessageTooLong
	}
	if msg.ReceiverID == 0 {
		log.Printf("[Broadcast] Validation failed: receiver id is required")
		return errors.New("receiver id is required")
	}

	msgToStore := &models.Message{
		ID:         primitive.NewObjectID(),
		SenderID:   senderID,
		ReceiverID: msg.ReceiverID,
		Message:    msg.Message,
		Timestamp:  time.Now().UTC(),
	}

	start := time.Now()
	if err := w.message.InsertMessage(ctx, msgToStore); err != nil {
		log.Printf("[Broadcast] DB insert failed | Error: %v | Duration: %v", err, time.Since(start))
		return fmt.Errorf("failed to store message")
	}
	log.Printf("[Broadcast] DB insert success | Sender: %d | Receiver: %d | Duration: %v",
		senderID, msg.ReceiverID, time.Since(start))

	wsMsg := wsmodels.Message{
		SenderID:   senderID,
		ReceiverID: msg.ReceiverID,
		Message:    msg.Message,
		Timestamp:  time.Now().UTC(),
	}
	jsonMsg := mustJSON(wsMsg)

	w.lock.RLock()
	receiverClient, receiverExists := w.clients[msg.ReceiverID]
	w.lock.RUnlock()

	if receiverExists {
		select {
		case receiverClient.send <- jsonMsg:
			log.Printf("[Broadcast] enqueued to receiver %d", msg.ReceiverID)
		default:
			log.Printf("[Broadcast] receiver queue full (%d)", msg.ReceiverID)
			return ErrConnectionTimeout
		}
	} else {
		log.Printf("[Broadcast] receiver offline (%d)", msg.ReceiverID)
	}

	w.lock.RLock()
	senderClient, senderExists := w.clients[senderID]
	w.lock.RUnlock()

	if senderExists && senderID != msg.ReceiverID {
		select {
		case senderClient.send <- jsonMsg:
			log.Printf("[Broadcast] echoed to sender %d", senderID)
		default:
			log.Printf("[Broadcast] sender queue full (%d)", senderID)
		}
	}
	return nil
}
