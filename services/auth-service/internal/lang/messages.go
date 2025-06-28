package lang

import (
	"fmt"
)

// MessageKey представляет ключ сообщения
type MessageKey string

// Ключи сообщений
const (
	// Database messages
	DBConnectionError MessageKey = "db.connection.error"
	DBPingError       MessageKey = "db.ping.error"
	DBConnected       MessageKey = "db.connected"

	// Config messages
	JWTSecretMissing  MessageKey = "config.jwt_secret.missing"
	BCryptCostInvalid MessageKey = "config.bcrypt_cost.invalid"

	// Auth messages
	InvalidRequestFormat MessageKey = "auth.request.invalid_format"
	UserAlreadyExists    MessageKey = "auth.user.already_exists"
	InvalidCredentials   MessageKey = "auth.credentials.invalid"
	TokenNotProvided     MessageKey = "auth.token.not_provided"
	TokenInvalid         MessageKey = "auth.token.invalid"
	UserNotFound         MessageKey = "auth.user.not_found"
	InternalServerError  MessageKey = "auth.server.internal_error"
	UserRegistered       MessageKey = "auth.user.registered"
	UserLoggedIn         MessageKey = "auth.user.logged_in"

	// Validation messages
	ValidationFieldRequired MessageKey = "validation.field.required"
	ValidationEmailInvalid  MessageKey = "validation.email.invalid"
	ValidationPasswordMin   MessageKey = "validation.password.min"
	ValidationRoleInvalid   MessageKey = "validation.role.invalid"

	// Logging messages - Handler level
	LogRegistrationRequest MessageKey = "log.registration.request"
	LogLoginRequest        MessageKey = "log.login.request"
	LogValidationFailed    MessageKey = "log.validation.failed"
	LogRegistrationFailed  MessageKey = "log.registration.failed"
	LogRegistrationSuccess MessageKey = "log.registration.success"
	LogLoginFailed         MessageKey = "log.login.failed"
	LogLoginSuccess        MessageKey = "log.login.success"
	LogGetMeFailed         MessageKey = "log.getme.failed"
	LogGetMeSuccess        MessageKey = "log.getme.success"
	LogParseRequestFailed  MessageKey = "log.parse.request.failed"

	// Logging messages - Service level
	LogAttemptingRegistration MessageKey = "log.service.attempting.registration"
	LogCheckEmailExists       MessageKey = "log.service.check.email.exists"
	LogEmailAlreadyExists     MessageKey = "log.service.email.already.exists"
	LogPasswordHashError      MessageKey = "log.service.password.hash.error"
	LogUserCreateError        MessageKey = "log.service.user.create.error"
	LogJWTGenerateError       MessageKey = "log.service.jwt.generate.error"
	LogRegistrationComplete   MessageKey = "log.service.registration.complete"
	LogAttemptingLogin        MessageKey = "log.service.attempting.login"
	LogDatabaseErrorLogin     MessageKey = "log.service.database.error.login"
	LogUserNotFoundLogin      MessageKey = "log.service.user.not.found.login"
	LogInvalidPassword        MessageKey = "log.service.invalid.password"
	LogLoginComplete          MessageKey = "log.service.login.complete"
	LogJWTParseError          MessageKey = "log.service.jwt.parse.error"
	LogJWTInvalid             MessageKey = "log.service.jwt.invalid"
	LogUserFetchError         MessageKey = "log.service.user.fetch.error"
	LogUserNotFoundValidation MessageKey = "log.service.user.not.found.validation"

	// Logging messages - Repository level
	LogUserCreateSuccess MessageKey = "log.repo.user.create.success"
	LogUserCreateFailed  MessageKey = "log.repo.user.create.failed"
	LogUserNotFoundRepo  MessageKey = "log.repo.user.not.found"
	LogDatabaseError     MessageKey = "log.repo.database.error"
	LogEmailExistsCheck  MessageKey = "log.repo.email.exists.check"

	// Logging messages - Middleware level
	LogJWTMissingHeader     MessageKey = "log.jwt.missing.header"
	LogJWTInvalidFormat     MessageKey = "log.jwt.invalid.format"
	LogJWTValidationFailed  MessageKey = "log.jwt.validation.failed"
	LogJWTValidationSuccess MessageKey = "log.jwt.validation.success"
)

// Messages интерфейс для получения сообщений
type Messages interface {
	Get(key MessageKey, args ...interface{}) string
	GetValidationError(field, tag, param string) string
}

// MessageProvider базовая реализация провайдера сообщений
type MessageProvider struct {
	messages map[MessageKey]string
}

// NewMessageProvider создает новый провайдер сообщений
func NewMessageProvider(messages map[MessageKey]string) *MessageProvider {
	return &MessageProvider{messages: messages}
}

// Get возвращает сообщение по ключу с возможностью форматирования
func (m *MessageProvider) Get(key MessageKey, args ...interface{}) string {
	if msg, exists := m.messages[key]; exists {
		if len(args) > 0 {
			// Используем fmt.Sprintf для форматирования
			return fmt.Sprintf(msg, args...)
		}
		return msg
	}
	return string(key) // Fallback на ключ если сообщение не найдено
}

// GetValidationError возвращает сообщение об ошибке валидации
func (m *MessageProvider) GetValidationError(field, tag, param string) string {
	switch tag {
	case "required":
		return m.Get(ValidationFieldRequired) + ": " + field
	case "email":
		return m.Get(ValidationEmailInvalid) + ": " + field
	case "min":
		return m.Get(ValidationPasswordMin) + ": " + field + " (мин. " + param + " символов)"
	case "oneof":
		return m.Get(ValidationRoleInvalid) + ": " + field
	default:
		return "Ошибка валидации поля: " + field
	}
}
