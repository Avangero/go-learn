package main

import (
	"log"

	"github.com/avangero/auth-service/internal/config"
	"github.com/avangero/auth-service/internal/database"
	"github.com/avangero/auth-service/internal/handlers"
	"github.com/avangero/auth-service/internal/lang/ru"
	"github.com/avangero/auth-service/internal/repositories"
	"github.com/avangero/auth-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Инициализация системы сообщений
	messages := ru.NewRussianMessages()

	// Загрузка конфигурации
	configLoader := config.NewLoader(messages)
	cfg, err := configLoader.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	// Подключение к базе данных
	connectionManager := database.NewConnectionManager(messages)
	db := connectionManager.Connect(cfg)

	// Инициализация слоев приложения (Dependency Injection)
	userRepo := repositories.NewUserRepository(db, messages)
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret, cfg.BCryptCost, messages)

	// Создание Fiber приложения
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Настройка маршрутов
	handlers.SetupRoutes(app, authService, messages)

	// Запуск сервера
	log.Printf("🚀 Сервер запущен на порту %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
