package websocketservice

import (
	"log"
	"net/http"

	websocket "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket"
	"github.com/gin-gonic/gin"
)

func (s *serviceImpl) HandleConnections(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ws, err := s.manager.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := &websocket.Client{
		Conn:   ws,
		UserID: userID.(string),
		Send:   make(chan []byte, 256),
	}

	s.manager.register <- client

	go client.WritePump()
	go client.ReadPump(s.manager)
}
