package repositories

import (
	"context"
	"database/sql"
	"log"

	"github.com/avangero/auth-service/internal/lang"
	"github.com/avangero/auth-service/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
}

// userRepository реализация UserRepository
type userRepository struct {
	db       *sqlx.DB
	messages lang.Messages
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *sqlx.DB, messages lang.Messages) UserRepository {
	return &userRepository{
		db:       db,
		messages: messages,
	}
}

// Create создает нового пользователя в БД
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, created_at) 
		VALUES (:id, :email, :password_hash, :role, :created_at)`

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		log.Printf(r.messages.Get(lang.LogUserCreateFailed), user.Email, err)
		return err
	}

	log.Printf(r.messages.Get(lang.LogUserCreateSuccess), user.Email)
	return nil
}

// GetByEmail находит пользователя по email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE email = $1"

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf(r.messages.Get(lang.LogUserNotFoundRepo), "email", email)
			return nil, nil // Возвращаем nil, nil при отсутствии пользователя
		}
		log.Printf(r.messages.Get(lang.LogDatabaseError), "email", email, err)
		return nil, err
	}

	return &user, nil
}

// GetByID находит пользователя по ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE id = $1"

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf(r.messages.Get(lang.LogUserNotFoundRepo), "ID", id.String())
			return nil, nil // Возвращаем nil, nil при отсутствии пользователя
		}
		log.Printf(r.messages.Get(lang.LogDatabaseError), "ID", id.String(), err)
		return nil, err
	}

	return &user, nil
}

// EmailExists проверяет существование пользователя с данным email
func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"

	err := r.db.GetContext(ctx, &exists, query, email)
	if err != nil {
		log.Printf(r.messages.Get(lang.LogEmailExistsCheck), email, err)
		return false, err
	}

	return exists, nil
}
