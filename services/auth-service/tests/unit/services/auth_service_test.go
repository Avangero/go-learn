package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/avangero/auth-service/internal/lang/ru"
	"github.com/avangero/auth-service/internal/models"
	"github.com/avangero/auth-service/internal/models/requests"
	"github.com/avangero/auth-service/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository для тестирования
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func TestAuthService_Register_Success(t *testing.T) {
	// Подготовка
	mockRepo := new(MockUserRepository)
	messages := ru.NewRussianMessages()
	authService := services.NewAuthService(mockRepo, "test-secret", 4, messages)

	req := &requests.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Role:     "employee",
	}

	// Настройка моков
	mockRepo.On("EmailExists", mock.Anything, req.Email).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	// Выполнение
	tokenResponse, err := authService.Register(context.Background(), req)

	// Проверка
	require.NoError(t, err)
	assert.NotEmpty(t, tokenResponse.Token)
	assert.Equal(t, req.Email, tokenResponse.User.Email)
	assert.Equal(t, req.Role, tokenResponse.User.Role)
	assert.NotEmpty(t, tokenResponse.User.ID)
	assert.NotEmpty(t, tokenResponse.User.Created)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	// Подготовка
	mockRepo := new(MockUserRepository)
	messages := ru.NewRussianMessages()
	authService := services.NewAuthService(mockRepo, "test-secret", 4, messages)

	req := &requests.RegisterRequest{
		Email:    "existing@example.com",
		Password: "password123",
		Role:     "employee",
	}

	// Настройка моков
	mockRepo.On("EmailExists", mock.Anything, req.Email).Return(true, nil)

	// Выполнение
	tokenResponse, err := authService.Register(context.Background(), req)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, tokenResponse)
	assert.Contains(t, err.Error(), "уже существует")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_RepositoryError(t *testing.T) {
	// Подготовка
	mockRepo := new(MockUserRepository)
	messages := ru.NewRussianMessages()
	authService := services.NewAuthService(mockRepo, "test-secret", 4, messages)

	req := &requests.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Role:     "employee",
	}

	// Настройка моков
	mockRepo.On("EmailExists", mock.Anything, req.Email).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("database error"))

	// Выполнение
	tokenResponse, err := authService.Register(context.Background(), req)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, tokenResponse)
	assert.Contains(t, err.Error(), "database error")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	// Подготовка
	mockRepo := new(MockUserRepository)
	messages := ru.NewRussianMessages()
	authService := services.NewAuthService(mockRepo, "test-secret", 4, messages)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 4)

	existingUser := &models.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     "employee",
		Created:  time.Now(),
	}

	req := &requests.LoginRequest{
		Email:    existingUser.Email,
		Password: password,
	}

	// Настройка моков
	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	// Выполнение
	tokenResponse, err := authService.Login(context.Background(), req)

	// Проверка
	require.NoError(t, err)
	assert.NotEmpty(t, tokenResponse.Token)
	assert.Equal(t, existingUser.ID, tokenResponse.User.ID)
	assert.Equal(t, existingUser.Email, tokenResponse.User.Email)
	assert.Equal(t, existingUser.Role, tokenResponse.User.Role)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	// Подготовка
	mockRepo := new(MockUserRepository)
	messages := ru.NewRussianMessages()
	authService := services.NewAuthService(mockRepo, "test-secret", 4, messages)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), 4)

	existingUser := &models.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     "employee",
		Created:  time.Now(),
	}

	req := &requests.LoginRequest{
		Email:    existingUser.Email,
		Password: "wrong-password",
	}

	// Настройка моков
	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	// Выполнение
	tokenResponse, err := authService.Login(context.Background(), req)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, tokenResponse)
	assert.Contains(t, err.Error(), "Неверный email или пароль")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Подготовка
	mockRepo := new(MockUserRepository)
	messages := ru.NewRussianMessages()
	authService := services.NewAuthService(mockRepo, "test-secret", 4, messages)

	req := &requests.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	// Настройка моков - теперь возвращаем nil, nil для пользователя не найден
	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(nil, nil)

	// Выполнение
	tokenResponse, err := authService.Login(context.Background(), req)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, tokenResponse)
	assert.Contains(t, err.Error(), "Неверный email или пароль")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	// Подготовка
	mockRepo := new(MockUserRepository)
	messages := ru.NewRussianMessages()
	authService := services.NewAuthService(mockRepo, "test-secret", 4, messages)

	user := &models.User{
		ID:      uuid.New(),
		Email:   "test@example.com",
		Role:    "employee",
		Created: time.Now(),
	}

	// Создаем валидный токен
	token, err := authService.GenerateToken(user)
	require.NoError(t, err)

	// Настройка моков
	mockRepo.On("GetByID", mock.Anything, user.ID).Return(user, nil)

	// Выполнение
	validatedUser, err := authService.ValidateToken(context.Background(), token)

	// Проверка
	require.NoError(t, err)
	assert.Equal(t, user.ID, validatedUser.ID)
	assert.Equal(t, user.Email, validatedUser.Email)
	assert.Equal(t, user.Role, validatedUser.Role)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	// Подготовка
	mockRepo := new(MockUserRepository)
	messages := ru.NewRussianMessages()
	authService := services.NewAuthService(mockRepo, "test-secret", 4, messages)

	invalidToken := "invalid.jwt.token"

	// Выполнение
	user, err := authService.ValidateToken(context.Background(), invalidToken)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, user)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken_UserNotFound(t *testing.T) {
	// Подготовка
	mockRepo := new(MockUserRepository)
	messages := ru.NewRussianMessages()
	authService := services.NewAuthService(mockRepo, "test-secret", 4, messages)

	user := &models.User{
		ID:      uuid.New(),
		Email:   "test@example.com",
		Role:    "employee",
		Created: time.Now(),
	}

	// Создаем валидный токен
	token, err := authService.GenerateToken(user)
	require.NoError(t, err)

	// Настройка моков - пользователь не найден в БД
	mockRepo.On("GetByID", mock.Anything, user.ID).Return(nil, nil)

	// Выполнение
	validatedUser, err := authService.ValidateToken(context.Background(), token)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, validatedUser)
	assert.Contains(t, err.Error(), "Пользователь не найден")

	mockRepo.AssertExpectations(t)
}
