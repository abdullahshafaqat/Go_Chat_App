package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	senderID, err := strconv.Atoi(userID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID format"})
		return
	}

	msg := models.Message{
		ID:         primitive.NewObjectID(), // Generate new ID here
		SenderID:   senderID,
		ReceiverID: req.ReceiverID,
		Message:    req.Content,
		Timestamp:  time.Now(),
	}

	if err := r.messageservice.SendMessage(c, &msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent successfully",
		"data":    msg,
	})
}