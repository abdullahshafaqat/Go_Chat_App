package router

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

func (r *Router) StartWebSocketServer(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		log.Println("userID not found in context")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade failed:", err)
		return
	}

	log.Printf("WebSocket route hit with userID param: %s", userID)
	go func() {
		err = r.WebSocketSvc.AddConn(userID.(string), wsConn, c)
		if err != nil {
			log.Println("AddConn error:", err)
		}
	}()
}

func (r *Router) SaveBackendConnection(c *gin.Context) {
	userID := "-1" // backend/system default ID

	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade failed:", err)
		return
	}

	log.Printf("Backend WebSocket route hit with userID param: %s", userID)
	go func() {
		err = r.WebSocketSvc.AddConn(userID, wsConn, c)
		if err != nil {
			log.Println("AddConn error:", err)
		}
	}()
}
