package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/avangero/auth-service/internal/lang"
)

// ConfigValidator валидатор конфигурации
type ConfigValidator struct {
	messages lang.Messages
}

// NewConfigValidator создает новый валидатор конфигурации
func NewConfigValidator(messages lang.Messages) *ConfigValidator {
	return &ConfigValidator{
		messages: messages,
	}
}

// Validate валидирует конфигурацию
func (v *ConfigValidator) Validate(cfg *Config) error {
	// Проверка JWT Secret
	if cfg.JWT.Secret == "" {
		return errors.New(v.messages.Get(lang.JWTSecretMissing))
	}

	// Парсинг и валидация BCrypt cost
	bcryptCostStr := v.getEnv("BCRYPT_COST", "12")
	var err error
	cfg.BCryptCost, err = strconv.Atoi(bcryptCostStr)
	if err != nil {
		return errors.New(v.messages.Get(lang.BCryptCostInvalid) + ": " + err.Error())
	}

	// Проверка диапазона BCrypt cost
	if cfg.BCryptCost < 4 || cfg.BCryptCost > 31 {
		return errors.New(v.messages.Get(lang.BCryptCostInvalid) + ": должно быть от 4 до 31")
	}

	return nil
}

// getEnv возвращает значение переменной окружения или defaultValue
func (v *ConfigValidator) getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
