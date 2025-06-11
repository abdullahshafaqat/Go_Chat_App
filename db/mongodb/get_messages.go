package mongodb

import (
	"context"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *dbImpl) GetMessages(c *gin.Context, senderID int) ([]models.Message, error) {
	filter := bson.M{"sender_id": senderID}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := db.collection.Find(c, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var messages []models.Message
	if err = cursor.All(c, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
