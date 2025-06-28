package validators

import (
	"fmt"
	"strings"

	"github.com/avangero/auth-service/internal/lang"
	"github.com/go-playground/validator/v10"
)

// AuthValidator валидатор для аутентификации
type AuthValidator struct {
	validator *validator.Validate
	messages  lang.Messages
}

// NewAuthValidator создает новый валидатор
func NewAuthValidator(messages lang.Messages) *AuthValidator {
	return &AuthValidator{
		validator: validator.New(),
		messages:  messages,
	}
}

// Validate валидирует структуру и возвращает отформатированные ошибки
func (v *AuthValidator) Validate(s interface{}) error {
	if err := v.validator.Struct(s); err != nil {
		return v.formatValidationErrors(err)
	}
	return nil
}

// ValidationError кастомная ошибка валидации
type ValidationError struct {
	Message string
	Field   string
	Tag     string
}

func (e ValidationError) Error() string {
	return e.Message
}

// formatValidationErrors форматирует ошибки валидации
func (v *AuthValidator) formatValidationErrors(err error) error {
	var errorMessages []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			message := v.messages.GetValidationError(
				fieldError.Field(),
				fieldError.Tag(),
				fieldError.Param(),
			)
			errorMessages = append(errorMessages, message)
		}
	}

	return &ValidationError{
		Message: fmt.Sprintf("Ошибки валидации: %s", strings.Join(errorMessages, "; ")),
	}
}
