package config

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string         `validate:"required"`
	Server   ServerConfig   `validate:"required"`
	Database DatabaseConfig `validate:"required"`
}

type ServerConfig struct {
	Port string `validate:"required,numeric"`
}

type DatabaseConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required,numeric"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Name     string `validate:"required"`
	SSLMode  string `validate:"required"`
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func validateConfig(cnf *Config) error {
	validate := validator.New()
	err := validate.Struct(cnf)

	if err != nil {
		return err
	}
	return nil
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	conf := &Config{
		Env: getEnv("ENV", "development"),
		Server: ServerConfig{
			Port: getEnv("PORT", "8001"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", ""),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}

	if err := validateConfig(conf); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return conf, nil
}
