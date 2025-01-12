package api

import (
	"chatterbox/internal/models"
	"chatterbox/internal/services"
	"chatterbox/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{UserService: userService}
}

func (controller *UserController) CreateUser(c *fiber.Ctx) error {
	var userRequest models.UserRequest
	if err := c.BodyParser(&userRequest); err != nil {
		utils.LogError("Failed to parse request body: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if !utils.ValidateRequired(userRequest.Username) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Username is required")
	}
	if !utils.ValidateEmail(userRequest.Email) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid email address")
	}
	if !utils.ValidatePassword(userRequest.Password) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Password must be at least 8 characters")
	}

	user, err := controller.UserService.CreateUser(c.Context(), userRequest)
	if err != nil {
		utils.LogError("Failed to create user: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	utils.LogInfo("User created successfully: " + user.Username)
	return utils.SuccessResponse(c, user, "User created successfully")
}

func (controller *UserController) GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("id")

	if !utils.ValidateRequired(userID) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User ID is required")
	}

	user, err := controller.UserService.GetUserByID(c.Context(), userID)
	if err != nil {
		utils.LogError("Failed to retrieve user: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	utils.LogInfo("User retrieved successfully: " + user.Username)
	return utils.SuccessResponse(c, user, "User retrieved successfully")
}
