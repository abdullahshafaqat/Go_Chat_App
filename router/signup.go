package router

import (
	"net/http"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

func (r *routerImpl) SignUp(c *gin.Context) {
	var Newuser models.UserSignup
	if err := c.ShouldBindJSON(&Newuser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if err := r.authservice.SignUp(c, &Newuser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created", "user_name": Newuser.Username})
}
