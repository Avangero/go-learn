package middleware

import (
	"log"
	"strings"

	"github.com/avangero/auth-service/internal/lang"
	"github.com/avangero/auth-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// JWTMiddleware создает middleware для проверки JWT токенов
func JWTMiddleware(authService services.AuthService, messages lang.Messages) fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientIP := c.IP()

		// Получаем заголовок Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Printf(messages.Get(lang.LogJWTMissingHeader), clientIP)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": messages.Get(lang.TokenNotProvided),
			})
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf(messages.Get(lang.LogJWTInvalidFormat), clientIP)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": messages.Get(lang.TokenInvalid),
			})
		}

		tokenString := parts[1]

		// Валидируем токен через AuthService
		user, err := authService.ValidateToken(c.Context(), tokenString)
		if err != nil {
			log.Printf(messages.Get(lang.LogJWTValidationFailed), clientIP, err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": messages.Get(lang.TokenInvalid),
			})
		}

		// Сохраняем пользователя в контексте для использования в handlers
		c.Locals("user", user)
		log.Printf(messages.Get(lang.LogJWTValidationSuccess), clientIP, user.Email)

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
