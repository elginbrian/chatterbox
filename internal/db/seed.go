package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedDatabase(dbPool *pgxpool.Pool) error {
	queries := []string{
		`INSERT INTO users (username, email, password_hash) VALUES ('john_doe', 'john@example.com', 'hashed_password') ON CONFLICT DO NOTHING`,
		`INSERT INTO chat_rooms (name, is_group) VALUES ('General', TRUE) ON CONFLICT DO NOTHING`,
	}

	for _, query := range queries {
		if _, err := dbPool.Exec(context.Background(), query); err != nil {
			log.Printf("Failed to execute seed query: %v", err)
			return err
		}
	}

	log.Println("Database seeded successfully.")
	return nil
}
