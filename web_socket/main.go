package main

import (
	"context"
	"log"

	messageservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/message_service"
	mongodb "github.com/abdullahshafaqat/Go_Chat_App.git/db/mongodb"
	wsrouter "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/router"
	websocketimpl "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/web_socket_impl"
)

func main() {
	mongoCollection := mongodb.ConnectToDb()
	defer func() {
		if err := mongoCollection.Database().Client().Disconnect(context.Background()); err != nil {
			log.Printf("MongoDB disconnect error: %v", err)
		}
	}()
	log.Println("Successfully connected to MongoDB - user-messages collection")

	wsService := websocketimpl.NewWebSocketService(mongodb.NewDB(mongoCollection))
	messageService := messageservice.NewMessageService(mongodb.NewDB(mongoCollection))

	websocketRouter := wsrouter.WSRouter(messageService, wsService)

	log.Println("WebSocket server running on :8004")
	if err := websocketRouter.Engine.Run(":8004"); err != nil {
		log.Fatal(err)
	}
}
 