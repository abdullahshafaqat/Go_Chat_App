package db

import (
	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type DB interface {
	CreateUser(c *gin.Context, user *models.UserSignup) error
	GetUserByEmail(email string) (string, string, error)
	GetID(email string) (string, error)
}

type dbImpl struct {
	db *sqlx.DB
}

func NewDB(db *sqlx.DB) DB {
	return &dbImpl{db: db}
}
