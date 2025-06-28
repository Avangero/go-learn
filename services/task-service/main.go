package main

import (
	// Для форматирования сообщений об ошибках
	"fmt"
	"log" // Для простого логирования на старте
	"os"  // Для чтения переменных окружения

	"github.com/go-playground/validator/v10" // Для валидации
	"github.com/gofiber/fiber/v2"            // Импортируем Fiber
	"github.com/google/uuid"
)

// Task представляет модель задачи
type Task struct {
	ID          string `json:"id"`                                                     // `json:"id"` - это тег, который указывает, как поле должно быть названо в JSON
	Title       string `json:"title" validate:"required"`                              // Title обязательно для заполнения
	Description string `json:"description"`                                            // Description не обязательно
	Status      string `json:"status" validate:"required,oneof=TODO IN_PROGRESS DONE"` // Status обязательно и должен быть одним из разрешенных значений
}

// Глобальный валидатор
var validate = validator.New()

// Временное хранилище задач в памяти (пока без БД)
var tasks = map[string]Task{} // Карта для быстрого поиска задач по ID

// validateTask валидирует структуру Task и возвращает читаемые сообщения об ошибках
func validateTask(task *Task) error {
	err := validate.Struct(task)
	if err != nil {
		// Формируем понятные сообщения об ошибках
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("Поле '%s' обязательно для заполнения", err.Field()))
			case "oneof":
				errorMessages = append(errorMessages, fmt.Sprintf("Поле '%s' должно быть одним из: TODO, IN_PROGRESS, DONE", err.Field()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("Поле '%s' не прошло валидацию: %s", err.Field(), err.Tag()))
			}
		}
		return fmt.Errorf("Ошибки валидации: %v", errorMessages)
	}
	return nil
}

func main() {
	// Инициализация Fiber приложения
	app := fiber.New()

	// 1. Приветственный маршрут
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Привет, я сервис задач!")
	})

	// 2. Создание задачи (POST /tasks)
	app.Post("/tasks", createTask)

	// 3. Получение всех задач (GET /tasks)
	app.Get("/tasks", getTasks)

	// 4. Получение задачи по ID (GET /tasks/:id)
	app.Get("/tasks/:id", getTaskByID)

	// 5. Обновление задачи (PUT /tasks/:id)
	app.Put("/tasks/:id", updateTask)

	// 6. Удаление задачи (DELETE /tasks/:id)
	app.Delete("/tasks/:id", deleteTask)

	// Читаем порт из переменной окружения, по умолчанию 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Task Service запущен на порту %s", port)
	log.Fatal(app.Listen(":" + port)) // Запускаем сервер
}

// createTask - обработчик для создания новой задачи
func createTask(c *fiber.Ctx) error {
	task := new(Task) // Создаем указатель на новую структуру Task

	// Парсим тело запроса в структуру Task.
	// Fiber автоматически обрабатывает JSON.
	if err := c.BodyParser(task); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса",
		})
	}

	// Валидируем полученные данные
	if err := validateTask(task); err != nil {
		log.Printf("Ошибка валидации: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Генерация простого ID (в реальном приложении это будет UUID)
	task.ID = uuid.New().String()
	tasks[task.ID] = *task // Сохраняем задачу в карте

	log.Printf("Задача создана: %+v", task) // %+v выводит структуру со всеми полями
	return c.Status(fiber.StatusCreated).JSON(task)
}

// getTasks - обработчик для получения всех задач
func getTasks(c *fiber.Ctx) error {
	// Преобразуем карту задач в срез задач для вывода
	allTasks := []Task{}
	for _, task := range tasks {
		allTasks = append(allTasks, task)
	}
	return c.JSON(allTasks)
}

// getTaskByID - обработчик для получения задачи по ID
func getTaskByID(c *fiber.Ctx) error {
	id := c.Params("id") // Получаем ID из параметров маршрута

	task, ok := tasks[id] // Ищем задачу в карте
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Задача не найдена",
		})
	}
	return c.JSON(task)
}

// updateTask - обработчик для обновления задачи
func updateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	_, ok := tasks[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Задача не найдена",
		})
	}

	updatedTask := new(Task)
	if err := c.BodyParser(updatedTask); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса",
		})
	}

	// Обновляем только те поля, которые пришли в запросе
	// В реальном приложении нужно быть осторожным с этим, чтобы не перезаписать важные поля
	existingTask := tasks[id]
	if updatedTask.Title != "" {
		existingTask.Title = updatedTask.Title
	}
	if updatedTask.Description != "" {
		existingTask.Description = updatedTask.Description
	}
	if updatedTask.Status != "" {
		existingTask.Status = updatedTask.Status
	}

	// Валидируем обновленную задачу
	if err := validateTask(&existingTask); err != nil {
		log.Printf("Ошибка валидации при обновлении: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	tasks[id] = existingTask // Сохраняем обновленную задачу

	log.Printf("Задача обновлена: %+v", existingTask)
	return c.JSON(existingTask)
}

// deleteTask - обработчик для удаления задачи
func deleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	_, ok := tasks[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Задача не найдена",
		})
	}

	delete(tasks, id) // Удаляем задачу из карты
	log.Printf("Задача удалена: %s", id)
	return c.SendStatus(fiber.StatusNoContent) // 204 No Content
}
