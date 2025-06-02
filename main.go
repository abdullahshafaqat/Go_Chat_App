package main

import (
	"log"

	authservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/auth_service"
	db "github.com/abdullahshafaqat/Go_Chat_App.git/db/postgres"
	"github.com/abdullahshafaqat/Go_Chat_App.git/router"
	"github.com/gin-gonic/gin"
)

func main() {
	database := db.Database()
	defer database.Close()
	dbLayer := db.NewDB(database)
	serviceLayer := authservice.NewAuthService(dbLayer)
	routerLayer := router.NewRouter(serviceLayer)

	r := gin.Default()
	routerLayer.DefineRoutes(r)

	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
