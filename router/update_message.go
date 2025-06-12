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

	var request models.UpdateMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	updatedMsg, _ := r.messageservice.UpdateMessage(c, request.ID, senderID, request.Message)

	c.JSON(http.StatusOK, gin.H{
		"Message updated successfully": updatedMsg,
	})
}
