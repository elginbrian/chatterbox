package api

import (
	"chatterbox/internal/auth"
	"chatterbox/internal/models"
	"chatterbox/internal/services"
	"chatterbox/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	UserService *services.UserService
}

func NewAuthController(userService *services.UserService) *AuthController {
	return &AuthController{UserService: userService}
}

func (controller *AuthController) Login(c *fiber.Ctx) error {
	var loginRequest models.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		utils.LogError("Failed to parse request body: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if !utils.ValidateRequired(loginRequest.Username) || !utils.ValidateRequired(loginRequest.Password) {
		utils.LogError("Missing username or password")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Username and password are required")
	}

	user, err := controller.UserService.GetUserByID(c.Context(), loginRequest.Username)
	if err != nil {
		utils.LogError("Invalid credentials for user: " + loginRequest.Username)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid username or password")
	}

	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		utils.LogError("Failed to generate token: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	utils.LogInfo("User logged in successfully: " + user.Username)
	return utils.SuccessResponse(c, fiber.Map{"token": token}, "Login successful")
}
