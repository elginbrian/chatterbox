package middleware

import (
	"chatterbox/internal/db"
	"context"

	"github.com/gofiber/fiber/v2"
)

func SessionMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	ctx := context.Background()
	userID, err := db.GetRedisClient().Get(ctx, token).Result()

	if err != nil || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid session"})
	}

	c.Locals("userID", userID)
	return c.Next()
}