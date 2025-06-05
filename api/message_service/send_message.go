package messageservice

import (
	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

func (s *serviceImpl) SendMessage(c *gin.Context, msg *models.Message) error {
	return s.mongodb.InsertMessage(c, msg)
}
