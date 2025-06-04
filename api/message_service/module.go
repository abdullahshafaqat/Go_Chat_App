package messageservice

import (
	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
)

type Service interface {
	SendMessage(message *models.Message) error
}

