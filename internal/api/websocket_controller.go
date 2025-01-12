package api

import (
	"chatterbox/pkg/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type WebSocketController struct {
	Hub *websocket.Hub
}

func NewWebSocketController(hub *websocket.Hub) *WebSocketController {
	return &WebSocketController{Hub: hub}
}

func (controller *WebSocketController) HandleWebSocket(c *fiber.Ctx) error {
	conn, err := websocket.Upgrade(c)
	if err != nil {
		return err
	}

	client := websocket.NewClient(conn, controller.Hub)
	controller.Hub.register <- client

	go client.ReadMessages()

	select {}
}
