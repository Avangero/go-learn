package responses

import "github.com/avangero/auth-service/internal/models"

// TokenResponse представляет ответ с JWT токеном
type TokenResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// StatusResponse представляет ответ о статусе
type StatusResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Version string `json:"version"`
}

// ValidationResponse представляет ответ валидации токена
type ValidationResponse struct {
	Valid bool        `json:"valid"`
	User  models.User `json:"user"`
}
