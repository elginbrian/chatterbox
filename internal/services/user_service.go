package services

import (
	"context"

	"chatterbox/internal/models"
	"chatterbox/internal/repositories"
)

type UserService struct {
	UserRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (service *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return service.UserRepo.CreateUser(ctx, user)
}

func (service *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return service.UserRepo.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
    return s.UserRepo.GetUserByUsername(ctx, username)
}