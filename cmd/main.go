package main

import (
	"context"
	"log"

	"chatterbox/internal/api"
	"chatterbox/internal/middleware"
	"chatterbox/internal/repositories"
	"chatterbox/internal/services"
	"chatterbox/internal/utils"
	"chatterbox/pkg/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	db, err := pgxpool.Connect(context.Background(), "postgres://username:password@localhost:5432/chatterbox")
	if err != nil {
		utils.LogError("Unable to connect to database: " + err.Error())
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	userService := services.NewUserService(userRepo)
	chatService := services.NewChatService(chatRepo)
	messageService := services.NewMessageService(messageRepo)

	userController := api.NewUserController(userService)
	chatController := api.NewChatController(chatService)
	messageController := api.NewMessageController(messageService)
	authController := api.NewAuthController(userService)

	hub := websocket.NewHub()

	go hub.Run()

	app := fiber.New()

	app.Post("/login", authController.Login)

	protected := app.Group("/api", middleware.AuthMiddleware)
	protected.Post("/users", userController.CreateUser)
	protected.Get("/users/:id", userController.GetUserByID)
	protected.Post("/chatrooms", chatController.CreateChatRoom)
	protected.Get("/chatrooms/:id", chatController.GetChatByID)
	protected.Post("/messages", messageController.CreateMessage)
	protected.Get("/messages/:chat_id", messageController.GetMessagesByChatID)

	app.Get("/ws", api.NewWebSocketController(hub).HandleWebSocket)

	if err := app.Listen(":8080"); err != nil {
		utils.LogError("Error starting server: " + err.Error())
		log.Fatalf("Error starting server: %v", err)
	}
}