package authservice

import (
	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

type Service interface {
	SignUp(c *gin.Context, user *models.UserSignup) error
	Login(c *gin.Context, login *models.UserLogin) error
}

type serviceImpl struct {
	db interface {
		CreateUser(c *gin.Context, user *models.UserSignup) error
		GetUserByEmail(c *gin.Context, email string) (*models.UserLogin, error)
	}
}

func NewAuthService(db interface {
	CreateUser(c *gin.Context, user *models.UserSignup) error
	GetUserByEmail(c *gin.Context, email string) (*models.UserLogin, error)
}) Service {
	return &serviceImpl{db: db}
}
