package messageservice

import (
	"github.com/abdullahshafaqat/Go_Chat_App.git/db/mongodb"
	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

type Service interface {
	SendMessage(c *gin.Context, msg *models.Message) error
}

type serviceImpl struct {
	mongodb mongodb.Database
}

func NewMessageService(mongodb mongodb.Database) Service {
	return &serviceImpl{mongodb: mongodb}
}
