package router

import (
	"net/http"
	"time"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

func (r *routerImpl) SendMessage(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := r.authservice.BearerToken(authHeader)

	isValid, _, err := r.authservice.Authorize(token)
	if err != nil || !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
		return
	}

	var msg models.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	msg.Timestamp = time.Now()

	if err := r.messageservice.SendMessage(c, &msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}
