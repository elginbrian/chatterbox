package utils

import (
	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(c *fiber.Ctx, data interface{}, message string) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  "error",
		"message": message,
	})
}
