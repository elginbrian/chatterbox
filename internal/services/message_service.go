package services

import (
	"context"

	"chatterbox/internal/models"
	"chatterbox/internal/repositories"
)

type MessageService struct {
	MessageRepo *repositories.MessageRepository
}

func NewMessageService(messageRepo *repositories.MessageRepository) *MessageService {
	return &MessageService{MessageRepo: messageRepo}
}

func (service *MessageService) CreateMessage(ctx context.Context, message *models.Message) error {
	return service.MessageRepo.CreateMessage(ctx, message)
}

func (service *MessageService) GetMessagesByChatID(ctx context.Context, chatID int) ([]models.Message, error) {
	return service.MessageRepo.GetMessagesByChatID(ctx, chatID)
}
