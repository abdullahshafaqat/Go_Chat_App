package main

import (
	"context"
	"log"

	authservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/auth_service"
	mongodb "github.com/abdullahshafaqat/Go_Chat_App.git/db/mongodb" // Add this import
	db "github.com/abdullahshafaqat/Go_Chat_App.git/db/postgres"
	"github.com/abdullahshafaqat/Go_Chat_App.git/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize PostgreSQL
	database := db.Database()
	defer database.Close()
	dbLayer := db.NewDB(database)

	// Test MongoDB connection
	mongoCollection := mongodb.ConnectToUserMessages()
	defer func() {
		if err := mongoCollection.Database().Client().Disconnect(context.Background()); err != nil {
			log.Printf("MongoDB disconnect error: %v", err)
		}
	}()
	log.Println("Successfully connected to MongoDB - user-messages collection")

	// Continue with your service and router setup
	serviceLayer := authservice.NewAuthService(dbLayer)
	routerLayer := router.NewRouter(serviceLayer)

	r := gin.Default()
	routerLayer.DefineRoutes(r)

	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
