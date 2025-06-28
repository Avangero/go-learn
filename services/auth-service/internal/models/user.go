package models

import (
	"time"

	"github.com/google/uuid"
)

// User представляет модель пользователя в системе
type User struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Email    string    `json:"email" db:"email" validate:"required,email"`
	Password string    `json:"-" db:"password_hash"` // не возвращается в JSON
	Role     string    `json:"role" db:"role" validate:"required,oneof=employee manager"`
	Created  time.Time `json:"created_at" db:"created_at"`
}

// RegisterRequest представляет запрос на регистрацию
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=employee manager"`
}

// LoginRequest представляет запрос на вход
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// TokenResponse представляет ответ с JWT токеном
type TokenResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
