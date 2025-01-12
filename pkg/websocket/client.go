package websocket

import (
	"chatterbox/internal/utils"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Hub  *Hub
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:   conn.RemoteAddr().String(),
		Conn: conn,
		Hub:  hub,
	}
}

func (client *Client) ReadMessages() {
	defer func() {
		client.Hub.unregister <- client
		client.Conn.Close()
	}()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			utils.LogError("Failed to read message from client: " + err.Error())
			break
		}

		client.Hub.BroadcastMessage(client, msg)
	}
}

func (client *Client) SendMessage(msg []byte) {
	if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		utils.LogError("Failed to send message to client: " + err.Error())
	}
}
