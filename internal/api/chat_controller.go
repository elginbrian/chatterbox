package api

import (
	"chatterbox/internal/models"
	"chatterbox/internal/services"
	"chatterbox/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ChatController struct {
	ChatService *services.ChatService
}

func NewChatController(chatService *services.ChatService) *ChatController {
	return &ChatController{ChatService: chatService}
}

func (controller *ChatController) CreateChatRoom(c *fiber.Ctx) error {
	var chatRequest models.ChatRequest
	if err := c.BodyParser(&chatRequest); err != nil {
		utils.LogError("Failed to parse request body: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if !utils.ValidateRequired(chatRequest.Name) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Chat room name is required")
	}

	chatRoom := &models.ChatRoom{}
	chatRoom.Name = chatRequest.Name
	err := controller.ChatService.CreateChatRoom(c.Context(), chatRoom)
	if err != nil {
		utils.LogError("Failed to create chat room: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create chat room")
	}

	utils.LogInfo("Chat room created successfully: " + chatRoom.Name)
	return utils.SuccessResponse(c, chatRoom, "Chat room created successfully")
}

func (controller *ChatController) GetChatByID(c *fiber.Ctx) error {
	chatID := c.Params("id")

	if !utils.ValidateRequired(chatID) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Chat ID is required")
	}

	chatIDInt, err := strconv.Atoi(chatID)
	if err != nil {
		utils.LogError("Invalid chat ID: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid chat ID")
	}

	chatRoom, err := controller.ChatService.GetChatByID(c.Context(), chatIDInt)
	if err != nil {
		utils.LogError("Failed to retrieve chat room: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Chat room not found")
	}

	utils.LogInfo("Chat room retrieved successfully: " + chatRoom.Name)
	return utils.SuccessResponse(c, chatRoom, "Chat room retrieved successfully")
}