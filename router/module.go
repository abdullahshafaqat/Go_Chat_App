package router

import (
	authservice "github.com/abdullahshafaqat/Go_ChatApp.git/api/auth_service"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine      *gin.Engine
	AuthService authservice.AuthService
}

