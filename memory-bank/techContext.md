# Технический контекст

## Используемые технологии

### Основной стек
- **Go**: 1.23.4 (основной язык)
- **Fiber**: v2.52.8 (веб-фреймворк)
- **Validator**: v10.22.1 (валидация данных)
- **UUID**: v1.6.0 (генерация уникальных ID)

### Инфраструктура
- **Docker**: Контейнеризация сервисов
- **Docker Compose**: Оркестрация для разработки
- **Air**: Live-reloading для разработки (настроен в контейнере)

## Настройка разработки

### Go Workspace
```bash
# go.work
go 1.23.4
use ./services/task-service
```

### Зависимости проекта
```bash
# Основной go.mod
module go-task-managment-system
go 1.23.4

# task-service/go.mod
module github.com/avangero/task-service
go 1.23.4
require:
  - github.com/go-playground/validator/v10 v10.22.1
  - github.com/gofiber/fiber/v2 v2.52.8
  - github.com/google/uuid v1.6.0
```

### Docker конфигурация
```yaml
# docker-compose.yml
services:
  task-service:
    build: ./services/task-service
    ports: ["8080:8080"]
    volumes: ["./services/task-service:/app"]
    environment: ["PORT=8080"]
```

## Технические ограничения

### Версионирование
- **Go**: Требуется 1.23.4 или совместимая версия
- **Docker**: Современная версия для multi-stage builds
- **Air**: Совместимость с Go версией проверена

### Производительность
- **In-memory storage**: Быстро, но не персистентно
- **Single-threaded**: Пока один процесс на сервис
- **Memory limit**: Ограничено RAM контейнера

### Безопасность
- **Validation**: Базовая валидация входящих данных
- **CORS**: Не настроен
- **Authentication**: Отсутствует
- **HTTPS**: Не настроен

## Инструменты разработки

### Сборка и запуск
```bash
# Локальная разработка
make dev          # Запуск через Docker Compose
make build        # Сборка всех сервисов
make clean        # Очистка

# Прямой запуск
cd services/task-service
go run main.go
```

### Тестирование API
```bash
# Создание задачи
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","status":"TODO"}'

# Получение всех задач
curl http://localhost:8080/tasks

# Получение задачи по ID
curl http://localhost:8080/tasks/{id}
```

## Паттерны использования инструментов

### Fiber Framework
```go
// Инициализация
app := fiber.New()

// Middleware (будущее)
// app.Use(logger.New())
// app.Use(cors.New())

// Routing
app.Get("/", handler)
app.Post("/tasks", createTask)
app.Get("/tasks", getTasks)
app.Get("/tasks/:id", getTaskByID)
app.Put("/tasks/:id", updateTask)
app.Delete("/tasks/:id", deleteTask)
```

### Validator
```go
// Глобальный экземпляр
var validate = validator.New()

// Структура с тегами
type Task struct {
    Title  string `validate:"required"`
    Status string `validate:"required,oneof=TODO IN_PROGRESS DONE"`
}

// Валидация с обработкой ошибок
func validateTask(task *Task) error {
    if err := validate.Struct(task); err != nil {
        // Преобразование в понятные сообщения
        return formatValidationErrors(err)
    }
    return nil
}
```

### Логирование
```go
// Текущий подход
log.Printf("Задача создана: %+v", task)
log.Printf("Ошибка валидации: %v", err)

// Планируется структурированное логирование
// с использованием slog или logrus
```

## Конфигурация окружения

### Переменные окружения
- **PORT**: Порт для HTTP сервера (по умолчанию 8080)
- **Будущие**: DATABASE_URL, LOG_LEVEL, etc.

### Файлы конфигурации
- **go.work**: Конфигурация workspace
- **docker-compose.yml**: Конфигурация для разработки
- **Dockerfile**: Конфигурация сборки сервиса

## Планы технического развития

### Ближайшие улучшения
1. **Структурированное логирование**: slog или logrus
2. **Конфигурация**: viper для управления конфигурацией
3. **Middleware**: Добавление стандартных middleware
4. **Тестирование**: Unit и integration тесты

### Будущие технологии
1. **PostgreSQL**: Для персистентного хранилища
2. **Redis**: Для кеширования
3. **Prometheus**: Для метрик
4. **Jaeger**: Для трассировки
5. **gRPC**: Для межсервисного взаимодействия 