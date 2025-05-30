package authservice

import (
	"github.com/gin-gonic/gin"
	"github.com/abdullahshafaqat/Go_ChatApp.git/models"
)

type AuthService interface {
	SignUp(c *gin.Context, req *models.UserSignup) *models.UserSignup
}