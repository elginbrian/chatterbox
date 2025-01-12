package routes

import (
	"chatterbox/internal/api"
	"chatterbox/pkg/websocket"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, hub *websocket.Hub) {
	websocketController := api.NewWebSocketController(hub)

	app.Get("/ws", websocketController.HandleWebSocket)
}
