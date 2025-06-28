# Системные паттерны

## Архитектура системы

### Микросервисная архитектура
```
┌─────────────────┐
│   task-service  │
│     :8080       │
└─────────────────┘
```

**Текущее состояние**: Один сервис, но инфраструктура готова для расширения

### Go Workspace Structure
```
go-task-management-system/
├── go.work              # Workspace configuration
├── go.mod               # Main module
├── services/
│   └── task-service/    # Individual microservice
│       ├── go.mod       # Service-specific dependencies
│       ├── main.go      # Service entry point
│       └── Dockerfile   # Service containerization
└── docker-compose.yml   # Orchestration
```

## Ключевые технические решения

### 1. Модульная организация
- **Go Workspace**: Использование go.work для управления множественными модулями
- **Изолированные зависимости**: Каждый сервис имеет свой go.mod
- **Единая оркестрация**: Docker Compose для всех сервисов

### 2. Контейнеризация
- **Docker**: Каждый сервис в отдельном контейнере
- **Live-reloading**: Настроен volume mount для разработки
- **Порты**: Четкое разделение портов между сервисами

### 3. HTTP API паттерны
- **RESTful routes**: Стандартные HTTP методы и пути
- **JSON**: Единый формат обмена данными
- **Status codes**: Правильное использование HTTP статусов

## Паттерны проектирования

### 1. Data Transfer Objects (DTO)
```go
type Task struct {
    ID          string `json:"id"`
    Title       string `json:"title" validate:"required"`
    Description string `json:"description"`
    Status      string `json:"status" validate:"required,oneof=TODO IN_PROGRESS DONE"`
}
```

### 2. Repository Pattern (упрощенный)
```go
var tasks = map[string]Task{} // In-memory storage
```

### 3. Validation Pattern
```go
var validate = validator.New()

func validateTask(task *Task) error {
    // Structured validation with readable errors
}
```

### 4. Handler Pattern
```go
func createTask(c *fiber.Ctx) error {
    // 1. Parse request
    // 2. Validate data
    // 3. Process business logic
    // 4. Return response
}
```

## Отношения компонентов

### Request Flow
```
Client → Fiber Router → Handler → Validator → Storage → Response
```

### Error Handling Flow
```
Error → Logger → Structured Message → HTTP Status → Client
```

## Критические пути реализации

### 1. Создание задачи
1. **Parsing**: `c.BodyParser(task)`
2. **Validation**: `validateTask(task)`
3. **ID Generation**: `uuid.New().String()`
4. **Storage**: `tasks[task.ID] = *task`
5. **Response**: `c.Status(201).JSON(task)`

### 2. Обработка ошибок
1. **Capture**: Перехват ошибки на каждом этапе
2. **Log**: Логирование для отладки
3. **Transform**: Преобразование в понятное сообщение
4. **Return**: Возврат с правильным HTTP статусом

## Архитектурные принципы

### 1. Separation of Concerns
- **Routing**: Fiber handles HTTP concerns
- **Validation**: Validator handles data validation
- **Business Logic**: Handlers contain business rules
- **Storage**: Separate storage layer (even if in-memory)

### 2. Dependency Injection (готовность)
- **Global validator**: Один экземпляр валидатора
- **Конфигурация**: Через переменные окружения

### 3. Extensibility
- **Модульность**: Легко добавить новые сервисы
- **Стандартизация**: Единый подход к обработке HTTP
- **Контейнеризация**: Готовность к оркестрации

## Планы развития архитектуры

### Ближайшие улучшения
1. **Database Layer**: Замена in-memory на PostgreSQL
2. **Config Management**: Централизованная конфигурация
3. **Middleware**: Добавление логирования, CORS, etc.

### Будущие сервисы
1. **user-service**: Управление пользователями
2. **notification-service**: Уведомления
3. **api-gateway**: Единая точка входа 