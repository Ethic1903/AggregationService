package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DBConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type LoggerConfig struct {
	Level  string
	Format string
}

func MustLoad() *Config {
	_ = godotenv.Load(".env")

	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "7071"),
		},
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "admin"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "subscription_service_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
