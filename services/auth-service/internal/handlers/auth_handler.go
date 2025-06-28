package handlers

import (
	"log"

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
	clientIP := c.IP()
	log.Printf(h.messages.Get(lang.LogRegistrationRequest), clientIP)

	var req requests.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf(h.messages.Get(lang.LogParseRequestFailed), clientIP, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": h.messages.Get(lang.InvalidRequestFormat),
		})
	}

	// Валидация
	if err := h.validator.Validate(&req); err != nil {
		log.Printf(h.messages.Get(lang.LogValidationFailed), clientIP, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   h.messages.Get(lang.InvalidRequestFormat),
			"details": err.Error(),
		})
	}

	// Регистрация
	response, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		log.Printf(h.messages.Get(lang.LogRegistrationFailed), clientIP, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Printf(h.messages.Get(lang.LogRegistrationSuccess), clientIP, req.Email)
	return c.JSON(response)
}

// Login аутентифицирует пользователя
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	clientIP := c.IP()
	log.Printf(h.messages.Get(lang.LogLoginRequest), clientIP)

	var req requests.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf(h.messages.Get(lang.LogParseRequestFailed), clientIP, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": h.messages.Get(lang.InvalidRequestFormat),
		})
	}

	// Валидация
	if err := h.validator.Validate(&req); err != nil {
		log.Printf(h.messages.Get(lang.LogValidationFailed), clientIP, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   h.messages.Get(lang.InvalidRequestFormat),
			"details": err.Error(),
		})
	}

	// Аутентификация
	response, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		log.Printf(h.messages.Get(lang.LogLoginFailed), clientIP, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Printf(h.messages.Get(lang.LogLoginSuccess), clientIP, req.Email)
	return c.JSON(response)
}

// GetMe возвращает информацию о текущем пользователе
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	clientIP := c.IP()
	user := c.Locals("user")
	if user == nil {
		log.Printf(h.messages.Get(lang.LogGetMeFailed), clientIP)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": h.messages.Get(lang.UserNotFound),
		})
	}

	log.Printf(h.messages.Get(lang.LogGetMeSuccess), clientIP, user.(*models.User).Email)
	return c.JSON(user)
}

// ValidateToken валидирует JWT токен
func (h *AuthHandler) ValidateToken(c *fiber.Ctx) error {
	clientIP := c.IP()

	// Получаем пользователя из контекста (установлен в middleware)
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		log.Printf("Token validation failed: user not found in context from IP %s", clientIP)
		return c.Status(fiber.StatusUnauthorized).JSON(responses.ErrorResponse{
			Error: h.messages.Get(lang.TokenInvalid),
		})
	}

	log.Printf("Token validation successful for IP %s, user: %s", clientIP, user.Email)
	return c.JSON(responses.ValidationResponse{
		Valid: true,
		User:  *user,
	})
}
