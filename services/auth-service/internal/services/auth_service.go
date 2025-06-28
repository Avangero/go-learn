package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/avangero/auth-service/internal/lang"
	"github.com/avangero/auth-service/internal/models"
	"github.com/avangero/auth-service/internal/models/requests"
	"github.com/avangero/auth-service/internal/models/responses"
	"github.com/avangero/auth-service/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService интерфейс для сервиса аутентификации
type AuthService interface {
	Register(ctx context.Context, req *requests.RegisterRequest) (*responses.TokenResponse, error)
	Login(ctx context.Context, req *requests.LoginRequest) (*responses.TokenResponse, error)
	ValidateToken(ctx context.Context, tokenString string) (*models.User, error)
	GenerateToken(user *models.User) (string, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

// JWTClaims представляет данные в JWT токене
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// authService реализация AuthService
type authService struct {
	userRepo   repositories.UserRepository
	jwtSecret  string
	bcryptCost int
	messages   lang.Messages
}

// NewAuthService создает новый экземпляр AuthService
func NewAuthService(userRepo repositories.UserRepository, jwtSecret string, bcryptCost int, messages lang.Messages) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		bcryptCost: bcryptCost,
		messages:   messages,
	}
}

// Register регистрирует нового пользователя
func (s *authService) Register(ctx context.Context, req *requests.RegisterRequest) (*responses.TokenResponse, error) {
	log.Printf(s.messages.Get(lang.LogAttemptingRegistration), req.Email)

	// Проверяем существование пользователя
	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		log.Printf(s.messages.Get(lang.LogCheckEmailExists), req.Email, err)
		return nil, err
	}
	if exists {
		log.Printf(s.messages.Get(lang.LogEmailAlreadyExists), req.Email)
		return nil, errors.New(s.messages.Get(lang.UserAlreadyExists))
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.bcryptCost)
	if err != nil {
		log.Printf(s.messages.Get(lang.LogPasswordHashError), req.Email, err)
		return nil, err
	}

	// Создаем пользователя
	user := &models.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		Created:  time.Now(),
	}

	// Сохраняем в БД
	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Printf(s.messages.Get(lang.LogUserCreateError), req.Email, err)
		return nil, err
	}

	// Генерируем токен
	token, err := s.generateJWT(user)
	if err != nil {
		log.Printf(s.messages.Get(lang.LogJWTGenerateError), req.Email, err)
		return nil, err
	}

	log.Printf(s.messages.Get(lang.LogRegistrationComplete), req.Email)
	return &responses.TokenResponse{
		Token: token,
		User:  *user,
	}, nil
}

// Login аутентифицирует пользователя
func (s *authService) Login(ctx context.Context, req *requests.LoginRequest) (*responses.TokenResponse, error) {
	log.Printf(s.messages.Get(lang.LogAttemptingLogin), req.Email)

	// Ищем пользователя
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Printf(s.messages.Get(lang.LogDatabaseErrorLogin), req.Email, err)
		// Всегда возвращаем общую ошибку для безопасности
		return nil, errors.New(s.messages.Get(lang.InvalidCredentials))
	}

	// Проверяем, найден ли пользователь
	if user == nil {
		log.Printf(s.messages.Get(lang.LogUserNotFoundLogin), req.Email)
		return nil, errors.New(s.messages.Get(lang.InvalidCredentials))
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf(s.messages.Get(lang.LogInvalidPassword), req.Email)
		return nil, errors.New(s.messages.Get(lang.InvalidCredentials))
	}

	// Генерируем токен
	token, err := s.generateJWT(user)
	if err != nil {
		log.Printf(s.messages.Get(lang.LogJWTGenerateError), req.Email, err)
		return nil, err
	}

	log.Printf(s.messages.Get(lang.LogLoginComplete), req.Email)
	return &responses.TokenResponse{
		Token: token,
		User:  *user,
	}, nil
}

// ValidateToken проверяет валидность JWT токена и возвращает пользователя
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		log.Printf(s.messages.Get(lang.LogJWTParseError), err)
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		log.Printf(s.messages.Get(lang.LogJWTInvalid))
		return nil, errors.New(s.messages.Get(lang.TokenInvalid))
	}

	// Получаем актуальные данные пользователя из БД
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		log.Printf(s.messages.Get(lang.LogUserFetchError), claims.UserID.String(), err)
		return nil, err
	}

	if user == nil {
		log.Printf(s.messages.Get(lang.LogUserNotFoundValidation), claims.UserID.String())
		return nil, errors.New(s.messages.Get(lang.UserNotFound))
	}

	return user, nil
}

// GenerateToken генерирует JWT токен для пользователя (публичный метод)
func (s *authService) GenerateToken(user *models.User) (string, error) {
	return s.generateJWT(user)
}

// GetUserByID получает пользователя по ID
func (s *authService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// generateJWT генерирует JWT токен для пользователя (приватный метод)
func (s *authService) generateJWT(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
