package config

// Config содержит все настройки приложения
type Config struct {
	Port       string
	Database   DatabaseConfig
	JWT        JWTConfig
	BCryptCost int
}

// DatabaseConfig содержит настройки подключения к БД
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// JWTConfig содержит настройки JWT
type JWTConfig struct {
	Secret string
}
