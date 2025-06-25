package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	messageservice "github.com/abdullahshafaqat/Go_Chat_App.git/api/message_service"
	mongodb "github.com/abdullahshafaqat/Go_Chat_App.git/db/mongodb"
	wsrouter "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/router"
	websocketimpl "github.com/abdullahshafaqat/Go_Chat_App.git/web_socket/web_socket_impl"
)

func main() {

	mongoCollection := mongodb.ConnectToDb()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoCollection.Database().Client().Disconnect(ctx); err != nil {
			log.Printf("MongoDB disconnect error: %v", err)
		}
	}()

	log.Println("Successfully connected to MongoDB")

	wsService := websocketimpl.NewWebSocketService(mongodb.NewDB(mongoCollection))
	messageService := messageservice.NewMessageService(mongodb.NewDB(mongoCollection))

	websocketRouter := wsrouter.WSRouter(messageService, wsService)

	server := &http.Server{
		Addr:    ":8004",
		Handler: websocketRouter.Engine,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("WebSocket server running on :8004")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server shutdown completed")
}