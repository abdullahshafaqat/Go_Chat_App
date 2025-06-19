package messageservice

import (
	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

func (s *serviceImpl) GetMessages(c *gin.Context, userID int) ([]models.Message, error) {
	return s.mongodb.GetMessages(c, userID)
}
