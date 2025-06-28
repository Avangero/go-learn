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

		// Logging messages - Handler level
		lang.LogRegistrationRequest: "Запрос регистрации с IP: %s",
		lang.LogLoginRequest:        "Запрос входа с IP: %s",
		lang.LogValidationFailed:    "Ошибка валидации запроса с IP %s: %v",
		lang.LogRegistrationFailed:  "Регистрация не удалась для IP %s: %v",
		lang.LogRegistrationSuccess: "Регистрация успешна для IP %s, email: %s",
		lang.LogLoginFailed:         "Вход не удался для IP %s: %v",
		lang.LogLoginSuccess:        "Вход успешен для IP %s, email: %s",
		lang.LogGetMeFailed:         "GetMe не удался: пользователь не найден в контексте с IP %s",
		lang.LogGetMeSuccess:        "GetMe успешен для IP %s, пользователь: %s",
		lang.LogParseRequestFailed:  "Ошибка парсинга запроса с IP %s: %v",

		// Logging messages - Service level
		lang.LogAttemptingRegistration: "Попытка регистрации пользователя с email: %s",
		lang.LogCheckEmailExists:       "Ошибка проверки существования email %s: %v",
		lang.LogEmailAlreadyExists:     "Регистрация не удалась: пользователь с email %s уже существует",
		lang.LogPasswordHashError:      "Ошибка хеширования пароля для пользователя %s: %v",
		lang.LogUserCreateError:        "Ошибка создания пользователя %s: %v",
		lang.LogJWTGenerateError:       "Ошибка генерации JWT для пользователя %s: %v",
		lang.LogRegistrationComplete:   "Регистрация пользователя успешно завершена для email: %s",
		lang.LogAttemptingLogin:        "Попытка входа пользователя с email: %s",
		lang.LogDatabaseErrorLogin:     "Ошибка БД при входе для email %s: %v",
		lang.LogUserNotFoundLogin:      "Вход не удался: пользователь не найден с email %s",
		lang.LogInvalidPassword:        "Вход не удался: неверный пароль для email %s",
		lang.LogLoginComplete:          "Вход пользователя успешно завершен для email: %s",
		lang.LogJWTParseError:          "Ошибка парсинга JWT токена: %v",
		lang.LogJWTInvalid:             "Недействительный JWT токен или claims",
		lang.LogUserFetchError:         "Ошибка получения пользователя по ID %s при валидации токена: %v",
		lang.LogUserNotFoundValidation: "Пользователь не найден с ID %s при валидации токена",

		// Logging messages - Repository level
		lang.LogUserCreateSuccess: "Пользователь успешно создан с email %s",
		lang.LogUserCreateFailed:  "Ошибка создания пользователя с email %s: %v",
		lang.LogUserNotFoundRepo:  "Пользователь не найден с %s: %s",
		lang.LogDatabaseError:     "Ошибка БД при операции с %s %s: %v",
		lang.LogEmailExistsCheck:  "Ошибка БД при проверке существования email %s: %v",

		// Logging messages - Middleware level
		lang.LogJWTMissingHeader:     "JWT middleware: отсутствует заголовок Authorization с IP %s",
		lang.LogJWTInvalidFormat:     "JWT middleware: неверный формат заголовка Authorization с IP %s",
		lang.LogJWTValidationFailed:  "JWT middleware: валидация токена не удалась с IP %s: %v",
		lang.LogJWTValidationSuccess: "JWT middleware: валидация токена успешна для IP %s, пользователь: %s",
	}

	return lang.NewMessageProvider(messages)
}
