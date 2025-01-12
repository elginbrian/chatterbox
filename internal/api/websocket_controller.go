package api

import (
	websocketPkg "chatterbox/pkg/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type WebSocketController struct {
	Hub *websocketPkg.Hub
}

func NewWebSocketController(hub *websocketPkg.Hub) *WebSocketController {
	return &WebSocketController{Hub: hub}
}

func (controller *WebSocketController) HandleWebSocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return websocketPkg.HandleWebSocket(func(conn *websocket.Conn) {
			client := websocketPkg.NewClient(conn, controller.Hub)
			controller.Hub.Register(client)

			go client.ReadMessages()

			select {}
		})(c)
	}
	return fiber.ErrUpgradeRequired
}
