package api

import (
	"chatterbox/internal/db"
	"chatterbox/internal/models"
	"chatterbox/internal/services"
	"chatterbox/internal/utils"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MessageController struct {
	MessageService *services.MessageService
}

func NewMessageController(messageService *services.MessageService) *MessageController {
	return &MessageController{MessageService: messageService}
}

func (controller *MessageController) CreateMessage(c *fiber.Ctx) error {
	var messageRequest models.MessageRequest
	if err := c.BodyParser(&messageRequest); err != nil {
		utils.LogError("Failed to parse request body: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if !utils.ValidateRequired(messageRequest.Content) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Message content is required")
	}

	message := &models.Message{
		Content:  messageRequest.Content,
		ChatID:   messageRequest.ChatID,
		SenderID: messageRequest.SenderID,
	}
	if err := controller.MessageService.CreateMessage(c.Context(), message); err != nil {
		utils.LogError("Failed to create message: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create message")
	}

	ctx := context.Background()
	cacheKey := "messages:chat:" + strconv.Itoa(message.ChatID)
	if err := db.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		utils.LogError("Failed to invalidate cache for chat room: " + strconv.Itoa(message.ChatID))
	}

	utils.LogInfo("Message created successfully in chat room: " + strconv.Itoa(message.ChatID))
	return utils.SuccessResponse(c, message, "Message created successfully")
}

func (controller *MessageController) GetMessagesByChatID(c *fiber.Ctx) error {
	chatID := c.Params("chat_id")

	if !utils.ValidateRequired(chatID) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Chat ID is required")
	}

	chatIDInt, err := strconv.Atoi(chatID)
	if err != nil {
		utils.LogError("Invalid chat ID: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid chat ID")
	}

	ctx := context.Background()
	cacheKey := "messages:chat:" + chatID
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {

		var cachedMessages []models.Message
		if err := json.Unmarshal([]byte(cachedData), &cachedMessages); err == nil {
			utils.LogInfo("Messages retrieved from cache for chat room: " + chatID)
			return utils.SuccessResponse(c, cachedMessages, "Messages retrieved successfully")
		}
	}

	messages, err := controller.MessageService.GetMessagesByChatID(c.Context(), chatIDInt)
	if err != nil {
		utils.LogError("Failed to retrieve messages: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Messages not found")
	}

	messagesJSON, err := json.Marshal(messages)
	if err == nil {
		err = db.RedisClient.Set(ctx, cacheKey, messagesJSON, 10*time.Minute).Err()
		if err != nil {
			utils.LogError("Failed to cache messages for chat room: " + chatID)
		}
	}

	utils.LogInfo("Messages retrieved successfully for chat room: " + chatID)
	return utils.SuccessResponse(c, messages, "Messages retrieved successfully")
}