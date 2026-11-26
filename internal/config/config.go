package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config stores all application configurations.
type Config struct {
	AppName string
	AppEnv  string
	AppPort int

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret string
}

// LoadConfig loads environment variables.
func LoadConfig() (*Config, error) {

	// Load .env file.
	// Environment variables already set by the OS are preserved.
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load .env: %w", err)
		}
	}

	// Read application port.
	portStr := getEnv("APP_PORT", "8080")

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid APP_PORT: %q", portStr)
	}

	// Read database port.
	dbPortStr := getEnv("DB_PORT", "5432")

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil || dbPort < 1 || dbPort > 65535 {
		return nil, fmt.Errorf("invalid DB_PORT: %q", dbPortStr)
	}

	// Create configuration object.
	cfg := &Config{
		AppName: getEnv("APP_NAME", "shoe-shop-api"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: port,

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "shoe_shop"),

		JWTSecret: getEnv("JWT_SECRET", ""),
	}

	return cfg, nil
}

// getEnv reads an environment variable or returns a default value.
func getEnv(key string, fallback string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
