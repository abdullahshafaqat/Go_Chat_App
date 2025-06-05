package mongodb

import (
	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database interface {
	InsertMessage(c *gin.Context, msg *models.Message) error
}

type dbImpl struct {
	collection *mongo.Collection
}

func NewDB(collection *mongo.Collection) Database {
	return &dbImpl{collection: collection}
}
