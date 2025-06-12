package mongodb

import (
	"errors"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (db *dbImpl) UpdateMessage(c *gin.Context, filter bson.M, update bson.M) (*models.Message, error) {
	result, err := db.collection.UpdateOne(c, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("message not found")
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("message found but not modified")
	}


	var updatedMsg models.Message
	err = db.collection.FindOne(c, filter).Decode(&updatedMsg)
	if err != nil {
		return nil, err
	}

	return &updatedMsg, nil
}
