package router

import (
	"net/http"
	"strconv"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

func (r *routerImpl) UpdateMessage(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenString := r.authservice.BearerToken(authHeader)

	userID, err := r.authservice.DecodeToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	isValid, _, err := r.authservice.Authorize(tokenString)
	if err != nil || !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	senderID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	messageID := c.Param("_id")

	var msg models.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	updatedMsg, err := r.messageservice.UpdateMessage(c, messageID, senderID, &msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Message updated successfully",
		"data":    updatedMsg,
	})
}
