package websocket

import "chatterbox/internal/utils"

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
			utils.LogInfo("New client connected: " + client.ID)

		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				client.Conn.Close()
				utils.LogInfo("Client disconnected: " + client.ID)
			}

		case msg := <-hub.broadcast:
			for client := range hub.clients {
				client.SendMessage(msg)
			}
		}
	}
}

func (hub *Hub) BroadcastMessage(client *Client, msg []byte) {
	utils.LogInfo("Broadcasting message from client: " + client.ID)
	hub.broadcast <- msg
}

func (h *Hub) Register(client *Client) {
    h.register <- client
}