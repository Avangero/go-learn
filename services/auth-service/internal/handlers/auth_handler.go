package handlers

import (
	"github.com/avangero/auth-service/internal/lang"
	"github.com/avangero/auth-service/internal/models"
	"github.com/avangero/auth-service/internal/models/requests"
	"github.com/avangero/auth-service/internal/models/responses"
	"github.com/avangero/auth-service/internal/services"
	"github.com/avangero/auth-service/internal/validators"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler обработчик для аутентификации
type AuthHandler struct {
	authService services.AuthService
	validator   *validators.AuthValidator
	messages    lang.Messages
}

// NewAuthHandler создает новый обработчик аутентификации
func NewAuthHandler(authService services.AuthService, messages lang.Messages) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validators.NewAuthValidator(messages),
		messages:    messages,
	}
}

// GetStatus возвращает статус сервиса
func (h *AuthHandler) GetStatus(c *fiber.Ctx) error {
	return c.JSON(responses.StatusResponse{
		Service: "Auth Service",
		Status:  "running",
		Version: "1.0.0",
	})
}

// Register регистрирует нового пользователя
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req requests.RegisterRequest

	// Парсим JSON запрос
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error: h.messages.Get(lang.InvalidRequestFormat),
		})
	}

	// Валидируем данные
	if err := h.validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error: err.Error(),
		})
	}

	// Регистрируем пользователя
	tokenResponse, err := h.authService.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(responses.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(tokenResponse)
}

// Login аутентифицирует пользователя
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req requests.LoginRequest

	// Парсим JSON запрос
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error: h.messages.Get(lang.InvalidRequestFormat),
		})
	}

	// Валидируем данные
	if err := h.validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error: err.Error(),
		})
	}

	// Аутентифицируем пользователя
	tokenResponse, err := h.authService.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(tokenResponse)
}

// GetMe возвращает информацию о текущем пользователе
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	// Получаем пользователя из контекста (установлен в middleware)
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.ErrorResponse{
			Error: h.messages.Get(lang.TokenInvalid),
		})
	}

	return c.JSON(user)
}

// ValidateToken валидирует JWT токен
func (h *AuthHandler) ValidateToken(c *fiber.Ctx) error {
	// Получаем пользователя из контекста (установлен в middleware)
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.ErrorResponse{
			Error: h.messages.Get(lang.TokenInvalid),
		})
	}

	return c.JSON(responses.ValidationResponse{
		Valid: true,
		User:  *user,
	})
}
