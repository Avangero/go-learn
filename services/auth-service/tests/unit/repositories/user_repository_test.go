package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/avangero/auth-service/internal/lang/ru"
	"github.com/avangero/auth-service/internal/models"
	"github.com/avangero/auth-service/internal/repositories"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	// Используем реальную БД для тестов (можно заменить на testcontainers)
	db, err := sqlx.Connect("postgres", "postgres://test_user:test_password@localhost:5432/test_auth_db?sslmode=disable")
	require.NoError(t, err)

	// Очищаем таблицу перед каждым тестом
	_, err = db.Exec("DELETE FROM users")
	require.NoError(t, err)

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	messages := ru.NewRussianMessages()
	repo := repositories.NewUserRepository(db, messages)

	user := &models.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     "employee",
		Created:  time.Now(),
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	// Проверяем, что пользователь создан
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE email = $1", user.Email)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	messages := ru.NewRussianMessages()
	repo := repositories.NewUserRepository(db, messages)

	// Создаем тестового пользователя
	originalUser := &models.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     "employee",
		Created:  time.Now(),
	}

	err := repo.Create(context.Background(), originalUser)
	require.NoError(t, err)

	// Получаем пользователя по email
	user, err := repo.GetByEmail(context.Background(), "test@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, originalUser.Email, user.Email)
	assert.Equal(t, originalUser.ID, user.ID)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	messages := ru.NewRussianMessages()
	repo := repositories.NewUserRepository(db, messages)

	// Пытаемся получить несуществующего пользователя
	user, err := repo.GetByEmail(context.Background(), "nonexistent@example.com")
	assert.NoError(t, err)
	assert.Nil(t, user) // Должен вернуть nil, nil при отсутствии пользователя
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	messages := ru.NewRussianMessages()
	repo := repositories.NewUserRepository(db, messages)

	// Создаем тестового пользователя
	originalUser := &models.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     "employee",
		Created:  time.Now(),
	}

	err := repo.Create(context.Background(), originalUser)
	require.NoError(t, err)

	// Получаем пользователя по ID
	user, err := repo.GetByID(context.Background(), originalUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, originalUser.Email, user.Email)
	assert.Equal(t, originalUser.ID, user.ID)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	messages := ru.NewRussianMessages()
	repo := repositories.NewUserRepository(db, messages)

	// Пытаемся получить несуществующего пользователя
	nonExistentID := uuid.New()
	user, err := repo.GetByID(context.Background(), nonExistentID)
	assert.NoError(t, err)
	assert.Nil(t, user) // Должен вернуть nil, nil при отсутствии пользователя
}

func TestUserRepository_EmailExists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	messages := ru.NewRussianMessages()
	repo := repositories.NewUserRepository(db, messages)

	// Создаем тестового пользователя
	user := &models.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     "employee",
		Created:  time.Now(),
	}

	err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Проверяем существование email
	exists, err := repo.EmailExists(context.Background(), "test@example.com")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Проверяем несуществующий email
	exists, err = repo.EmailExists(context.Background(), "nonexistent@example.com")
	assert.NoError(t, err)
	assert.False(t, exists)
}
