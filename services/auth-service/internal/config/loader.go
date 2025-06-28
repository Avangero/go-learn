package config

import (
	"os"
	"strconv"

	"github.com/avangero/auth-service/internal/lang"
)

// Loader загружает конфигурацию из переменных окружения
type Loader struct {
	messages lang.Messages
}

// NewLoader создает новый загрузчик конфигурации
func NewLoader(messages lang.Messages) *Loader {
	return &Loader{messages: messages}
}

// Load загружает конфигурацию из переменных окружения
func (l *Loader) Load() (*Config, error) {
	cfg := &Config{
		Port: l.getEnv("PORT", "8081"),
		Database: DatabaseConfig{
			Host:     l.getEnv("DB_HOST", "localhost"),
			Port:     l.getEnv("DB_PORT", "5432"),
			User:     l.getEnv("DB_USER", "postgres"),
			Password: l.getEnv("DB_PASSWORD", "postgres"),
			Name:     l.getEnv("DB_NAME", "auth_db"),
		},
		JWT: JWTConfig{
			Secret: l.getEnv("JWT_SECRET", ""),
		},
	}

	// Валидация конфигурации
	if err := l.validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate валидирует конфигурацию
func (l *Loader) validate(cfg *Config) error {
	validator := NewConfigValidator(l.messages)
	return validator.Validate(cfg)
}

// getEnv возвращает значение переменной окружения или defaultValue
func (l *Loader) getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseInt парсит строку в int с обработкой ошибок
func (l *Loader) parseInt(value string, defaultValue int) (int, error) {
	if value == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(value)
}
