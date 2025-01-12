package websocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)


func HandleWebSocket(handler func(*websocket.Conn)) func(*fiber.Ctx) error {

    return websocket.New(func(c *websocket.Conn) {

        handler(c)

    })

}
