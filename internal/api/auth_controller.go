package api

import (
	"chatterbox/internal/auth"
	"chatterbox/internal/db"
	"chatterbox/internal/models"
	"chatterbox/internal/services"
	"chatterbox/internal/utils"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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

	user, err := controller.UserService.GetUserByUsername(c.Context(), loginRequest.Username)
	if err != nil {
		utils.LogError("Invalid credentials for user: " + loginRequest.Username)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)); err != nil {
		utils.LogError("Invalid password for user: " + loginRequest.Username)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid username or password")
	}

	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		utils.LogError("Failed to generate token: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	ctx := context.Background()
	if err := db.RedisClient.Set(ctx, token, user.ID, 24*time.Hour).Err(); err != nil {
		utils.LogError("Failed to store token in Redis: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to process login")
	}

	utils.LogInfo("User logged in successfully: " + user.Username)
	return utils.SuccessResponse(c, fiber.Map{"token": token}, "Login successful")
}

func (controller *AuthController) Register(c *fiber.Ctx) error {
	var registerRequest models.RegisterRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		utils.LogError("Failed to parse request body: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if !utils.ValidateRequired(registerRequest.Username) || !utils.ValidateRequired(registerRequest.Password) {
		utils.LogError("Missing username or password")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Username and password are required")
	}

	existingUser, _ := controller.UserService.GetUserByUsername(c.Context(), registerRequest.Username)
	if existingUser != nil {
		utils.LogError("User already exists: " + registerRequest.Username)
		return utils.ErrorResponse(c, fiber.StatusConflict, "Username is already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError("Failed to hash password: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to process registration")
	}

	user := &models.User{
		Username:     registerRequest.Username,
		PasswordHash: string(hashedPassword),
	}
	if err := controller.UserService.CreateUser(c.Context(), user); err != nil {
		utils.LogError("Failed to create user: " + err.Error())
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	utils.LogInfo("User registered successfully: " + user.Username)
	return utils.SuccessResponse(c, fiber.Map{"username": user.Username}, "Registration successful")
}