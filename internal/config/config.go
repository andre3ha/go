package config

import (
	"fmt"
	"os"
)

// Config - главная структура конфигурации приложения.
type Config struct {
	DB     DBConfig
	Server ServerConfig
}

// DBConfig - параметры подключения к PostgreSQL.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// ServerConfig - параметры HTTP-сервера.
type ServerConfig struct {
	Port string
}

// New - загружает конфиг из переменных окружения.
func New() (*Config, error) {
	return &Config{
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
	}, nil
}

// DSN - строка подключения к PostgreSQL для sqlx.Connect.
func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}
