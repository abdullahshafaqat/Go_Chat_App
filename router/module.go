package router

import (
	authservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/auth_service"
	messageservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/message_service"
	websocketservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/web_socket_service"
	"github.com/gin-gonic/gin"
)

type Router interface {
	DefineRoutes(r *gin.Engine)
}

type routerImpl struct {
	authservice    authservice.Service
	messageservice messageservice.Service
	wsService      websocketservice.Service
}

func NewRouter(authservice authservice.Service, messageservice messageservice.Service, ws websocketservice.Service) Router {
	return &routerImpl{
		authservice:    authservice,
		messageservice: messageservice,
		wsService:      ws,
	}
}
