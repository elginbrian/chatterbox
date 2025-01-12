package repositories

import (
	"context"
	"fmt"

	"chatterbox/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository struct {
	DB *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{DB: db}
}

func (repo *ChatRepository) CreateChatRoom(ctx context.Context, chat *models.ChatRoom) error {
	query := "INSERT INTO chat_rooms (name, is_group) VALUES ($1, $2) RETURNING id"
	err := repo.DB.QueryRow(ctx, query, chat.Name, chat.IsGroup).Scan(&chat.ID)
	if err != nil {
		return fmt.Errorf("failed to create chat room: %v", err)
	}
	return nil
}

func (repo *ChatRepository) GetChatByID(ctx context.Context, id int) (*models.ChatRoom, error) {
	var chat models.ChatRoom
	query := "SELECT id, name, is_group, created_at FROM chat_rooms WHERE id=$1"
	err := repo.DB.QueryRow(ctx, query, id).Scan(&chat.ID, &chat.Name, &chat.IsGroup, &chat.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat room by ID: %v", err)
	}
	return &chat, nil
}
