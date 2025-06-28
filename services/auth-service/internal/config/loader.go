package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

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
		JWT: JWTConfig{
			Secret: l.getEnv("JWT_SECRET", ""),
		},
	}

	// Загружаем BCRYPT_COST
	bcryptCost, err := l.parseInt(l.getEnv("BCRYPT_COST", "12"), 12)
	if err != nil {
		return nil, fmt.Errorf("недопустимое значение BCRYPT_COST: %v", err)
	}
	cfg.BCryptCost = bcryptCost

	// Загружаем конфигурацию БД - поддерживаем DATABASE_URL и отдельные переменные
	if err := l.loadDatabaseConfig(cfg); err != nil {
		return nil, err
	}

	// Валидация конфигурации
	if err := l.validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// loadDatabaseConfig загружает конфигурацию БД из DATABASE_URL или отдельных переменных
func (l *Loader) loadDatabaseConfig(cfg *Config) error {
	// Проверяем DATABASE_URL сначала
	if databaseURL := l.getEnv("DATABASE_URL", ""); databaseURL != "" {
		return l.parseDatabaseURL(cfg, databaseURL)
	}

	// Используем отдельные переменные
	cfg.Database = DatabaseConfig{
		Host:     l.getEnv("DB_HOST", "localhost"),
		Port:     l.getEnv("DB_PORT", "5432"),
		User:     l.getEnv("DB_USER", "postgres"),
		Password: l.getEnv("DB_PASSWORD", "postgres"),
		Name:     l.getEnv("DB_NAME", "auth_db"),
	}

	return nil
}

// parseDatabaseURL парсит DATABASE_URL в DatabaseConfig
func (l *Loader) parseDatabaseURL(cfg *Config, databaseURL string) error {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return err
	}

	cfg.Database = DatabaseConfig{
		Host: u.Hostname(),
		Port: u.Port(),
		User: u.User.Username(),
		Name: strings.TrimPrefix(u.Path, "/"),
	}

	// Получаем пароль
	if password, ok := u.User.Password(); ok {
		cfg.Database.Password = password
	}

	// Если порт не указан, используем стандартный
	if cfg.Database.Port == "" {
		cfg.Database.Port = "5432"
	}

	return nil
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
