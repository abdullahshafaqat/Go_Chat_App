package router

import (
	authservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/auth_service"
	messageservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/message_service"
	websocketservice "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/web_socket_impl"
	"github.com/gin-gonic/gin"
)

type Router interface {
	DefineRoutes(r *gin.Engine)
}

type routerImpl struct {
	authservice      authservice.Service
	messageservice   messageservice.Service
	webSocketService websocketservice.WebSocketService
}

func NewRouter(authservice authservice.Service, messageservice messageservice.Service, webSocketService websocketservice.WebSocketService) Router {
	return &routerImpl{
		authservice:      authservice,
		messageservice:   messageservice,
		webSocketService: webSocketService,
	}
}
