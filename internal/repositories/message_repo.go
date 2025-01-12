package repositories

import (
	"context"
	"fmt"

	"chatterbox/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	DB *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{DB: db}
}

func (repo *MessageRepository) CreateMessage(ctx context.Context, message *models.Message) error {
	query := "INSERT INTO messages (chat_id, sender_id, content) VALUES ($1, $2, $3) RETURNING id"
	err := repo.DB.QueryRow(ctx, query, message.ChatID, message.SenderID, message.Content).Scan(&message.ID)
	if err != nil {
		return fmt.Errorf("failed to create message: %v", err)
	}
	return nil
}

func (repo *MessageRepository) GetMessagesByChatID(ctx context.Context, chatID int) ([]models.Message, error) {
	var messages []models.Message
	query := "SELECT id, chat_id, sender_id, content, created_at, is_read FROM messages WHERE chat_id=$1"
	rows, err := repo.DB.Query(ctx, query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by chat ID: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Content, &message.CreatedAt, &message.IsRead); err != nil {
			return nil, fmt.Errorf("failed to scan message: %v", err)
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over messages: %v", err)
	}

	return messages, nil
}
