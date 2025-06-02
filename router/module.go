package router

import (
	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

type Router interface {
	DefineRoutes(r *gin.Engine)
}

type routerImpl struct {
	service Service
}

type Service interface {
	SignUp(c *gin.Context, user *models.UserSignup) error
	Login(c *gin.Context, login *models.UserLogin) error
}

func NewRouter(service Service) Router {
	return &routerImpl{service: service}
}
