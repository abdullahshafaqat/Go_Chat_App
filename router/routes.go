package router

import "github.com/gin-gonic/gin"

func (r *routerImpl) DefineRoutes(router *gin.Engine) {
	router.POST("/signup", r.SignUp)
	router.POST("/login", r.Login)
	router.POST("/auth", r.Authorize)
	router.POST("/refresh", r.RefreshToken)
	router.POST("/messages", r.SendMessage)
	router.GET("/get", r.GetMessages)
}
