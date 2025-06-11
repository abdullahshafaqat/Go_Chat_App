package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *routerImpl) GetMessages(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID "})
		return
	}
	messages, err := r.messageservice.GetMessages(c, senderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
