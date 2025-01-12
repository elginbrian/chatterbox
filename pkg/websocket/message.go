package websocket

type Message struct {
	ClientID string `json:"client_id"`
	Content  string `json:"content"`
}

func NewMessage(clientID, content string) *Message {
	return &Message{
		ClientID: clientID,
		Content:  content,
	}
}