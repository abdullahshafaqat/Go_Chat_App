package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	wsmodels "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/models"

	"github.com/gin-gonic/gin"
)

func (r *routerImpl) SendMessage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	fmt.Println("Sending message to:", req.ReceiverID)

	senderID, err := strconv.Atoi(userID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID format"})
		return
	}

	
	incoming := wsmodels.IncomingMessage{
		ReceiverID: req.ReceiverID,
		Message:    req.Content,
	}


	if err := r.webSocketService.BroadcastMessage(c.Request.Context(), senderID, incoming); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}
