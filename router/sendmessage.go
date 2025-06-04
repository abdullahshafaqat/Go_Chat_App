package router

import (
	"net/http"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

func (r *routerImpl) SendMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}


}
