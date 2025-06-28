package handlers

import (
	"github.com/avangero/auth-service/internal/lang"
	"github.com/avangero/auth-service/internal/middleware"
	"github.com/avangero/auth-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes настраивает маршруты приложения
func SetupRoutes(app *fiber.App, authService services.AuthService, messages lang.Messages) {
	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Создаем обработчик с зависимостями
	authHandler := NewAuthHandler(authService, messages)

	// Создаем JWT middleware
	jwtMiddleware := middleware.JWTMiddleware(authService, messages)

	// Публичные маршруты
	app.Get("/", authHandler.GetStatus)

	// API группа
	api := app.Group("/api/v1")

	// Публичные маршруты аутентификации
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	// Защищенные маршруты
	protected := api.Use(jwtMiddleware)
	protected.Get("/me", authHandler.GetMe)
	protected.Post("/validate", authHandler.ValidateToken)
}
