package wsrouter

import (
    "time"

    messageservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/message_service"
    websocketimpl "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/web_socket_impl"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

type Router struct {
    Engine          *gin.Engine
    MessageService  messageservice.Service
    WebSocketService websocketimpl.WebSocketService
}

func WSRouter(MessageService messageservice.Service, WebSocketService websocketimpl.WebSocketService) *Router {
    engine := gin.Default()

    engine.Use(cors.New(cors.Config{
        AllowAllOrigins:  true,
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    router := &Router{
        Engine:          engine,
        MessageService:  MessageService,
        WebSocketService: WebSocketService,
    }

    router.defineWebSocketRouter()
    return router
}