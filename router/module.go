package router

import (
	authservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/auth_service"
	"github.com/gin-gonic/gin"
)

type Router interface {
	DefineRoutes(r *gin.Engine)
}

type routerImpl struct {
	authservice authservice.Service
}

func NewRouter(authservice authservice.Service) Router {
	return &routerImpl{

		authservice: authservice,
	}
}
