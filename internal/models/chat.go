package models

import "time"

type ChatRoom struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IsGroup   bool      `json:"is_group"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatRequest struct {
    Name string `json:"name"`
}