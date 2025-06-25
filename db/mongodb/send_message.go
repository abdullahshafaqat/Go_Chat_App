package mongodb

import (
	"context"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
)

func (db *dbImpl) InsertMessage(ctx context.Context, msg *models.Message) error {
	_, err := db.collection.InsertOne(ctx, msg)
	return err
}
