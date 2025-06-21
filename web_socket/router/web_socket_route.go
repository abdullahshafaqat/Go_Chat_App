package wsrouter

import (
	"log"
	"net/http"
	"strconv"

	websocketimpl "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/web_socket_impl"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *Router) HandleWebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		userIDStr, ok := userIDValue.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID is not a valid string"})
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return
		}

		log.Printf("WebSocket connection established for user: %d", userID)
		websocketimpl.HandleConnection(userID, conn, r.WebSocketService)
	}
}
