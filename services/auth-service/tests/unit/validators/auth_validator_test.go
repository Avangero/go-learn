package validators_test

import (
	"testing"

	"github.com/avangero/auth-service/internal/lang/ru"
	"github.com/avangero/auth-service/internal/models/requests"
	"github.com/avangero/auth-service/internal/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthValidator_Validate_RegisterRequest_Success(t *testing.T) {
	// Подготовка
	messages := ru.NewRussianMessages()
	validator := validators.NewAuthValidator(messages)

	req := &requests.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Role:     "employee",
	}

	// Выполнение
	err := validator.Validate(req)

	// Проверка
	assert.NoError(t, err)
}

func TestAuthValidator_Validate_RegisterRequest_InvalidEmail(t *testing.T) {
	// Подготовка
	messages := ru.NewRussianMessages()
	validator := validators.NewAuthValidator(messages)

	req := &requests.RegisterRequest{
		Email:    "invalid-email",
		Password: "password123",
		Role:     "employee",
	}

	// Выполнение
	err := validator.Validate(req)

	// Проверка
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email адресом")
}

func TestAuthValidator_Validate_RegisterRequest_ShortPassword(t *testing.T) {
	// Подготовка
	messages := ru.NewRussianMessages()
	validator := validators.NewAuthValidator(messages)

	req := &requests.RegisterRequest{
		Email:    "test@example.com",
		Password: "123", // слишком короткий
		Role:     "employee",
	}

	// Выполнение
	err := validator.Validate(req)

	// Проверка
	require.Error(t, err)
	assert.Contains(t, err.Error(), "минимум")
}

func TestAuthValidator_Validate_RegisterRequest_InvalidRole(t *testing.T) {
	// Подготовка
	messages := ru.NewRussianMessages()
	validator := validators.NewAuthValidator(messages)

	req := &requests.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Role:     "admin", // недопустимая роль
	}

	// Выполнение
	err := validator.Validate(req)

	// Проверка
	require.Error(t, err)
	assert.Contains(t, err.Error(), "разрешенных значений")
}

func TestAuthValidator_Validate_LoginRequest_Success(t *testing.T) {
	// Подготовка
	messages := ru.NewRussianMessages()
	validator := validators.NewAuthValidator(messages)

	req := &requests.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Выполнение
	err := validator.Validate(req)

	// Проверка
	assert.NoError(t, err)
}

func TestAuthValidator_Validate_LoginRequest_MissingFields(t *testing.T) {
	tests := []struct {
		name     string
		request  *requests.LoginRequest
		errorMsg string
	}{
		{
			name:     "Missing email",
			request:  &requests.LoginRequest{Password: "password123"},
			errorMsg: "обязательно для заполнения",
		},
		{
			name:     "Missing password",
			request:  &requests.LoginRequest{Email: "test@example.com"},
			errorMsg: "обязательно для заполнения",
		},
		{
			name:     "Invalid email format",
			request:  &requests.LoginRequest{Email: "invalid", Password: "password123"},
			errorMsg: "email адресом",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Подготовка
			messages := ru.NewRussianMessages()
			validator := validators.NewAuthValidator(messages)

			// Выполнение
			err := validator.Validate(tt.request)

			// Проверка
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}
