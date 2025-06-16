package router

import (
	"time"

	messageservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/message_service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSockets interface {
	AddConn(userID string, conn *websocket.Conn, c *gin.Context) error
}

type Router struct {
	Engine       *gin.Engine
	WebSocketSvc WebSockets
}

func NewRouter(messageService messageservice.Service, websocket WebSockets) *Router {
	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r := &Router{
		Engine:       engine,
		WebSocketSvc: websocket,
	}

	r.defineWebSocketRouter()

	return r
}
