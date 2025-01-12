package repositories

import (
	"context"
	"fmt"

	"chatterbox/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password_hash, created_at FROM users WHERE id=$1"
	err := repo.DB.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %v", err)
	}
	return &user, nil
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id"
	err := repo.DB.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}
