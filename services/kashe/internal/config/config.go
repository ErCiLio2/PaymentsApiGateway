package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port string

	// Database
	DatabaseURL string
	RedisURL    string

	// Proxypay API compat
	APIKey        string
	PosID         int
	WebhookSecret string
	BaseURL       string
	WebURL        string

	// Auth
	JWTSecret          string
	GoogleClientID     string
	GoogleClientSecret string

	// Email
	ResendAPIKey string
	EmailFrom    string

	// Encriptação (secrets TOTP em DB)
	EncryptionKey string

	// Admin
	AdminPassword string
	AdminKey      string
}

func Load() (*Config, error) {
	secret := getEnv("PROXYPAY_WEBHOOK_SECRET", "")
	if len(secret) < 32 {
		return nil, fmt.Errorf("PROXYPAY_WEBHOOK_SECRET deve ter pelo menos 32 caracteres")
	}

	return &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		APIKey:             getEnv("API_KEY", "dev-key"),
		PosID:              getEnvInt("POS_ID", 1),
		WebhookSecret:      secret,
		BaseURL:            getEnv("BASE_URL", "http://localhost:8080"),
		WebURL:             getEnv("WEB_URL", "http://localhost:5173"),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		ResendAPIKey:       getEnv("RESEND_API_KEY", ""),
		EmailFrom:          getEnv("EMAIL_FROM", "noreply@kashe.dev"),
		EncryptionKey:      getEnv("ENCRYPTION_KEY", ""),
		AdminPassword:      getEnv("ADMIN_PASSWORD", "admin"),
		AdminKey:           getEnv("ADMIN_KEY", "admin"),
	}, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
