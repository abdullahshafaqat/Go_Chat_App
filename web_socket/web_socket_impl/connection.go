package websocketimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	wsmodels "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (w *webSocketService) AddClient(userID int, conn *websocket.Conn) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	
	if existingClient, exists := w.clients[userID]; exists {
		close(existingClient.closeChan)
		existingClient.conn.Close()
		delete(w.clients, userID)
	}

	client := &Client{
		conn:      conn,
		send:      make(chan []byte, 256),
		closeChan: make(chan struct{}),
	}

	conn.SetReadLimit(MaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(PongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})
	conn.SetPingHandler(func(message string) error {
		conn.SetWriteDeadline(time.Now().Add(WriteWait))
		return conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(WriteWait))
	})

	w.clients[userID] = client
	log.Printf("Client connected - UserID: %d, Active connections: %d", userID, len(w.clients))

	
	go w.writePump(userID, client)

	
	w.notifyConnectionStatus(userID, true)
	return nil
}

func (w *webSocketService) writePump(userID int, client *Client) {
	ticker := time.NewTicker(PingInterval)
	defer func() {
		ticker.Stop()
		client.conn.Close()
		w.RemoveClient(userID)
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Write error for user %d: %v", userID, err)
				return
			}
			writer.Write(message)
			writer.Close()

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Ping failed for user %d: %v", userID, err)
				return
			}

		case <-client.closeChan:
			return
		}
	}
}

func (w *webSocketService) RemoveClient(userID int) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if client, exists := w.clients[userID]; exists {
		close(client.closeChan)
		client.conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			time.Now().Add(WriteWait),
		)
		client.conn.Close()
		delete(w.clients, userID)
		log.Printf("Client disconnected - UserID: %d", userID)

	
		w.notifyConnectionStatus(userID, false)
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

func (w *webSocketService) BroadcastMessage(ctx context.Context,senderID int,msg wsmodels.IncomingMessage,) error{
    if len(msg.Message) > MaxMessageSize {
        return ErrMessageTooLong
    }


    msgToStore := &models.Message{
        ID:         primitive.NewObjectID(),
        SenderID:   senderID,
        ReceiverID: msg.ReceiverID,
        Message:    msg.Message,
        Timestamp:  time.Now().UTC(),
    }


    ginCtx := &gin.Context{
        Request: &http.Request{
            URL:    &url.URL{Path: "/ws"},
            Method: "WS",
        },
    }
    ginCtx.Request = ginCtx.Request.WithContext(ctx)

    if err := w.message.InsertMessage(ginCtx, msgToStore); err != nil {
        log.Printf("Failed to store message: %v", err)
        return fmt.Errorf("failed to store message")
    }

    wsMsg := wsmodels.MessageResponse{
        SenderID:   senderID,
        ReceiverID: msg.ReceiverID,
        Message:    msg.Message,
        Timestamp:  time.Now().UTC(),
    }
    w.lock.RLock()
    receiverClient, receiverExists := w.clients[msg.ReceiverID]
    w.lock.RUnlock()

    if receiverExists {
        select {
        case receiverClient.send <- mustJSON(wsMsg):
        default:
            return ErrConnectionTimeout
        }
    }


    w.lock.RLock()
    senderClient, senderExists := w.clients[senderID]
    w.lock.RUnlock()

    if senderExists {
        select {
        case senderClient.send <- mustJSON(wsMsg):
        default:
            log.Printf("Could not send message back to sender %d", senderID)
        }
    }

    return nil
}


func (w *webSocketService) GetOnlineUsers() []int {
	w.lock.RLock()
	defer w.lock.RUnlock()

	users := make([]int, 0, len(w.clients))
	for userID := range w.clients {
		users = append(users, userID)
	}
	return users
}

func mustJSON(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
