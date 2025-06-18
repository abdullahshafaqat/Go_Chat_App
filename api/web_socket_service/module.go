package websocketservice

import "github.com/gin-gonic/gin"

type Service interface {
	HandleConnections(c *gin.Context)
}

type serviceImpl struct {
	manager *websocket.Manager
}

func NewService(manager *websocket.Manager) Service {
	return &serviceImpl{manager: manager}
}
