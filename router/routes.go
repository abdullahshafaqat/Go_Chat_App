package router

import (
	"time"

	"github.com/abdullahshafaqat/Go_Chat_App.git/middelwares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (r *routerImpl) DefineRoutes(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"}, // React app origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/signup", r.SignUp)
	router.POST("/login", r.Login)
	router.POST("/refresh", r.RefreshToken)
	protected := router.Group("/")
	protected.Use(middelwares.AuthMiddleware()) 
	{
		protected.POST("/messages", r.SendMessage)
		protected.GET("/get_message", r.GetMessages)
		protected.GET("/update/:_id", r.UpdateMessage)
	}

}
