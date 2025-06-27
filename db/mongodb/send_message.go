package mongodb

import (
	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

func (db *dbImpl) InsertMessage(c *gin.Context, msg *models.Message) error {
	_, err := db.collection.InsertOne(c, msg)
	return err
}
