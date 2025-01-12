package middleware

import (
	"chatterbox/internal/auth"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing token")
	}

	token = strings.TrimPrefix(token, "Bearer ")

	claims, err := auth.ValidateToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
	}

	c.Locals("user", claims.Username)
	return c.Next()
}
