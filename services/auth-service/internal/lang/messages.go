package lang

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
			// Простое форматирование - можно расширить через fmt.Sprintf
			return msg
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
