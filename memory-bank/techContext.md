# Technical Context - Портал Обучения для Сотрудников

## Основной Технологический Стек

### Реализованные Технологии (Auth Service ✅)
- **Go 1.23.4** - основной язык
- **Fiber v2.52.8** - HTTP веб-фреймворк
- **PostgreSQL 15** - база данных (auth_db, user_db, course_db)
- **sqlx v1.4.0** - SQL расширения для Go
- **JWT v5.2.0** - аутентификация между сервисами
- **bcrypt** - безопасное хеширование паролей
- **validator v10.22.1** - валидация данных
- **testify/mock** - мокирование для тестов

### Архитектурные Решения

**Clean Architecture Structure (Auth Service):**
```
main.go (Composition Root)
    ↓
internal/
├── config/            # Конфигурация (3 файла по SRP)
├── database/          # Подключение к PostgreSQL
├── lang/              # Интернационализация (NEW)
├── handlers/          # HTTP обработчики
├── services/          # Бизнес-логика
├── repositories/      # Доступ к данным
├── models/            # Разделено на requests/responses
├── middleware/        # JWT и CORS
└── validators/        # Валидация с i18n
```

**Ключевые Улучшения:**
- ✅ **Single Responsibility:** config разделен на loader/validator/config
- ✅ **Интернационализация:** lang/ru система сообщений
- ✅ **Разделение моделей:** requests/responses/entities
- ✅ **100% тестирование:** 27 тестов, 78% покрытие

## Система Централизованной Интернационализации (ПОЛНАЯ)

**Архитектура централизованных сообщений для ВСЕГО приложения:**
```go
// internal/lang/messages.go
type MessageKey string
type Messages interface {
    Get(key MessageKey) string
}

type MessageProvider struct {
    messages map[MessageKey]string
}

func (p *MessageProvider) Get(key MessageKey) string {
    message, exists := p.messages[key]
    if !exists {
        return string(key)
    }
    return message
}
```

**37 ключей сообщений по архитектурным уровням:**
```go
// API Messages
const (
    UserNotFound         MessageKey = "auth.user.not_found"
    UserAlreadyExists    MessageKey = "auth.user.already_exists"
    InvalidCredentials   MessageKey = "auth.invalid.credentials"
)

// Handler Level Logging (10 ключей)
const (
    LogRegistrationRequest  MessageKey = "log.registration.request"
    LogLoginRequest         MessageKey = "log.login.request"
    LogValidationFailed     MessageKey = "log.validation.failed"
    LogRegistrationSuccess  MessageKey = "log.registration.success"
    LogLoginSuccess         MessageKey = "log.login.success"
    // ... и другие
)

// Service Level Logging (13 ключей)
const (
    LogAttemptingRegistration MessageKey = "log.service.attempting.registration"
    LogEmailAlreadyExists     MessageKey = "log.service.email.already.exists"
    LogRegistrationComplete   MessageKey = "log.service.registration.complete"
    // ... и другие
)

// Repository Level Logging (5 ключей)
const (
    LogUserCreateSuccess MessageKey = "log.repo.user.create.success"
    LogUserCreateFailed  MessageKey = "log.repo.user.create.failed"
    LogDatabaseError     MessageKey = "log.repo.database.error"
    // ... и другие
)

// Middleware Level Logging (4 ключа)
const (
    LogJWTValidationSuccess MessageKey = "log.jwt.validation.success"
    LogJWTValidationFailed  MessageKey = "log.jwt.validation.failed"
    // ... и другие
)
```

**Интеграция во все компоненты архитектуры:**
```go
// Все конструкторы требуют messages параметр
func NewUserRepository(db *sqlx.DB, messages lang.Messages) UserRepository
func NewAuthService(userRepo repositories.UserRepository, jwtSecret string, bcryptCost int, messages lang.Messages) AuthService
func NewAuthHandler(authService services.AuthService, messages lang.Messages) *AuthHandler
func JWTMiddleware(authService services.AuthService, messages lang.Messages) fiber.Handler

// main.go - единый источник truth
messages := ru.NewRussianMessages()
userRepo := repositories.NewUserRepository(db, messages)
authService := services.NewAuthService(userRepo, cfg.JWT.Secret, cfg.BCryptCost, messages)
authHandler := handlers.NewAuthHandler(authService, messages)
```

**Русские переводы с форматированием:**
```go
// internal/lang/ru/messages.go
"log.registration.request":        "Запрос регистрации с IP: %s",
"log.service.attempting.registration": "Попытка регистрации пользователя с email: %s",
"log.repo.user.create.success":    "Пользователь успешно создан с email %s",
"log.jwt.validation.success":      "JWT middleware: валидация токена успешна для IP %s, пользователь: %s",
```

**Критические преимущества:**
- ✅ **0 hardcoded текстов** в коде (включая логи)
- ✅ **Типобезопасные ключи** - MessageKey константы
- ✅ **Форматирование параметров** через fmt.Sprintf
- ✅ **Многоуровневое логирование** по архитектурным слоям
- ✅ **100% готовность** к добавлению новых языков

## Структура Тестирования с Lang Интеграцией

**Полное покрытие тестами (27 тестов) - интегрированы с lang системой:**
```
tests/unit/
├── config/            # 4 теста (loader, validator)
├── validators/        # 8 тестов (все сценарии валидации)
├── repositories/      # 7 тестов (с sqlmock + messages)
└── services/          # 8 тестов (с testify/mock + messages)
```

**Техники тестирования с lang поддержкой:**
- **sqlmock** для тестирования репозиториев с messages интеграцией
- **testify/mock** для мокирования сервисов с messages параметром
- **bcrypt cost 4** для быстрых тестов
- **Параллельное выполнение** тестов
- **Централизованные сообщения** во всех тестах

**Паттерн тестирования с lang системой:**
```go
// Каждый тест требует messages объект
func TestUserRepository_Create(t *testing.T) {
    messages := ru.NewRussianMessages()
    repo := repositories.NewUserRepository(db, messages)
    // тест использует централизованные сообщения
}

func TestAuthService_Register_Success(t *testing.T) {
    mockRepo := new(MockUserRepository)
    messages := ru.NewRussianMessages()
    authService := services.NewAuthService(mockRepo, "test-secret", 4, messages)
    // сервис логирует через lang систему во время теста
}
```

**Обновленные expectations:**
- Сообщения с правильным регистром через lang систему
- Передача messages во все mock объекты
- Проверка корректности интернационализированных сообщений
- Все 27 тестов проходят с lang интеграцией

## Безопасность

**Реализованные паттерны:**
- **Параметризованные запросы** - защита от SQL-инъекций
- **bcrypt cost 12** в production, 4 в тестах
- **JWT с 24-часовым TTL** 
- **Role-based access** (employee/manager)
- **Валидация на boundaries** (HTTP endpoints)

## Конфигурация

**Env-based конфигурация:**
```bash
# Обязательные переменные
JWT_SECRET=your-secret-key
DATABASE_URL=postgres://user:pass@localhost/auth_db

# Опциональные с defaults
PORT=8081
BCRYPT_COST=12
```

**Структура конфигурации:**
```go
// internal/config/config.go - структуры
type Config struct {
    Port     string
    Database DatabaseConfig
    JWT      JWTConfig
}

// internal/config/loader.go - загрузка
// internal/config/validator.go - валидация
```

## Контейнеризация

**Docker Compose архитектура:**
```yaml
services:
  auth-service:
    build: ./services/auth-service
    ports: ["8081:8081"]
    depends_on:
      postgres:
        condition: service_healthy
        
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: auth_db
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user -d auth_db"]
```

## Следующие Технические Задачи

**User Service (8082):**
- Применить Clean Architecture + ВСЕ стандарты Auth Service
- **Централизованная lang система** с самого начала
- **Messages объект** во всех компонентах (repository, service, handler)
- **Многоуровневое логирование** через lang ключи
- JWT валидацию из Auth Service
- 100% покрытие тестами с lang интеграцией

**Course Service (8083):**
- Применить Clean Architecture + ВСЕ стандарты Auth Service
- **Централизованная lang система** с самого начала  
- **Messages объект** во всех компонентах
- **Многоуровневое логирование** через lang ключи
- Система курсов/модулей/навыков
- Интеграция с User Service для прогресса

**API Gateway:**
- Единая точка входа с централизованными сообщениями
- Маршрутизация по сервисам с lang логированием
- Централизованная авторизация через lang систему

**Критические требования для всех новых сервисов:**
- ✅ **НЕТ hardcoded текстов** вообще
- ✅ **Lang система** с первого коммита
- ✅ **Messages parameter** в ВСЕХ конструкторах
- ✅ **Типобезопасные ключи** для всех сообщений
- ✅ **Многоуровневое логирование** по архитектурным слоям