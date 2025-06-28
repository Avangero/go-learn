package config_test

import (
	"os"
	"testing"

	"github.com/avangero/auth-service/internal/config"
	"github.com/avangero/auth-service/internal/lang/ru"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoader_Load_Success(t *testing.T) {
	// Подготовка
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("BCRYPT_COST", "10")
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("BCRYPT_COST")
	}()

	messages := ru.NewRussianMessages()
	loader := config.NewLoader(messages)

	// Выполнение
	cfg, err := loader.Load()

	// Проверка
	require.NoError(t, err)
	assert.Equal(t, "8081", cfg.Port)
	assert.Equal(t, "test-secret-key", cfg.JWT.Secret)
	assert.Equal(t, 10, cfg.BCryptCost)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "auth_db", cfg.Database.Name)
}

func TestLoader_Load_MissingJWTSecret(t *testing.T) {
	// Подготовка
	os.Unsetenv("JWT_SECRET")

	messages := ru.NewRussianMessages()
	loader := config.NewLoader(messages)

	// Выполнение
	cfg, err := loader.Load()

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "JWT_SECRET не установлен")
}

func TestLoader_Load_InvalidBCryptCost(t *testing.T) {
	// Подготовка
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("BCRYPT_COST", "invalid")
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("BCRYPT_COST")
	}()

	messages := ru.NewRussianMessages()
	loader := config.NewLoader(messages)

	// Выполнение
	cfg, err := loader.Load()

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "BCRYPT_COST")
}

func TestLoader_Load_BCryptCostOutOfRange(t *testing.T) {
	tests := []struct {
		name       string
		bcryptCost string
	}{
		{"Too low", "3"},
		{"Too high", "32"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Подготовка
			os.Setenv("JWT_SECRET", "test-secret")
			os.Setenv("BCRYPT_COST", tt.bcryptCost)
			defer func() {
				os.Unsetenv("JWT_SECRET")
				os.Unsetenv("BCRYPT_COST")
			}()

			messages := ru.NewRussianMessages()
			loader := config.NewLoader(messages)

			// Выполнение
			cfg, err := loader.Load()

			// Проверка
			assert.Error(t, err)
			assert.Nil(t, cfg)
			assert.Contains(t, err.Error(), "должно быть от 4 до 31")
		})
	}
}

func TestLoader_Load_DatabaseURL_Success(t *testing.T) {
	// Подготовка
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("DATABASE_URL", "postgres://testuser:testpass@testhost:5433/testdb")
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("DATABASE_URL")
	}()

	messages := ru.NewRussianMessages()
	loader := config.NewLoader(messages)

	// Выполнение
	cfg, err := loader.Load()

	// Проверка
	require.NoError(t, err)
	assert.Equal(t, "testhost", cfg.Database.Host)
	assert.Equal(t, "5433", cfg.Database.Port)
	assert.Equal(t, "testuser", cfg.Database.User)
	assert.Equal(t, "testpass", cfg.Database.Password)
	assert.Equal(t, "testdb", cfg.Database.Name)
}

func TestLoader_Load_DatabaseURL_Priority(t *testing.T) {
	// Подготовка - DATABASE_URL должен иметь приоритет над отдельными переменными
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("DATABASE_URL", "postgres://urluser:urlpass@urlhost:5433/urldb")
	os.Setenv("DB_HOST", "separatehost")
	os.Setenv("DB_USER", "separateuser")
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_USER")
	}()

	messages := ru.NewRussianMessages()
	loader := config.NewLoader(messages)

	// Выполнение
	cfg, err := loader.Load()

	// Проверка - должны использоваться значения из DATABASE_URL
	require.NoError(t, err)
	assert.Equal(t, "urlhost", cfg.Database.Host)
	assert.Equal(t, "urluser", cfg.Database.User)
	assert.Equal(t, "urldb", cfg.Database.Name)
}
