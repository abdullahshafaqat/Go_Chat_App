package db

import (
	"github.com/abdullahshafaqat/Go_ChatApp.git/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type DBStorage interface {
	SignUp(c *gin.Context, req *models.UserSignup) *models.UserSignup
}

type StorageImpl struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) DBStorage {
	return &StorageImpl{
		db: db,
	}

}
