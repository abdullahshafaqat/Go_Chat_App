package mongodb

import (
	"context"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database interface {
	InsertMessage(ctx context.Context, msg *models.Message) error
	GetMessages(c *gin.Context, senderID int) ([]models.Message, error)
	UpdateMessage(c *gin.Context, filter bson.M, update bson.M) (*models.Message, error)
}

type dbImpl struct {
	collection *mongo.Collection
}

func NewDB(collection *mongo.Collection) Database {
	return &dbImpl{collection: collection}
}
