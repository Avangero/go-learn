package ru

import "github.com/avangero/auth-service/internal/lang"

// NewRussianMessages создает провайдер русских сообщений
func NewRussianMessages() lang.Messages {
	messages := map[lang.MessageKey]string{
		// Database
		lang.DBConnectionError: "Ошибка подключения к базе данных",
		lang.DBPingError:       "Ошибка проверки подключения к БД",
		lang.DBConnected:       "✅ Подключение к PostgreSQL успешно",

		// Config
		lang.JWTSecretMissing:  "JWT_SECRET не установлен",
		lang.BCryptCostInvalid: "Неверное значение BCRYPT_COST",

		// Auth
		lang.InvalidRequestFormat: "Неверный формат запроса",
		lang.UserAlreadyExists:    "Пользователь с таким email уже существует",
		lang.InvalidCredentials:   "Неверный email или пароль",
		lang.TokenNotProvided:     "Токен не предоставлен",
		lang.TokenInvalid:         "Недействительный токен",
		lang.UserNotFound:         "Пользователь не найден",
		lang.InternalServerError:  "Внутренняя ошибка сервера",
		lang.UserRegistered:       "✅ Новый пользователь зарегистрирован",
		lang.UserLoggedIn:         "✅ Пользователь вошел в систему",

		// Validation
		lang.ValidationFieldRequired: "Поле обязательно для заполнения",
		lang.ValidationEmailInvalid:  "Поле должно быть действительным email адресом",
		lang.ValidationPasswordMin:   "Поле должно содержать минимум символов",
		lang.ValidationRoleInvalid:   "Поле должно быть одним из разрешенных значений",
	}

	return lang.NewMessageProvider(messages)
}
