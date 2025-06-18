package main

import (
	"context"
	"log"

	authservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/auth_service"
	messageservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/message_service"
	websocketservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/web_socket_service"
	mongodb "github.com/abdullahshafaqat/Go_Chat_App.git/db/mongodb" // Add this import
	db "github.com/abdullahshafaqat/Go_Chat_App.git/db/postgres"
	"github.com/abdullahshafaqat/Go_Chat_App.git/router"
	"github.com/gin-gonic/gin"
)

func main() {
	database := db.Database()
	defer database.Close()
	dbLayer := db.NewDB(database)

	mongoCollection := mongodb.ConnectToDb()
	defer func() {
		if err := mongoCollection.Database().Client().Disconnect(context.Background()); err != nil {
			log.Printf("MongoDB disconnect error: %v", err)
		}
	}()
	log.Println("Successfully connected to MongoDB - user-messages collection")


	serviceLayer := authservice.NewAuthService(dbLayer)
	routerLayer := router.NewRouter(serviceLayer, messageservice.NewMessageService(mongodb.NewDB(mongoCollection)))

	r := gin.Default()
	routerLayer.DefineRoutes(r)

	log.Println("Server running on :8003")
	if err := r.Run(":8003"); err != nil {
		log.Fatal(err)
	}
}
