package router

import (
	"net/http"

	"github.com/abdullahshafaqat/Go_ChatApp.git/models"
	"github.com/gin-gonic/gin"
)

func (r *Router) SignUp(c *gin.Context) {
	var req *models.UserSignup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	signup := models.UserSignup{
		Email:    req.Email,
		Password: req.Password,
	}
	user := r.AuthService.SignUp(c, &signup)
	c.JSON(200, gin.H{
		"message": "User created successfully",
		"user":    user.Username,
	})
}
