# Auth Service

Микросервис аутентификации и авторизации для Портала Обучения Сотрудников.

## Быстрый старт

### 1. Настройка окружения

Скопируйте `.env.example` в `.env` и настройте переменные:

```bash
cp .env.example .env
```

Отредактируйте `.env` файл под ваши настройки:

```bash
# Минимальные обязательные настройки
JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters
DATABASE_URL=postgres://username:password@localhost:5432/auth_db?sslmode=disable
```

### 2. Локальный запуск

**Требования:**
- Go 1.23.4+
- PostgreSQL 15+

**Запуск:**

```bash
# Установка зависимостей
go mod download

# Запуск базы данных (если через Docker)
docker-compose up -d postgres

# Запуск сервиса
go run main.go
```

### 3. Запуск через Docker Compose

```bash
# Запуск всех сервисов
docker-compose up -d

# Просмотр логов
docker-compose logs -f auth-service

# Остановка
docker-compose down
```

## Конфигурация

### Файлы окружения

- **`.env`** - единый файл конфигурации для локальной разработки и Docker
- **`.env.example`** - шаблон с описанием переменных

### Как это работает

**Локальная разработка:**
```bash
# .env содержит
DB_HOST=localhost  # ← подключение к локальной БД
```

**Docker Compose:**
```bash
# Docker автоматически переопределяет:
DB_HOST=postgres   # ← подключение к БД в контейнере
```

Docker Compose использует тот же `.env` файл, но переопределяет `DB_HOST=postgres` через `environment` секцию.

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `PORT` | HTTP порт сервиса | `8081` |
| `JWT_SECRET` | Секретный ключ для JWT | **обязательно** |
| `DATABASE_URL` | URL подключения к PostgreSQL | **обязательно** |
| `BCRYPT_COST` | Стоимость хеширования паролей | `12` |
| `GO_ENV` | Тип окружения | `development` |

## API Endpoints

- `POST /api/v1/register` - Регистрация пользователя
- `POST /api/v1/login` - Вход в систему
- `GET /api/v1/me` - Информация о текущем пользователе
- `POST /api/v1/validate` - Валидация JWT токена
- `GET /` - Health check

## Архитектура

Сервис построен на принципах Clean Architecture:

```
internal/
├── config/          # Конфигурация
├── database/        # Подключение к БД
├── handlers/        # HTTP обработчики
├── lang/            # Интернационализация
├── middleware/      # Middleware
├── models/          # Модели данных
├── repositories/    # Репозитории
├── services/        # Бизнес-логика
└── validators/      # Валидация
```

## Тестирование

```bash
# Запуск всех тестов
go test ./...

# Запуск с покрытием
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Безопасность

- JWT токены с TTL 24 часа
- bcrypt хеширование паролей (cost 12)
- Защита от SQL инъекций
- Role-based доступ
- Валидация всех входящих данных

## Production готовность

- ✅ Clean Architecture
- ✅ Централизованная интернационализация  
- ✅ Context propagation
- ✅ Comprehensive testing (27 тестов)
- ✅ Multi-level logging
- ✅ Security best practices
- ✅ Docker support 