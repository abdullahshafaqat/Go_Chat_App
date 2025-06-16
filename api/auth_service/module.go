package authservice

import (
	db "github.com/abdullahshafaqat/Go_Chat_App.git/db/postgres"
	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

type Service interface {
	SignUp(c *gin.Context, user *models.UserSignup) error
	Login(c *gin.Context, login *models.UserLogin) error

}

type serviceImpl struct {
	database db.DB
}

func NewAuthService(db db.DB) Service {
	return &serviceImpl{database: db}
}
