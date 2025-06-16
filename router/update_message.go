package router

import (
	"net/http"
	"strconv"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

func (r *routerImpl) UpdateMessage(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	
	senderID, err := strconv.Atoi(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	var request models.UpdateMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	updatedMsg, err := r.messageservice.UpdateMessage(c, request.ID, senderID, request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Message updated successfully",
		"data":    updatedMsg,
	})
}
