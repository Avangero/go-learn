package database

import (
	"fmt"
	"log"

	"github.com/avangero/auth-service/internal/config"
	"github.com/avangero/auth-service/internal/lang"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// ConnectionManager управляет подключениями к базе данных
type ConnectionManager struct {
	messages lang.Messages
}

// NewConnectionManager создает новый менеджер подключений
func NewConnectionManager(messages lang.Messages) *ConnectionManager {
	return &ConnectionManager{
		messages: messages,
	}
}

// Connect создает подключение к PostgreSQL
func (cm *ConnectionManager) Connect(cfg *config.Config) *sqlx.DB {
	// Строка подключения к PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(cm.messages.Get(lang.DBConnectionError)+":", err)
	}

	// Проверяем подключение
	if err = db.Ping(); err != nil {
		log.Fatal(cm.messages.Get(lang.DBPingError)+":", err)
	}

	log.Println(cm.messages.Get(lang.DBConnected))
	return db
}

// Connect создает подключение к PostgreSQL (удобная функция для обратной совместимости)
func Connect(cfg *config.Config) *sqlx.DB {
	// Создаем временный message provider для обратной совместимости
	// В идеале это должно быть внедрено через DI
	messages := map[lang.MessageKey]string{
		lang.DBConnectionError: "Ошибка подключения к базе данных",
		lang.DBPingError:       "Ошибка проверки подключения к БД",
		lang.DBConnected:       "✅ Подключение к PostgreSQL успешно",
	}

	cm := NewConnectionManager(lang.NewMessageProvider(messages))
	return cm.Connect(cfg)
}
