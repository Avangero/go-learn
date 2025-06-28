package services

import (
	"errors"
	"time"

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
	Register(req *requests.RegisterRequest) (*responses.TokenResponse, error)
	Login(req *requests.LoginRequest) (*responses.TokenResponse, error)
	ValidateToken(tokenString string) (*models.User, error)
	GenerateToken(user *models.User) (string, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
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
}

// NewAuthService создает новый экземпляр AuthService
func NewAuthService(userRepo repositories.UserRepository, jwtSecret string, bcryptCost int) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		bcryptCost: bcryptCost,
	}
}

// Register регистрирует нового пользователя
func (s *authService) Register(req *requests.RegisterRequest) (*responses.TokenResponse, error) {
	// Проверяем существование пользователя
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.bcryptCost)
	if err != nil {
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
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Генерируем токен
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &responses.TokenResponse{
		Token: token,
		User:  *user,
	}, nil
}

// Login аутентифицирует пользователя
func (s *authService) Login(req *requests.LoginRequest) (*responses.TokenResponse, error) {
	// Ищем пользователя
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		// Всегда возвращаем общую ошибку для безопасности
		return nil, errors.New("неверный email или пароль")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	// Генерируем токен
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &responses.TokenResponse{
		Token: token,
		User:  *user,
	}, nil
}

// ValidateToken проверяет валидность JWT токена и возвращает пользователя
func (s *authService) ValidateToken(tokenString string) (*models.User, error) {
	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("недействительный токен")
	}

	// Получаем актуальные данные пользователя из БД
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GenerateToken генерирует JWT токен для пользователя (публичный метод)
func (s *authService) GenerateToken(user *models.User) (string, error) {
	return s.generateJWT(user)
}

// GetUserByID получает пользователя по ID
func (s *authService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(id)
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
