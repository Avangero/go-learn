package repositories

import (
	"github.com/avangero/auth-service/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
	EmailExists(email string) (bool, error)
}

// userRepository реализация UserRepository
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

// Create создает нового пользователя в БД
func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, created_at) 
		VALUES (:id, :email, :password_hash, :role, :created_at)`

	_, err := r.db.NamedExec(query, user)
	return err
}

// GetByEmail находит пользователя по email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE email = $1"

	err := r.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID находит пользователя по ID
func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE id = $1"

	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// EmailExists проверяет существование пользователя с данным email
func (r *userRepository) EmailExists(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"

	err := r.db.Get(&exists, query, email)
	if err != nil {
		return false, err
	}

	return exists, nil
}
