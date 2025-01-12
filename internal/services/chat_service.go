package services

import (
	"context"

	"chatterbox/internal/models"
	"chatterbox/internal/repositories"
)

type ChatService struct {
	ChatRepo *repositories.ChatRepository
}

func NewChatService(chatRepo *repositories.ChatRepository) *ChatService {
	return &ChatService{ChatRepo: chatRepo}
}

func (service *ChatService) CreateChatRoom(ctx context.Context, chat *models.ChatRoom) error {
	return service.ChatRepo.CreateChatRoom(ctx, chat)
}

func (service *ChatService) GetChatByID(ctx context.Context, id int) (*models.ChatRoom, error) {
	return service.ChatRepo.GetChatByID(ctx, id)
}
