package config

import (
	"os"
)

type Config struct {
	Port            string
	GoogleClientID  string
	GoogleSecret    string
	JWTSecret       string
	WebSocketPort   string
	FrontendURL     string
	DatabaseURL     string
	Environment     string
}

func LoadConfig() *Config {
	return &Config{
		Port:           getEnvOrDefault("PORT", "8080"),
		GoogleClientID: getEnvOrDefault("GOOGLE_CLIENT_ID", ""),
		GoogleSecret:   getEnvOrDefault("GOOGLE_SECRET", ""),
		JWTSecret:      getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		WebSocketPort:  getEnvOrDefault("WS_PORT", "8081"),
		FrontendURL:    getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
		DatabaseURL:    getEnvOrDefault("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/polling?sslmode=disable"),
		Environment:    getEnvOrDefault("ENV", "development"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 