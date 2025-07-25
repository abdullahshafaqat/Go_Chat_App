package wsrouter

import (
	"log"
	"net/http"
	"time"

	// Adjust imports as needed
	websocketimpl "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/web_socket_impl"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		allowedOrigin := r.Header.Get("Origin") == "http://localhost:8080"
		if !allowedOrigin {
			log.Printf("[WSHandler] Rejected connection from origin: %s", r.Header.Get("Origin"))
		}
		return allowedOrigin
	},
}


func (r *Router) HandleWebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		logPrefix := "[WSHandler]"

		userID, exists := c.Get("userID")
		if !exists {
			log.Printf("%s User ID not found in context", logPrefix)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			return
		}
		userIDInt, ok := userID.(int)
		if !ok {
			log.Printf("%s Invalid user ID format in context: %v", logPrefix, userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}

		
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("%s WebSocket upgrade failed for user %d: %v", logPrefix, userIDInt, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket connection failed"})
			return
		}
		log.Printf("%s WebSocket connection established for user %d (took %v)",
			logPrefix, userIDInt, time.Since(startTime))

	
		go func(wsService websocketimpl.WebSocketService) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("%s Recovered panic in HandleConnection for user %d: %v",
						logPrefix, userIDInt, r)
				}
			}()

			websocketimpl.HandleConnection(userIDInt, conn, wsService)
		}(r.WebSocketService)
	}
}
