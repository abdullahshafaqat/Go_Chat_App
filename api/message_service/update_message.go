package messageservice

import (
	"errors"
	"time"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *serviceImpl) UpdateMessage(c *gin.Context, messageID string, senderID int, NewMessage string) (*models.Message, error) {

	objID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return nil, errors.New("invalid message ID format")
	}

	filter := bson.M{
		"_id":       objID,
		"sender_id": senderID,
	}

	update := bson.M{
		"$set": bson.M{
			"message":   NewMessage,
			"edited":    true,
			"edited_at": time.Now(),
		},
	}

	updatedMsg, err := s.mongodb.UpdateMessage(c, filter, update)
	if err != nil {
		return nil, err
	}

	if updatedMsg == nil {
		return nil, errors.New("message not found or not owned by user")
	}

	return updatedMsg, nil
}
