package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *routerImpl) GetMessages(c *gin.Context) {
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

	messages, err := r.messageservice.GetMessages(c, senderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range messages {
		if messages[i].Message == "" {
			messages[i].Message = "No content"
		}
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
