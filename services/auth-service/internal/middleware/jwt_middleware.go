package middleware

import (
	"strings"

	"github.com/avangero/auth-service/internal/lang"
	"github.com/avangero/auth-service/internal/models/responses"
	"github.com/avangero/auth-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// JWTMiddleware создает middleware для проверки JWT токенов
func JWTMiddleware(authService services.AuthService, messages lang.Messages) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем токен из заголовка Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.ErrorResponse{
				Error: messages.Get(lang.TokenNotProvided),
			})
		}

		// Проверяем формат "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.ErrorResponse{
				Error: messages.Get(lang.TokenInvalid),
			})
		}

		// Валидируем токен
		user, err := authService.ValidateToken(tokenParts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.ErrorResponse{
				Error: messages.Get(lang.TokenInvalid),
			})
		}

		// Сохраняем пользователя в контексте
		c.Locals("user", user)
		return c.Next()
	}
}

// GetUserID извлекает ID пользователя из контекста
func GetUserID(c *fiber.Ctx) uuid.UUID {
	return c.Locals("user_id").(uuid.UUID)
}

// GetUserEmail извлекает email пользователя из контекста
func GetUserEmail(c *fiber.Ctx) string {
	return c.Locals("user_email").(string)
}

// GetUserRole извлекает роль пользователя из контекста
func GetUserRole(c *fiber.Ctx) string {
	return c.Locals("user_role").(string)
}
