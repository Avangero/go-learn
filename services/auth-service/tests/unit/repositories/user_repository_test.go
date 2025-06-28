package repositories_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/avangero/auth-service/internal/models"
	"github.com/avangero/auth-service/internal/repositories"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create_Success(t *testing.T) {
	// Подготовка mock базы данных
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.NewUserRepository(sqlxDB)

	user := &models.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     "employee",
		Created:  time.Now(),
	}

	// Ожидаем выполнение INSERT запроса
	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.ID, user.Email, user.Password, user.Role, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Выполнение
	err = repo.Create(user)

	// Проверка
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Create_DatabaseError(t *testing.T) {
	// Подготовка mock базы данных
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.NewUserRepository(sqlxDB)

	user := &models.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     "employee",
		Created:  time.Now(),
	}

	// Ожидаем ошибку при выполнении INSERT
	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.ID, user.Email, user.Password, user.Role, sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	// Выполнение
	err = repo.Create(user)

	// Проверка
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail_Success(t *testing.T) {
	// Подготовка mock базы данных
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.NewUserRepository(sqlxDB)

	expectedUser := &models.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     "employee",
		Created:  time.Now(),
	}

	// Ожидаем выполнение SELECT запроса
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "role", "created_at"}).
		AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Password, expectedUser.Role, expectedUser.Created)

	mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\$1").
		WithArgs(expectedUser.Email).
		WillReturnRows(rows)

	// Выполнение
	user, err := repo.GetByEmail(expectedUser.Email)

	// Проверка
	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Password, user.Password)
	assert.Equal(t, expectedUser.Role, user.Role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	// Подготовка mock базы данных
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.NewUserRepository(sqlxDB)

	email := "nonexistent@example.com"

	// Ожидаем выполнение SELECT запроса без результатов
	mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\$1").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	// Выполнение
	user, err := repo.GetByEmail(email)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID_Success(t *testing.T) {
	// Подготовка mock базы данных
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.NewUserRepository(sqlxDB)

	expectedUser := &models.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     "employee",
		Created:  time.Now(),
	}

	// Ожидаем выполнение SELECT запроса
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "role", "created_at"}).
		AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Password, expectedUser.Role, expectedUser.Created)

	mock.ExpectQuery("SELECT \\* FROM users WHERE id = \\$1").
		WithArgs(expectedUser.ID).
		WillReturnRows(rows)

	// Выполнение
	user, err := repo.GetByID(expectedUser.ID)

	// Проверка
	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_EmailExists_True(t *testing.T) {
	// Подготовка mock базы данных
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.NewUserRepository(sqlxDB)

	email := "existing@example.com"

	// Ожидаем выполнение EXISTS запроса
	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE email = \\$1\\)").
		WithArgs(email).
		WillReturnRows(rows)

	// Выполнение
	exists, err := repo.EmailExists(email)

	// Проверка
	require.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_EmailExists_False(t *testing.T) {
	// Подготовка mock базы данных
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.NewUserRepository(sqlxDB)

	email := "nonexistent@example.com"

	// Ожидаем выполнение EXISTS запроса
	rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE email = \\$1\\)").
		WithArgs(email).
		WillReturnRows(rows)

	// Выполнение
	exists, err := repo.EmailExists(email)

	// Проверка
	require.NoError(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}
