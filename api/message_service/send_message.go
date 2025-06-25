package messageservice

import (
	"context"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
)

func (s *serviceImpl) SendMessage(ctx context.Context, msg *models.Message) error {
	return s.mongodb.InsertMessage(ctx, msg)
}
